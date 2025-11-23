package tokenservice

import (
	"context"
	"fmt"
	"time"

	pb "github.com/ZaiiiRan/job_search_service/auth-service/gen/go/user_service/v1"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/config/settings"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/domain/token"
	userversion "github.com/ZaiiiRan/job_search_service/auth-service/internal/domain/user_version"
	uow "github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/transport/redis"
	"github.com/ZaiiiRan/job_search_service/common/pkg/ctxmetadata"
	claims "github.com/ZaiiiRan/job_search_service/common/pkg/jwt"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type TokenService interface {
	GenerateApplicant(ctx context.Context, uow *uow.UnitOfWork, applicant *pb.Applicant, applicantVersion *userversion.UserVersion,
		existedRefreshToken *token.Token) (access *token.Token, refresh *token.Token, err error)

	GenerateEmployer(ctx context.Context, uow *uow.UnitOfWork, employer *pb.Employer, employerVersion *userversion.UserVersion,
		existedRefreshToken *token.Token) (access *token.Token, refresh *token.Token, err error)

	ValidateApplicantRefreshToken(ctx context.Context, uow *uow.UnitOfWork, tokenStr string) (*token.Token, *userversion.UserVersion, error)
	ValidateEmployerRefreshToken(ctx context.Context, uow *uow.UnitOfWork, tokenStr string) (*token.Token, *userversion.UserVersion, error)
	ValidateApplicantAccessToken(ctx context.Context, tokenStr string) (*claims.ApplicantClaims, error)
	ValidateEmployerAccessToken(ctx context.Context, tokenStr string) (*claims.EmployerClaims, error)
	InvalidateApplicant(ctx context.Context, uow *uow.UnitOfWork, refreshStr string) error
	InvalidateEmployer(ctx context.Context, uow *uow.UnitOfWork, refreshStr string) error
	GetApplicantVersion(ctx context.Context, uow *uow.UnitOfWork, userId int64) (*userversion.UserVersion, error)
	GetEmployerVersion(ctx context.Context, uow *uow.UnitOfWork, userId int64) (*userversion.UserVersion, error)
	CreateApplicantVersion(ctx context.Context, uow *uow.UnitOfWork, applicant *pb.Applicant) (*userversion.UserVersion, error)
	CreateEmployerVersion(ctx context.Context, uow *uow.UnitOfWork, employer *pb.Employer) (*userversion.UserVersion, error)
}

type service struct {
	tokenDataProvider       *tokenDataProvider
	userVersionDataProvider *userVersionDataProvider
	jwtSettings             *settings.JWTSettings
	log                     *zap.SugaredLogger
}

func New(jwtSettings settings.JWTSettings, redis *redis.RedisClient, log *zap.SugaredLogger) TokenService {
	return &service{
		tokenDataProvider:       newTokenDataProvider(redis),
		userVersionDataProvider: newUserVersionDataProvider(redis),
		jwtSettings:             &jwtSettings,
		log:                     log,
	}
}

func (s *service) GenerateApplicant(
	ctx context.Context, uow *uow.UnitOfWork,
	applicant *pb.Applicant, applicantVersion *userversion.UserVersion,
	existedRefreshToken *token.Token,
) (*token.Token, *token.Token, error) {
	l := s.log.With("op", "generate_tokens_for_applicant", "req_id", ctxmetadata.GetReqIdFromContext(ctx), "applicant_id", applicant.Id)

	var version int
	if existedRefreshToken != nil {
		version = existedRefreshToken.Version()
	} else if applicantVersion != nil {
		version = applicantVersion.Version()
	} else {
		l.Errorw("token.get_user_version", "err", "user version or existed refresh token is not provided")
		return nil, nil, fmt.Errorf("user version or existed refresh token is not provided")
	}

	c := &claims.ApplicantClaims{
		Id:         applicant.Id,
		FirstName:  applicant.FirstName,
		LastName:   applicant.LastName,
		Patronymic: applicant.Patronymic,
		Email:      applicant.Email,
		IsActive:   applicant.IsActive,
		IsDeleted:  applicant.IsDeleted,
		Version:    version,
	}

	access, accessExp, err := signToken(c, []byte(s.jwtSettings.AccessTokenSecret), time.Duration(s.jwtSettings.AccessTokenTTL)*time.Second)
	if err != nil {
		l.Errorw("token.sign_access_failed", "err", err)
		return nil, nil, err
	}

	refresh, refreshExp, err := signToken(c, []byte(s.jwtSettings.RefreshTokenSecret), time.Duration(s.jwtSettings.RefreshTokenTTL)*time.Second)
	if err != nil {
		l.Errorw("token.sign_refresh_failed", "err", err)
		return nil, nil, err
	}

	accessToken := token.New(applicant.Id, access, token.AccessTokenType, 0, accessExp)

	var refreshToken *token.Token
	if existedRefreshToken != nil {
		if err := s.tokenDataProvider.DeleteApplicantTokenFromCache(ctx, existedRefreshToken.Token()); err != nil {
			l.Errorw("token.delete_existed_from_cache", "err", err)
			return nil, nil, err
		}
		refreshToken = existedRefreshToken
		refreshToken.SetToken(refresh, refreshExp)
	} else {
		refreshToken = token.New(applicant.Id, refresh, token.RefreshTokenType, version, refreshExp)
	}

	if err := s.tokenDataProvider.SaveApplicantToken(ctx, uow, refreshToken); err != nil {
		l.Errorw("token.save_token_failed", "err", err)
		return nil, nil, err
	}

	l.Infow("token.generated_tokens_for_applicant")
	return accessToken, refreshToken, nil
}

func (s *service) GenerateEmployer(
	ctx context.Context, uow *uow.UnitOfWork,
	employer *pb.Employer, employerVersion *userversion.UserVersion,
	existedRefreshToken *token.Token,
) (*token.Token, *token.Token, error) {
	l := s.log.With("op", "generate_tokens_for_employer", "req_id", ctxmetadata.GetReqIdFromContext(ctx), "employer_id", employer.Id)

	var version int
	if existedRefreshToken != nil {
		version = existedRefreshToken.Version()
	} else if employerVersion != nil {
		version = employerVersion.Version()
	} else {
		l.Errorw("token.get_user_version", "err", "user version or existed refresh token is not provided")
		return nil, nil, fmt.Errorf("user version or existed refresh token is not provided")
	}

	c := &claims.EmployerClaims{
		Id:          employer.Id,
		CompanyName: employer.CompanyName,
		Email:       employer.Email,
		IsActive:    employer.IsActive,
		IsDeleted:   employer.IsDeleted,
		Version:     version,
	}

	access, accessExp, err := signToken(c, []byte(s.jwtSettings.AccessTokenSecret), time.Duration(s.jwtSettings.AccessTokenTTL)*time.Second)
	if err != nil {
		l.Errorw("token.sign_access_failed", "err", err)
		return nil, nil, err
	}

	refresh, refreshExp, err := signToken(c, []byte(s.jwtSettings.RefreshTokenSecret), time.Duration(s.jwtSettings.RefreshTokenTTL)*time.Second)
	if err != nil {
		l.Errorw("token.sign_refresh_failed", "err", err)
		return nil, nil, err
	}

	accessToken := token.New(employer.Id, access, token.AccessTokenType, 0, accessExp)

	var refreshToken *token.Token
	if existedRefreshToken != nil {
		if err := s.tokenDataProvider.DeleteEmployerTokenFromCache(ctx, existedRefreshToken.Token()); err != nil {
			l.Errorw("token.delete_existed_from_cache", "err", err)
			return nil, nil, err
		}
		refreshToken = existedRefreshToken
		refreshToken.SetToken(refresh, refreshExp)
	} else {
		refreshToken = token.New(employer.Id, refresh, token.RefreshTokenType, version, refreshExp)
	}

	if err := s.tokenDataProvider.SaveEmployerToken(ctx, uow, refreshToken); err != nil {
		l.Errorw("token.save_token_failed", "err", err)
		return nil, nil, err
	}

	l.Infow("token.generated_tokens_for_employer")
	return accessToken, refreshToken, nil
}

func (s *service) ValidateApplicantRefreshToken(ctx context.Context, uow *uow.UnitOfWork, tokenStr string) (*token.Token, *userversion.UserVersion, error) {
	l := s.log.With("op", "validate_applicant_refresh_token", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	cl, err := claims.ParseApplicantToken(tokenStr, []byte(s.jwtSettings.RefreshTokenSecret))
	if err != nil {
		l.Warnw("token.refresh_token_parse_failed", "err", err)
		return nil, nil, claims.ErrInvalidToken
	}

	t, err := s.tokenDataProvider.GetApplicantToken(ctx, uow, tokenStr)
	if err != nil {
		l.Errorw("token.get_token_failed", "err", err)
		return nil, nil, err
	}
	if t == nil || cl.Id != t.UserId() || cl.Version != t.Version() {
		l.Warnw("token.refresh_token_invalid")
		return nil, nil, claims.ErrInvalidToken
	}

	applicantVersion, err := s.userVersionDataProvider.GetApplicantVersion(ctx, uow, t.UserId())
	if err != nil {
		l.Errorw("token.get_user_version_failed", "err", err)
		return nil, nil, err
	}
	if applicantVersion == nil || t.Version() != applicantVersion.Version() {
		l.Warnw("token.refresh_token_invalid")
		return nil, nil, claims.ErrInvalidToken
	}

	l.Infow("token.refresh_token_valid", "applicant_id", cl.Id)
	return t, applicantVersion, nil
}

func (s *service) ValidateEmployerRefreshToken(ctx context.Context, uow *uow.UnitOfWork, tokenStr string) (*token.Token, *userversion.UserVersion, error) {
	l := s.log.With("op", "validate_employer_refresh_token", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	cl, err := claims.ParseEmployerToken(tokenStr, []byte(s.jwtSettings.RefreshTokenSecret))
	if err != nil {
		l.Warnw("token.refresh_token_parse_failed", "err", err)
		return nil, nil, claims.ErrInvalidToken
	}

	t, err := s.tokenDataProvider.GetEmployerToken(ctx, uow, tokenStr)
	if err != nil {
		l.Errorw("token.get_token_failed", "err", err)
		return nil, nil, err
	}
	if t == nil || cl.Id != t.UserId() || cl.Version != t.Version() {
		l.Warnw("token.refresh_token_invalid")
		return nil, nil, claims.ErrInvalidToken
	}

	employerVersion, err := s.userVersionDataProvider.GetEmployerVersion(ctx, uow, t.UserId())
	if err != nil {
		l.Errorw("token.get_user_version_failed", "err", err)
		return nil, nil, err
	}
	if employerVersion == nil || t.Version() != employerVersion.Version() {
		l.Warnw("token.refresh_token_invalid")
		return nil, nil, claims.ErrInvalidToken
	}

	l.Infow("token.refresh_token_valid", "employer_id", cl.Id)
	return t, employerVersion, nil
}

func (s *service) ValidateApplicantAccessToken(ctx context.Context, tokenStr string) (*claims.ApplicantClaims, error) {
	l := s.log.With("op", "validate_applicant_access_token", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	cl, err := claims.ParseApplicantToken(tokenStr, []byte(s.jwtSettings.AccessTokenSecret))
	if err != nil {
		l.Warnw("token.access_token_parse_failed", "err", err)
		return nil, claims.ErrInvalidToken
	}

	l.Infow("token.access_token_valid", "applicant_id", cl.Id)
	return cl, nil
}

func (s *service) ValidateEmployerAccessToken(ctx context.Context, tokenStr string) (*claims.EmployerClaims, error) {
	l := s.log.With("op", "validate_employer_access_token", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	cl, err := claims.ParseEmployerToken(tokenStr, []byte(s.jwtSettings.AccessTokenSecret))
	if err != nil {
		l.Warnw("token.access_token_parse_failed", "err", err)
		return nil, claims.ErrInvalidToken
	}

	l.Infow("token.access_token_valid", "employer_id", cl.Id)
	return cl, nil
}

func (s *service) InvalidateApplicant(ctx context.Context, uow *uow.UnitOfWork, refreshStr string) error {
	l := s.log.With("op", "invalidate_applicant_refresh_token", "req_id", ctxmetadata.GetReqIdFromContext(ctx))
	err := s.tokenDataProvider.DeleteApplicantToken(ctx, uow, refreshStr)
	if err != nil {
		l.Errorw("token.delete_refresh_token_failed", "err", err)
		return err
	}

	l.Infow("token.invalidate_refresh_token.success")
	return nil
}

func (s *service) InvalidateEmployer(ctx context.Context, uow *uow.UnitOfWork, refreshStr string) error {
	l := s.log.With("op", "invalidate_employer_refresh_token", "req_id", ctxmetadata.GetReqIdFromContext(ctx))
	err := s.tokenDataProvider.DeleteEmployerToken(ctx, uow, refreshStr)
	if err != nil {
		l.Errorw("token.delete_refresh_token_failed", "err", err)
		return err
	}

	l.Infow("token.invalidate_refresh_token.success")
	return nil
}

func (s *service) GetApplicantVersion(ctx context.Context, uow *uow.UnitOfWork, userId int64) (*userversion.UserVersion, error) {
	l := s.log.With("op", "get_applicant_version", "req_id", ctxmetadata.GetReqIdFromContext(ctx), "applicant_id", userId)
	uv, err := s.userVersionDataProvider.GetApplicantVersion(ctx, uow, userId)
	if err != nil {
		l.Errorw("token.get_user_version_failed", "err", err)
		return nil, err
	}
	l.Infow("token.get_user_version.success")
	return uv, nil
}

func (s *service) GetEmployerVersion(ctx context.Context, uow *uow.UnitOfWork, userId int64) (*userversion.UserVersion, error) {
	l := s.log.With("op", "get_employer_version", "req_id", ctxmetadata.GetReqIdFromContext(ctx), "employer_id", userId)
	uv, err := s.userVersionDataProvider.GetEmployerVersion(ctx, uow, userId)
	if err != nil {
		l.Errorw("token.get_user_version_failed", "err", err)
		return nil, err
	}
	l.Infow("token.get_user_version.success")
	return uv, nil
}

func (s *service) CreateApplicantVersion(ctx context.Context, uow *uow.UnitOfWork, applicant *pb.Applicant) (*userversion.UserVersion, error) {
	l := s.log.With("op", "create_applicant_version", "req_id", ctxmetadata.GetReqIdFromContext(ctx), "applicant_id", applicant.Id)

	var uv *userversion.UserVersion

	existedUserVersion, err := s.userVersionDataProvider.GetApplicantVersion(ctx, uow, applicant.Id)
	if err != nil {
		l.Errorw("token.get_existed_user_version_failed", "err", err)
		return nil, err
	}

	if existedUserVersion != nil {
		uv = existedUserVersion
		uv.IncrementVersion()
	} else {
		uv = userversion.New(applicant.Id)
	}

	if err := s.userVersionDataProvider.SaveApplicantVersion(ctx, uow, uv); err != nil {
		l.Errorw("token.save_user_version", "err", err)
		return nil, err
	}

	l.Infow("token.create_user_version.success")
	return uv, nil
}

func (s *service) CreateEmployerVersion(ctx context.Context, uow *uow.UnitOfWork, employer *pb.Employer) (*userversion.UserVersion, error) {
	l := s.log.With("op", "create_employer_version", "req_id", ctxmetadata.GetReqIdFromContext(ctx), "employer_id", employer.Id)

	var uv *userversion.UserVersion

	existedUserVersion, err := s.userVersionDataProvider.GetEmployerVersion(ctx, uow, employer.Id)
	if err != nil {
		l.Errorw("token.get_existed_user_version_failed", "err", err)
		return nil, err
	}

	if existedUserVersion != nil {
		uv = existedUserVersion
		uv.IncrementVersion()
	} else {
		uv = userversion.New(employer.Id)
	}

	if err := s.userVersionDataProvider.SaveEmployerVersion(ctx, uow, uv); err != nil {
		l.Errorw("token.save_user_version", "err", err)
		return nil, err
	}

	l.Infow("token.create_user_version.success")
	return uv, nil
}

func signToken[T jwt.Claims](c T, key []byte, ttl time.Duration) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(ttl)
	safeNbf := now.Add(-10 * time.Second)
	switch v := any(c).(type) {
	case *claims.ApplicantClaims:
		v.RegisteredClaims = jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(safeNbf),
		}
	case *claims.EmployerClaims:
		v.RegisteredClaims = jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(safeNbf),
		}
	default:
		return "", expiresAt, fmt.Errorf("unknown claims type")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	str, err := token.SignedString(key)
	if err != nil {
		return "", expiresAt, err
	}

	return str, expiresAt, nil
}
