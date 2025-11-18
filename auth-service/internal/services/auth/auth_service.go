package authservice

import (
	"context"
	"errors"

	pb "github.com/ZaiiiRan/job_search_service/auth-service/gen/go/auth_service/v1"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/domain/password"
	uow "github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/unitofwork/postgres"
	codeservice "github.com/ZaiiiRan/job_search_service/auth-service/internal/services/code"
	passwordservice "github.com/ZaiiiRan/job_search_service/auth-service/internal/services/password"
	tokenservice "github.com/ZaiiiRan/job_search_service/auth-service/internal/services/token"
	userservice "github.com/ZaiiiRan/job_search_service/auth-service/internal/services/user_service"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/transport/postgres"
	"github.com/ZaiiiRan/job_search_service/common/pkg/ctxmetadata"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	RegisterApplicant(ctx context.Context, req *pb.RegisterApplicantRequest) (*pb.RegisterApplicantResponse, error)
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
		log:             log,
	}
}

func (s *service) RegisterApplicant(ctx context.Context, req *pb.RegisterApplicantRequest) (*pb.RegisterApplicantResponse, error) {
	l := s.log.With("op", "register_applicant", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	applicant, err := s.userService.CreateApplicant(ctx, req.Applicant)
	if err != nil {
		l.Errorw("auth.register_applicant_failed", "err", err)
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

	_, err = s.codeService.CreateApplicantActivationCode(ctx, uow, applicant)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	access, refresh, err := s.tokenService.GenerateApplicant(ctx, uow, applicant, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	trailer := metadata.Pairs(
		"x-access-token", access.Token(),
		"x-refresh-token", refresh.Token(),
	)

	grpc.SetTrailer(ctx, trailer)
	return &pb.RegisterApplicantResponse{Applicant: applicant}, nil
}
