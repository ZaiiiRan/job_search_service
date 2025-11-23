package authservice

import (
	"context"
	"errors"

	pb "github.com/ZaiiiRan/job_search_service/auth-service/gen/go/auth_service/v1"
	userv1 "github.com/ZaiiiRan/job_search_service/auth-service/gen/go/user_service/v1"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/domain/code"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/domain/password"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/domain/token"
	userversion "github.com/ZaiiiRan/job_search_service/auth-service/internal/domain/user_version"
	uow "github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/unitofwork/postgres"
	codeservice "github.com/ZaiiiRan/job_search_service/auth-service/internal/services/code"
	passwordservice "github.com/ZaiiiRan/job_search_service/auth-service/internal/services/password"
	tokenservice "github.com/ZaiiiRan/job_search_service/auth-service/internal/services/token"
	userservice "github.com/ZaiiiRan/job_search_service/auth-service/internal/services/user_service"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/transport/postgres"
	"github.com/ZaiiiRan/job_search_service/common/pkg/ctxmetadata"
	claims "github.com/ZaiiiRan/job_search_service/common/pkg/jwt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	RegisterApplicant(ctx context.Context, req *pb.RegisterApplicantRequest) (*pb.RegisterApplicantResponse, error)
	GetNewApplicantActivationCode(ctx context.Context, req *pb.GetNewApplicantActivationCodeRequest) (*pb.GetNewApplicantActivationCodeResponse, error)
	ActivateApplicant(ctx context.Context, req *pb.ActivateApplicantRequest) (*pb.ActivateApplicantResponse, error)
	LoginApplicant(ctx context.Context, req *pb.LoginApplicantRequest) (*pb.LoginApplicantResponse, error)
	RefreshApplicant(ctx context.Context, req *pb.RefreshApplicantRequest) (*pb.RefreshApplicantResponse, error)
	LogoutApplicant(ctx context.Context, req *pb.LogoutApplicantRequest) (*pb.LogoutApplicantResponse, error)
	GetResetApplicantPasswordCode(ctx context.Context, req *pb.GetResetApplicantPasswordCodeRequest) (*pb.GetResetApplicantPasswordCodeResponse, error)
	ResetApplicantPassword(ctx context.Context, req *pb.ResetApplicantPasswordRequest) (*pb.ResetApplicantPasswordResponse, error)
	ChangeApplicantPassword(ctx context.Context, req *pb.ChangeApplicantPasswordRequest) (*pb.ChangeApplicantPasswordResponse, error)
}

type service struct {
	codeService     codeservice.CodeService
	passwordService passwordservice.PasswordService
	tokenService    tokenservice.TokenService
	userService     userservice.UserService
	postgresClient  *postgres.PostgresClient
	log             *zap.SugaredLogger
}

func New(
	postgresClient *postgres.PostgresClient,
	codeSvc codeservice.CodeService, passwordSvc passwordservice.PasswordService,
	tokenSvc tokenservice.TokenService, userSvc userservice.UserService,
	log *zap.SugaredLogger,
) AuthService {
	return &service{
		codeService:     codeSvc,
		passwordService: passwordSvc,
		tokenService:    tokenSvc,
		userService:     userSvc,
		postgresClient:  postgresClient,
		log:             log,
	}
}

func (s *service) RegisterApplicant(ctx context.Context, req *pb.RegisterApplicantRequest) (*pb.RegisterApplicantResponse, error) {
	l := s.log.With("op", "register_applicant", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	applicant, err := s.userService.CreateApplicant(ctx, req.Applicant)
	if err != nil {
		return nil, err
	}

	uow := uow.New(s.postgresClient)
	defer uow.Close()
	_, err = uow.BeginTransaction(ctx)
	if err != nil {
		l.Errorw("auth.register_applicant_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	_, err = s.passwordService.CreateApplicantPassword(ctx, uow, applicant, req.Password)
	if err != nil {
		var pve *password.PasswordValidationError
		if errors.As(err, &pve) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	uv, err := s.tokenService.CreateApplicantVersion(ctx, uow, applicant)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	_, err = s.codeService.CreateApplicantActivationCode(ctx, uow, applicant)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	if err := s.generateApplicantTokens(ctx, uow, applicant, uv, nil); err != nil {
		return nil, err
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("auth.register_applicant_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &pb.RegisterApplicantResponse{Applicant: applicant}, nil
}

func (s *service) GetNewApplicantActivationCode(ctx context.Context, req *pb.GetNewApplicantActivationCodeRequest) (*pb.GetNewApplicantActivationCodeResponse, error) {
	l := s.log.With("op", "get_new_applicant_activation_code", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	applicant, err := s.getAndCheckApplicantForActivation(ctx)
	if err != nil {
		return nil, err
	}

	uow := uow.New(s.postgresClient)
	defer uow.Close()

	_, err = s.codeService.RegenerateApplicantActivationCode(ctx, uow, applicant)
	if err != nil {
		var cve *code.CodeValidationError
		if errors.As(err, &cve) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	l.Infow("auth.get_new_activation_code.success")

	return &pb.GetNewApplicantActivationCodeResponse{}, nil
}

func (s *service) ActivateApplicant(ctx context.Context, req *pb.ActivateApplicantRequest) (*pb.ActivateApplicantResponse, error) {
	l := s.log.With("op", "activate_applicant", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	applicant, err := s.getAndCheckApplicantForActivation(ctx)
	if err != nil {
		return nil, err
	}

	uow := uow.New(s.postgresClient)
	defer uow.Close()
	_, err = uow.BeginTransaction(ctx)
	if err != nil {
		l.Errorw("auth.activate_applicant_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	valid, err := s.codeService.CheckApplicantActivationCode(ctx, uow, applicant, req.Code)
	if err != nil {
		var cve *code.CodeValidationError
		if errors.As(err, &cve) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	if !valid {
		return nil, status.Errorf(codes.InvalidArgument, "invalid code")
	}

	applicant, err = s.userService.ActivateApplicant(ctx, applicant)
	if err != nil {
		return nil, err
	}

	// invalidate all refresh tokens
	applicantVersion := userversion.New(applicant.Id)

	if err := s.generateApplicantTokens(ctx, uow, applicant, applicantVersion, nil); err != nil {
		return nil, err
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("auth.activate_applicant_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	l.Infow("auth.activate_applicant.success")
	return &pb.ActivateApplicantResponse{Applicant: applicant}, nil
}

func (s *service) LoginApplicant(ctx context.Context, req *pb.LoginApplicantRequest) (*pb.LoginApplicantResponse, error) {
	l := s.log.With("op", "login_applicant", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	applicant, err := s.userService.GetApplicantByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if applicant == nil || applicant.IsDeleted {
		return nil, status.Errorf(codes.Unauthenticated, "invalid email or password")
	}

	uow := uow.New(s.postgresClient)
	defer uow.Close()

	valid, err := s.passwordService.CheckApplicantPassword(ctx, uow, applicant, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	if !valid {
		return nil, status.Errorf(codes.Unauthenticated, "invalid email or password")
	}

	uv, err := s.tokenService.GetApplicantVersion(ctx, uow, applicant.Id)
	if err != nil || uv == nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	if err := s.generateApplicantTokens(ctx, uow, applicant, uv, nil); err != nil {
		return nil, err
	}

	l.Infow("auth.login_applicant.success")
	return &pb.LoginApplicantResponse{Applicant: applicant}, nil
}

func (s *service) RefreshApplicant(ctx context.Context, req *pb.RefreshApplicantRequest) (*pb.RefreshApplicantResponse, error) {
	l := s.log.With("op", "login_applicant", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}

	refreshTokenStr := md.Get("x-refresh-token")
	if len(refreshTokenStr) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}

	uow := uow.New(s.postgresClient)
	defer uow.Close()

	refreshToken, _, err := s.tokenService.ValidateApplicantRefreshToken(ctx, uow, refreshTokenStr[0])
	if err != nil {
		if errors.Is(err, claims.ErrInvalidToken) {
			return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
		}
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	applicant, err := s.userService.GetApplicantById(ctx, refreshToken.UserId())
	if err != nil {
		return nil, err
	}
	if applicant == nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}

	if err := s.generateApplicantTokens(ctx, uow, applicant, nil, refreshToken); err != nil {
		return nil, err
	}

	l.Infow("auth.refresh_applicant.success")
	return &pb.RefreshApplicantResponse{}, nil
}

func (s *service) LogoutApplicant(ctx context.Context, req *pb.LogoutApplicantRequest) (*pb.LogoutApplicantResponse, error) {
	l := s.log.With("op", "logout_applicant", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}

	refreshTokenStr := md.Get("x-refresh-token")
	if len(refreshTokenStr) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}

	uow := uow.New(s.postgresClient)
	defer uow.Close()

	s.tokenService.InvalidateApplicant(ctx, uow, refreshTokenStr[0])

	s.clearTokens(ctx)
	l.Infow("auth.logout_applicant.success")
	return &pb.LogoutApplicantResponse{}, nil
}

func (s *service) GetResetApplicantPasswordCode(ctx context.Context, req *pb.GetResetApplicantPasswordCodeRequest) (*pb.GetResetApplicantPasswordCodeResponse, error) {
	l := s.log.With("op", "get_reset_applicant_password_code", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	applicant, err := s.userService.GetApplicantByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if applicant == nil || applicant.IsDeleted {
		return &pb.GetResetApplicantPasswordCodeResponse{}, nil
	}

	uow := uow.New(s.postgresClient)
	defer uow.Close()

	_, err = s.codeService.RegenerateApplicantResetPasswordCode(ctx, uow, applicant)
	if err != nil {
		var cve *code.CodeValidationError
		if errors.As(err, &cve) {
			return &pb.GetResetApplicantPasswordCodeResponse{}, nil
		}
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	l.Infow("auth.get_reset_applicant_password_code.success")
	return &pb.GetResetApplicantPasswordCodeResponse{}, nil
}

func (s *service) ResetApplicantPassword(ctx context.Context, req *pb.ResetApplicantPasswordRequest) (*pb.ResetApplicantPasswordResponse, error) {
	l := s.log.With("op", "reset_applicant_password", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	applicant, err := s.userService.GetApplicantByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if applicant == nil || applicant.IsDeleted {
		return nil, status.Errorf(codes.InvalidArgument, "invalid email or code")
	}

	uow := uow.New(s.postgresClient)
	defer uow.Close()
	_, err = uow.BeginTransaction(ctx)
	if err != nil {
		l.Errorw("auth.reset_applicant_password_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	valid, err := s.codeService.CheckApplicantResetPasswordCode(ctx, uow, applicant, req.Code)
	if err != nil {
		var cve *code.CodeValidationError
		if errors.As(err, &cve) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	if !valid {
		return nil, status.Errorf(codes.InvalidArgument, "invalid email or code")
	}

	_, err = s.passwordService.UpdateApplicantPassword(ctx, uow, applicant, req.NewPassword)
	if err != nil {
		var pve *password.PasswordValidationError
		if errors.As(err, &pve) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	// invalidate all active tokens
	applicantVersion := userversion.New(applicant.Id)

	if err := s.generateApplicantTokens(ctx, uow, applicant, applicantVersion, nil); err != nil {
		return nil, err
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("auth.reset_applicant_password_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	l.Infow("auth.reset_applicant_password_failed.success")
	return &pb.ResetApplicantPasswordResponse{Applicant: applicant}, nil
}

func (s *service) ChangeApplicantPassword(ctx context.Context, req *pb.ChangeApplicantPasswordRequest) (*pb.ChangeApplicantPasswordResponse, error) {
	l := s.log.With("op", "change_applicant_password", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	claims, _ := ctxmetadata.GetApplicantClaimsFromContext(ctx)
	if claims == nil || claims.IsDeleted {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}

	applicant, err := s.userService.GetApplicantById(ctx, claims.Id)
	if err != nil {
		return nil, err
	}
	if applicant == nil || applicant.IsDeleted {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}

	if req.OldPassword == req.NewPassword {
		return nil, status.Errorf(codes.InvalidArgument, "old and new passwords are the same")
	}

	uow := uow.New(s.postgresClient)
	defer uow.Close()

	valid, err := s.passwordService.CheckApplicantPassword(ctx, uow, applicant, req.OldPassword)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	if !valid {
		return nil, status.Errorf(codes.InvalidArgument, "invalid old password")
	}

	_, err = s.passwordService.UpdateApplicantPassword(ctx, uow, applicant, req.NewPassword)
	if err != nil {
		var pve *password.PasswordValidationError
		if errors.As(err, &pve) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	// invalidate all active tokens
	applicantVersion := userversion.New(applicant.Id)

	if err := s.generateApplicantTokens(ctx, uow, applicant, applicantVersion, nil); err != nil {
		return nil, err
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("auth.change_applicant_password_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	l.Infow("auth.change_applicant_password_failed.success")
	return &pb.ChangeApplicantPasswordResponse{}, nil
}

func (s *service) generateApplicantTokens(
	ctx context.Context, uow *uow.UnitOfWork,
	applicant *userv1.Applicant, applicantVersion *userversion.UserVersion,
	existedRefreshToken *token.Token,
) error {
	access, refresh, err := s.tokenService.GenerateApplicant(ctx, uow, applicant, applicantVersion, existedRefreshToken)
	if err != nil {
		return status.Errorf(codes.Internal, "internal server error")
	}

	trailer := metadata.Pairs(
		"x-access-token", access.Token(),
		"x-refresh-token", refresh.Token(),
	)

	grpc.SetTrailer(ctx, trailer)
	return nil
}

func (s *service) getAndCheckApplicantForActivation(ctx context.Context) (*userv1.Applicant, error) {
	claims, _ := ctxmetadata.GetApplicantClaimsFromContext(ctx)
	if claims == nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}
	if claims.IsActive {
		return nil, status.Errorf(codes.AlreadyExists, "applicant is already activated")
	}
	if claims.IsDeleted {
		return nil, status.Errorf(codes.PermissionDenied, "applicant is deleted")
	}

	applicant, err := s.userService.GetApplicantById(ctx, claims.Id)
	if err != nil {
		return nil, err
	}
	if applicant == nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}
	if applicant.IsActive {
		return nil, status.Errorf(codes.AlreadyExists, "applicant is already activated")
	}
	if applicant.IsDeleted {
		return nil, status.Errorf(codes.PermissionDenied, "applicant is deleted")
	}
	return applicant, nil
}

func (s *service) clearTokens(ctx context.Context) {
	trailer := metadata.Pairs(
		"x-access-token", "",
		"x-refresh-token", "",
	)
	grpc.SetTrailer(ctx, trailer)
}
