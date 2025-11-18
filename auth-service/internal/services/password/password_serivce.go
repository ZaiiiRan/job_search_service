package passwordservice

import (
	"context"
	"fmt"

	pb "github.com/ZaiiiRan/job_search_service/auth-service/gen/go/user_service/v1"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/domain/password"
	uow "github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/transport/redis"
	"github.com/ZaiiiRan/job_search_service/common/pkg/ctxmetadata"
	"go.uber.org/zap"
)

type PasswordService interface {
	CreateApplicantPassword(ctx context.Context, uow *uow.UnitOfWork, applicant *pb.Applicant, rawPassword string) (*password.Password, error)
	CreateEmployerPassword(ctx context.Context, uow *uow.UnitOfWork, employer *pb.Employer, rawPassword string) (*password.Password, error)
	CheckApplicantPassword(ctx context.Context, uow *uow.UnitOfWork, applicant *pb.Applicant, rawPassword string) (bool, error)
	CheckEmployerPassword(ctx context.Context, uow *uow.UnitOfWork, employer *pb.Employer, rawPassword string) (bool, error)
	UpdateApplicantPassword(ctx context.Context, uow *uow.UnitOfWork, applicant *pb.Applicant, rawPassword string) (*password.Password, error)
	UpdateEmployerPassword(ctx context.Context, uow *uow.UnitOfWork, employer *pb.Employer, rawPassword string) (*password.Password, error)
}

type service struct {
	dataProvider *passwordDataProvider
	log          *zap.SugaredLogger
}

func New(redis *redis.RedisClient, log *zap.SugaredLogger) PasswordService {
	return &service{
		dataProvider: newPasswordDataProvider(redis),
		log:          log,
	}
}

func (s *service) CreateApplicantPassword(
	ctx context.Context, uow *uow.UnitOfWork,
	applicant *pb.Applicant, rawPassword string,
) (*password.Password, error) {
	l := s.log.With("op", "create_applicant_password", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	password, err := password.New(applicant.Id, rawPassword)
	if err != nil {
		l.Warnw("password.create_password_failed.validation_error", "err", err)
		return nil, err
	}

	existedPassword, err := s.dataProvider.GetApplicantPasswordByUserId(ctx, uow, applicant.Id)
	if err != nil {
		l.Errorw("password.create_password_failed", "err", err)
		return nil, err
	}
	if existedPassword != nil {
		password.SetId(existedPassword.Id())
	}

	if err := s.dataProvider.SaveApplicantPassword(ctx, uow, password); err != nil {
		l.Errorw("password.create_password_failed", "err", err)
		return nil, err
	}

	l.Infow("password.create_password.success")
	return password, nil
}

func (s *service) CreateEmployerPassword(
	ctx context.Context, uow *uow.UnitOfWork,
	employer *pb.Employer, rawPassword string,
) (*password.Password, error) {
	l := s.log.With("op", "create_employer_password", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	password, err := password.New(employer.Id, rawPassword)
	if err != nil {
		l.Warnw("password.create_password_failed.validation_error", "err", err)
		return nil, err
	}

	existedPassword, err := s.dataProvider.GetEmployerPasswordByUserId(ctx, uow, employer.Id)
	if err != nil {
		l.Errorw("password.create_password_failed", "err", err)
		return nil, err
	}
	if existedPassword != nil {
		password.SetId(existedPassword.Id())
	}

	if err := s.dataProvider.SaveEmployerPassword(ctx, uow, password); err != nil {
		l.Errorw("password.create_password_failed", "err", err)
		return nil, err
	}

	l.Infow("password.create_password.success")
	return password, nil
}

func (s *service) CheckApplicantPassword(
	ctx context.Context, uow *uow.UnitOfWork,
	applicant *pb.Applicant, rawPassword string,
) (bool, error) {
	l := s.log.With("op", "check_applicant_password", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	password, err := s.dataProvider.GetApplicantPasswordByUserId(ctx, uow, applicant.Id)
	if err != nil {
		l.Errorw("password.check_password_failed", "err", err)
		return false, err
	}
	if password == nil {
		l.Errorw("password.check_password_failed", "err", "password not found")
		return false, fmt.Errorf("password not found")
	}

	correct := password.Check(rawPassword)
	l.Infow("password.check_password.success")
	return correct, nil
}

func (s *service) CheckEmployerPassword(
	ctx context.Context, uow *uow.UnitOfWork,
	employer *pb.Employer, rawPassword string,
) (bool, error) {
	l := s.log.With("op", "check_employer_password", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	password, err := s.dataProvider.GetEmployerPasswordByUserId(ctx, uow, employer.Id)
	if err != nil {
		l.Errorw("password.check_password_failed", "err", err)
		return false, err
	}
	if password == nil {
		l.Errorw("password.check_password_failed", "err", "password not found")
		return false, fmt.Errorf("password not found")
	}

	correct := password.Check(rawPassword)
	l.Infow("password.check_password.success")
	return correct, nil
}

func (s *service) UpdateApplicantPassword(
	ctx context.Context, uow *uow.UnitOfWork,
	applicant *pb.Applicant, rawPassword string,
) (*password.Password, error) {
	l := s.log.With("op", "update_applicant_password", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	password, err := s.dataProvider.GetApplicantPasswordByUserId(ctx, uow, applicant.Id)
	if err != nil {
		l.Errorw("password.update_password_failed", "err", err)
		return nil, err
	}
	if password == nil {
		l.Errorw("password.update_password_failed", "err", "password not found")
		return nil, fmt.Errorf("password not found")
	}

	if err := password.SetPassword(rawPassword); err != nil {
		l.Warnw("password.update_password_failed.validation_error", "err", err)
		return nil, err
	}

	if err := s.dataProvider.SaveApplicantPassword(ctx, uow, password); err != nil {
		l.Errorw("password.update_password_failed", "err", err)
		return nil, err
	}

	l.Infow("password.update_password.success")
	return password, nil
}

func (s *service) UpdateEmployerPassword(
	ctx context.Context, uow *uow.UnitOfWork,
	employer *pb.Employer, rawPassword string,
) (*password.Password, error) {
	l := s.log.With("op", "update_employer_password", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	password, err := s.dataProvider.GetEmployerPasswordByUserId(ctx, uow, employer.Id)
	if err != nil {
		l.Errorw("password.update_password_failed", "err", err)
		return nil, err
	}
	if password == nil {
		l.Errorw("password.update_password_failed", "err", "password not found")
		return nil, fmt.Errorf("password not found")
	}

	if err := password.SetPassword(rawPassword); err != nil {
		l.Warnw("password.update_password_failed.validation_error", "err", err)
		return nil, err
	}

	if err := s.dataProvider.SaveEmployerPassword(ctx, uow, password); err != nil {
		l.Errorw("password.update_password_failed", "err", err)
		return nil, err
	}

	l.Infow("password.update_password.success")
	return password, nil
}
