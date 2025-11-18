package codeservice

import (
	"context"

	pb "github.com/ZaiiiRan/job_search_service/auth-service/gen/go/user_service/v1"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/domain/code"
	uow "github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/transport/redis"
	"github.com/ZaiiiRan/job_search_service/common/pkg/ctxmetadata"
	"go.uber.org/zap"
)

type CodeService interface {
	CreateApplicantActivationCode(ctx context.Context, uow *uow.UnitOfWork, applicant *pb.Applicant) (*code.Code, error)
	CreateEmployerActivationCode(ctx context.Context, uow *uow.UnitOfWork, employer *pb.Employer) (*code.Code, error)
	CheckApplicantActivationCode(ctx context.Context, uow *uow.UnitOfWork, applicant *pb.Applicant, rawCode string) (bool, error)
	CheckEmployerActivationCode(ctx context.Context, uow *uow.UnitOfWork, employer *pb.Employer, rawCode string) (bool, error)
	RegenerateApplicantActivationCode(ctx context.Context, uow *uow.UnitOfWork, applicant *pb.Applicant) (*code.Code, error)
	RegenerateEmployerActivationCode(ctx context.Context, uow *uow.UnitOfWork, employer *pb.Employer) (*code.Code, error)
	CheckApplicantResetPasswordCode(ctx context.Context, uow *uow.UnitOfWork, applicant *pb.Applicant, rawCode string) (bool, error)
	CheckEmployerResetPasswordCode(ctx context.Context, uow *uow.UnitOfWork, employer *pb.Employer, rawCode string) (bool, error)
	RegenerateApplicantResetPasswordCode(ctx context.Context, uow *uow.UnitOfWork, applicant *pb.Applicant) (*code.Code, error)
	RegenerateEmployerResetPasswordCode(ctx context.Context, uow *uow.UnitOfWork, employer *pb.Employer) (*code.Code, error)
}

type service struct {
	dataProvider *codeDataProvider
	log          *zap.SugaredLogger
}

func New(redis *redis.RedisClient, log *zap.SugaredLogger) CodeService {
	return &service{
		dataProvider: newCodeDataProvider(redis),
		log:          log,
	}
}

func (s *service) CreateApplicantActivationCode(
	ctx context.Context, uow *uow.UnitOfWork,
	applicant *pb.Applicant,
) (*code.Code, error) {
	l := s.log.With("op", "create_applicant_activation_code", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	code, err := code.New(applicant.Id)
	if err != nil {
		l.Errorw("code.create_activation_code_failed", "err", err)
		return nil, err
	}

	existedCode, err := s.dataProvider.GetApplicantActivationCode(ctx, uow, applicant.Id)
	if err != nil {
		l.Errorw("code.create_activation_code_failed", "err", err)
		return nil, err
	}
	if existedCode != nil {
		code.SetId(existedCode.Id())
	}

	if err := s.dataProvider.SaveApplicantActivationCode(ctx, uow, code); err != nil {
		l.Errorw("code.create_activation_code_failed", "err", err)
		return nil, err
	}

	l.Infow("code.create_activation_code.success")
	return code, nil
}

func (s *service) CreateEmployerActivationCode(
	ctx context.Context, uow *uow.UnitOfWork,
	employer *pb.Employer,
) (*code.Code, error) {
	l := s.log.With("op", "create_employer_activation_code", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	code, err := code.New(employer.Id)
	if err != nil {
		l.Errorw("code.create_activation_code_failed", "err", err)
		return nil, err
	}

	existedCode, err := s.dataProvider.GetEmployerActivationCode(ctx, uow, employer.Id)
	if err != nil {
		l.Errorw("code.create_activation_code_failed", "err", err)
		return nil, err
	}
	if existedCode != nil {
		code.SetId(existedCode.Id())
	}

	if err := s.dataProvider.SaveEmployerActivationCode(ctx, uow, code); err != nil {
		l.Errorw("code.create_activation_code_failed", "err", err)
		return nil, err
	}

	l.Infow("code.create_activation_code.success")
	return code, nil
}

func (s *service) CheckApplicantActivationCode(
	ctx context.Context, uow *uow.UnitOfWork,
	applicant *pb.Applicant, rawCode string,
) (bool, error) {
	l := s.log.With("op", "check_applicant_activation_code", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	code, err := s.dataProvider.GetApplicantActivationCode(ctx, uow, applicant.Id)
	if err != nil {
		l.Errorw("code.check_activation_code_failed", "err", err)
		return false, err
	}
	if code == nil {
		l.Errorw("code.check_activation_code_failed", "err", "activation code not found")
		return false, nil
	}

	valid, err := code.CheckCode(rawCode)
	if err != nil {
		l.Warnw("code.check_activation_code_failed", "err", err)
		return false, err
	}
	if !valid {
		l.Warnw("code.check_activation_code_failed", "err", "invalid code")
		return false, nil
	}

	if err := s.dataProvider.DeleteApplicantActivationCode(ctx, uow, code); err != nil {
		l.Errorw("code.check_activation_code_failed", "err", err)
	}

	l.Infow("code.check_activation_code.success")
	return true, nil
}

func (s *service) CheckEmployerActivationCode(
	ctx context.Context, uow *uow.UnitOfWork,
	employer *pb.Employer, rawCode string,
) (bool, error) {
	l := s.log.With("op", "check_employer_activation_code", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	code, err := s.dataProvider.GetEmployerActivationCode(ctx, uow, employer.Id)
	if err != nil {
		l.Errorw("code.check_activation_code_failed", "err", err)
		return false, err
	}
	if code == nil {
		l.Errorw("code.check_activation_code_failed", "err", "activation code not found")
		return false, nil
	}

	valid, err := code.CheckCode(rawCode)
	if err != nil {
		l.Warnw("code.check_activation_code_failed", "err", err)
		return false, err
	}
	if !valid {
		l.Warnw("code.check_activation_code_failed", "err", "invalid code")
		return false, nil
	}

	if err := s.dataProvider.DeleteEmployerActivationCode(ctx, uow, code); err != nil {
		l.Errorw("code.check_activation_code_failed", "err", err)
	}

	l.Infow("code.check_activation_code.success")
	return true, nil
}

func (s *service) RegenerateApplicantActivationCode(
	ctx context.Context, uow *uow.UnitOfWork,
	applicant *pb.Applicant,
) (*code.Code, error) {
	l := s.log.With("op", "regenerate_applicant_activation_code", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	var c *code.Code
	existedCode, err := s.dataProvider.GetApplicantActivationCode(ctx, uow, applicant.Id)
	if err != nil {
		l.Errorw("code.regenerate_activation_code_failed", "err", err)
		return nil, err
	}
	if existedCode != nil {
		c = existedCode
		err = c.GenerateCode()
		if err != nil {
			l.Warnw("code.regenerate_activation_code_failed", "err", err)
			return nil, err
		}
	} else {
		c, err = code.New(applicant.Id)
		if err != nil {
			l.Errorw("code.regenerate_activation_code_failed", "err", err)
			return nil, err
		}
	}

	if err := s.dataProvider.SaveApplicantActivationCode(ctx, uow, c); err != nil {
		l.Errorw("code.regenerate_activation_code_failed", "err", err)
		return nil, err
	}

	l.Infow("code.regenerate_activation_code.success")
	return c, nil
}

func (s *service) RegenerateEmployerActivationCode(
	ctx context.Context, uow *uow.UnitOfWork,
	employer *pb.Employer,
) (*code.Code, error) {
	l := s.log.With("op", "regenerate_employer_activation_code", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	var c *code.Code
	existedCode, err := s.dataProvider.GetEmployerActivationCode(ctx, uow, employer.Id)
	if err != nil {
		l.Errorw("code.regenerate_activation_code_failed", "err", err)
		return nil, err
	}
	if existedCode != nil {
		c = existedCode
		err = c.GenerateCode()
		if err != nil {
			l.Warnw("code.regenerate_activation_code_failed", "err", err)
			return nil, err
		}
	} else {
		c, err = code.New(employer.Id)
		if err != nil {
			l.Errorw("code.regenerate_activation_code_failed", "err", err)
			return nil, err
		}
	}

	if err := s.dataProvider.SaveEmployerActivationCode(ctx, uow, c); err != nil {
		l.Errorw("code.regenerate_activation_code_failed", "err", err)
		return nil, err
	}

	l.Infow("code.regenerate_activation_code.success")
	return c, nil
}

func (s *service) CheckApplicantResetPasswordCode(
	ctx context.Context, uow *uow.UnitOfWork,
	applicant *pb.Applicant, rawCode string,
) (bool, error) {
	l := s.log.With("op", "check_applicant_reset_password_code", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	code, err := s.dataProvider.GetApplicantResetPasswordCode(ctx, uow, applicant.Id)
	if err != nil {
		l.Errorw("code.check_reset_password_code_failed", "err", err)
		return false, err
	}
	if code == nil {
		l.Errorw("code.check_reset_password_code_failed", "err", "reset password code not found")
		return false, nil
	}

	valid, err := code.CheckCode(rawCode)
	if err != nil {
		l.Warnw("code.check_reset_password_code_failed", "err", err)
		return false, err
	}
	if !valid {
		l.Warnw("code.check_reset_password_code_failed", "err", "invalid code")
		return false, nil
	}

	if err := s.dataProvider.DeleteApplicantResetPasswordCode(ctx, uow, code); err != nil {
		l.Errorw("code.check_reset_password_code_failed", "err", err)
	}

	l.Infow("code.check_reset_password_code.success")
	return true, nil
}

func (s *service) CheckEmployerResetPasswordCode(
	ctx context.Context, uow *uow.UnitOfWork,
	employer *pb.Employer, rawCode string,
) (bool, error) {
	l := s.log.With("op", "check_employer_reset_password_code", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	code, err := s.dataProvider.GetEmployerResetPasswordCode(ctx, uow, employer.Id)
	if err != nil {
		l.Errorw("code.check_reset_password_code_failed", "err", err)
		return false, err
	}
	if code == nil {
		l.Errorw("code.check_reset_password_code_failed", "err", "reset password code not found")
		return false, nil
	}

	valid, err := code.CheckCode(rawCode)
	if err != nil {
		l.Warnw("code.check_reset_password_code_failed", "err", err)
		return false, err
	}
	if !valid {
		l.Warnw("code.check_reset_password_code_failed", "err", "invalid code")
		return false, nil
	}

	if err := s.dataProvider.DeleteEmployerResetPasswordCode(ctx, uow, code); err != nil {
		l.Errorw("code.check_reset_password_code_failed", "err", err)
	}

	l.Infow("code.check_reset_password_code.success")
	return true, nil
}

func (s *service) RegenerateApplicantResetPasswordCode(
	ctx context.Context, uow *uow.UnitOfWork,
	applicant *pb.Applicant,
) (*code.Code, error) {
	l := s.log.With("op", "regenerate_applicant_reset_password_code", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	var c *code.Code
	existedCode, err := s.dataProvider.GetApplicantResetPasswordCode(ctx, uow, applicant.Id)
	if err != nil {
		l.Errorw("code.regenerate_reset_password_code_failed", "err", err)
		return nil, err
	}
	if existedCode != nil {
		c = existedCode
		err = c.GenerateCode()
		if err != nil {
			l.Warnw("code.regenerate_reset_password_code_failed", "err", err)
			return nil, err
		}
	} else {
		c, err = code.New(applicant.Id)
		if err != nil {
			l.Errorw("code.regenerate_reset_password_code_failed", "err", err)
			return nil, err
		}
	}

	if err := s.dataProvider.SaveApplicantResetPasswordCode(ctx, uow, c); err != nil {
		l.Errorw("code.regenerate_reset_password_code_failed", "err", err)
		return nil, err
	}

	l.Infow("code.regenerate_reset_password_code.success")
	return c, nil
}

func (s *service) RegenerateEmployerResetPasswordCode(
	ctx context.Context, uow *uow.UnitOfWork,
	employer *pb.Employer,
) (*code.Code, error) {
	l := s.log.With("op", "regenerate_employer_reset_password_code", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	var c *code.Code
	existedCode, err := s.dataProvider.GetEmployerResetPasswordCode(ctx, uow, employer.Id)
	if err != nil {
		l.Errorw("code.regenerate_reset_password_code_failed", "err", err)
		return nil, err
	}
	if existedCode != nil {
		c = existedCode
		err = c.GenerateCode()
		if err != nil {
			l.Warnw("code.regenerate_reset_password_code_failed", "err", err)
			return nil, err
		}
	} else {
		c, err = code.New(employer.Id)
		if err != nil {
			l.Errorw("code.regenerate_reset_password_code_failed", "err", err)
			return nil, err
		}
	}

	if err := s.dataProvider.SaveEmployerResetPasswordCode(ctx, uow, c); err != nil {
		l.Errorw("code.regenerate_reset_password_code_failed", "err", err)
		return nil, err
	}

	l.Infow("code.regenerate_reset_password_code.success")
	return c, nil
}
