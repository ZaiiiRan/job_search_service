package tokenservice

import (
	"context"
	"fmt"
	"time"

	pb "github.com/ZaiiiRan/job_search_service/auth-service/gen/go/user_service/v1"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/config/settings"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/domain/token"
	"github.com/ZaiiiRan/job_search_service/common/pkg/ctxmetadata"
	claims "github.com/ZaiiiRan/job_search_service/common/pkg/jwt"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type TokenService interface {
	GenerateApplicant(ctx context.Context, applicant *pb.Applicant, existedRefreshToken *token.Token) (access *token.Token, refresh *token.Token, err error)
	GenerateEmployer(ctx context.Context, employer *pb.Employer, existedRefreshToken *token.Token) (access *token.Token, refresh *token.Token, err error)
	ValidateApplicantRefreshToken(ctx context.Context, tokenStr string) (*token.Token, error)
	ValidateEmployerRefreshToken(ctx context.Context, tokenStr string) (*token.Token, error)
	ValidateApplicantAccessToken(ctx context.Context, tokenStr string) (*claims.ApplicantClaims, error)
	ValidateEmployerAccessToken(ctx context.Context, tokenStr string) (*claims.EmployerClaims, error)
	InvalidateApplicant(ctx context.Context, refreshStr string) error
	InvalidateEmployer(ctx context.Context, refreshStr string) error
}

type service struct {
	jwtSettings *settings.JWTSettings
	log         *zap.SugaredLogger
}

func New(jwtSettings settings.JWTSettings, log *zap.SugaredLogger) TokenService {
	return &service{
		jwtSettings: &jwtSettings,
		log:         log,
	}
}

func (s *service) GenerateApplicant(ctx context.Context, applicant *pb.Applicant, existedRefreshToken *token.Token,) (*token.Token, *token.Token, error) {
	l := s.log.With("op", "generate_tokens_for_applicant", "req_id", ctxmetadata.GetReqIdFromContext(ctx), "applicant_id", applicant.Id)

	c := &claims.ApplicantClaims{
		Id:         applicant.Id,
		FirstName:  applicant.FirstName,
		LastName:   applicant.LastName,
		Patronymic: applicant.Patronymic,
		Email:      applicant.Email,
		IsActive:   applicant.IsActive,
		IsDeleted:  applicant.IsDeleted,
	}

	access, accessExp, err := signToken(c, []byte(s.jwtSettings.AccessTokenSecret), time.Duration(s.jwtSettings.AccessTokenTTL))
	if err != nil {
		l.Errorw("token.sign_access_failed", "err", err)
		return nil, nil, err
	}

	refresh, refreshExp, err := signToken(c, []byte(s.jwtSettings.RefreshTokenSecret), time.Duration(s.jwtSettings.RefreshTokenTTL))
	if err != nil {
		l.Errorw("token.sign_refresh_failed", "err", err)
		return nil, nil, err
	}

	accessToken := token.New(applicant.Id, access, token.AccessTokenType, accessExp)

	var refreshToken *token.Token
	if existedRefreshToken != nil {
		refreshToken = existedRefreshToken
		refreshToken.SetToken(refresh, refreshExp)
	} else {
		refreshToken = token.New(applicant.Id, refresh, token.RefreshTokenType, refreshExp)
	}

	l.Infow("token.generated_tokens_for_applicant")
	return accessToken, refreshToken, nil
}

func (s *service) GenerateEmployer(ctx context.Context, employer *pb.Employer, existedRefreshToken *token.Token) (*token.Token, *token.Token, error) {
	l := s.log.With("op", "generate_tokens_for_employer", "req_id", ctxmetadata.GetReqIdFromContext(ctx), "employer_id", employer.Id)

	c := &claims.EmployerClaims{
		Id:          employer.Id,
		CompanyName: employer.CompanyName,
		Email:       employer.Email,
		IsActive:    employer.IsActive,
		IsDeleted:   employer.IsDeleted,
	}

	access, accessExp, err := signToken(c, []byte(s.jwtSettings.AccessTokenSecret), time.Duration(s.jwtSettings.AccessTokenTTL))
	if err != nil {
		l.Errorw("token.sign_access_failed", "err", err)
		return nil, nil, err
	}

	refresh, refreshExp, err := signToken(c, []byte(s.jwtSettings.RefreshTokenSecret), time.Duration(s.jwtSettings.RefreshTokenTTL))
	if err != nil {
		l.Errorw("token.sign_refresh_failed", "err", err)
		return nil, nil, err
	}

	accessToken := token.New(employer.Id, access, token.AccessTokenType, accessExp)

	var refreshToken *token.Token
	if existedRefreshToken != nil {
		refreshToken = existedRefreshToken
		refreshToken.SetToken(refresh, refreshExp)
	} else {
		refreshToken = token.New(employer.Id, refresh, token.RefreshTokenType, refreshExp)
	}

	l.Infow("token.generated_tokens_for_employer")

	return accessToken, refreshToken, nil
}

func (s *service) ValidateApplicantRefreshToken(ctx context.Context, tokenStr string) (*token.Token, error) {
	l := s.log.With("op", "validate_applicant_refresh_token", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	cl, err := claims.ParseApplicantToken(tokenStr, []byte(s.jwtSettings.RefreshTokenSecret))
	if err != nil {
		l.Warnw("token.refresh_token_parse_failed", "err", err)
		return nil, claims.ErrInvalidToken
	}

	// check token in repo

	l.Infow("token.refresh_token_valid", "applicant_id", cl.Id)
	return &token.Token{}, nil
}

func (s *service) ValidateEmployerRefreshToken(ctx context.Context, tokenStr string) (*token.Token, error) {
	l := s.log.With("op", "validate_employer_refresh_token", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	cl, err := claims.ParseEmployerToken(tokenStr, []byte(s.jwtSettings.RefreshTokenSecret))
	if err != nil {
		l.Warnw("token.refresh_token_parse_failed", "err", err)
		return nil, claims.ErrInvalidToken
	}

	// check token in repo

	l.Infow("token.refresh_token_valid", "employer_id", cl.Id)
	return &token.Token{}, nil
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

func (s *service) InvalidateApplicant(ctx context.Context, refreshStr string) error {
	l := s.log.With("op", "applicant_token_invalidate", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	// delete token from repo

	l.Infow("token.applicant_token_invalidate_ok")
	return nil
}

func (s *service) InvalidateEmployer(ctx context.Context, refreshStr string) error {
	l := s.log.With("op", "employer_token_invalidate", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	// delete token from repo

	l.Infow("token.employer_token_invalidate_ok")
	return nil
}

func signToken[T jwt.Claims](c T, key []byte, ttl time.Duration) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(ttl)
	switch v := any(c).(type) {
	case *claims.ApplicantClaims:
		v.RegisteredClaims = jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		}
	case *claims.EmployerClaims:
		v.RegisteredClaims = jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
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
