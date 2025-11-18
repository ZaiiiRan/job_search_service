package codeservice

import (
	"context"

	"github.com/ZaiiiRan/job_search_service/auth-service/internal/domain/code"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/impl"
	repo "github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/impl/postgres"
	cache "github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/impl/redis"
	dal "github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/models"
	uow "github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/transport/redis"
)

type codeDataProvider struct {
	redis *redis.RedisClient
}

func newCodeDataProvider(redis *redis.RedisClient) *codeDataProvider {
	return &codeDataProvider{
		redis: redis,
	}
}

func (p *codeDataProvider) SaveApplicantActivationCode(ctx context.Context, uow *uow.UnitOfWork, code *code.Code) error {
	return p.save(ctx, uow, code, repo.ApplicantActivationCodes, cache.ApplicantActivationCodesCache)
}

func (p *codeDataProvider) SaveApplicantResetPasswordCode(ctx context.Context, uow *uow.UnitOfWork, code *code.Code) error {
	return p.save(ctx, uow, code, repo.ApplicantResetPasswordCodes, cache.ApplicantResetPasswordCodesCache)
}

func (p *codeDataProvider) GetApplicantActivationCode(ctx context.Context, uow *uow.UnitOfWork, userId int64) (*code.Code, error) {
	return p.get(ctx, uow, userId, repo.ApplicantActivationCodes, cache.ApplicantActivationCodesCache)
}

func (p *codeDataProvider) GetApplicantResetPasswordCode(ctx context.Context, uow *uow.UnitOfWork, userId int64) (*code.Code, error) {
	return p.get(ctx, uow, userId, repo.ApplicantResetPasswordCodes, cache.ApplicantResetPasswordCodesCache)
}

func (p *codeDataProvider) DeleteApplicantActivationCode(ctx context.Context, uow *uow.UnitOfWork, code *code.Code) error {
	return p.delete(ctx, uow, code, repo.ApplicantActivationCodes, cache.ApplicantActivationCodesCache)
}

func (p *codeDataProvider) DeleteApplicantResetPasswordCode(ctx context.Context, uow *uow.UnitOfWork, code *code.Code) error {
	return p.delete(ctx, uow, code, repo.ApplicantResetPasswordCodes, cache.ApplicantResetPasswordCodesCache)
}

func (p *codeDataProvider) SaveEmployerActivationCode(ctx context.Context, uow *uow.UnitOfWork, code *code.Code) error {
	return p.save(ctx, uow, code, repo.EmployerActivationCodes, cache.EmployerActivationCodesCache)
}

func (p *codeDataProvider) SaveEmployerResetPasswordCode(ctx context.Context, uow *uow.UnitOfWork, code *code.Code) error {
	return p.save(ctx, uow, code, repo.EmployerResetPasswordCodes, cache.EmployerResetPasswordCodesCache)
}

func (p *codeDataProvider) GetEmployerActivationCode(ctx context.Context, uow *uow.UnitOfWork, userId int64) (*code.Code, error) {
	return p.get(ctx, uow, userId, repo.EmployerActivationCodes, cache.EmployerActivationCodesCache)
}

func (p *codeDataProvider) GetEmployerResetPasswordCode(ctx context.Context, uow *uow.UnitOfWork, userId int64) (*code.Code, error) {
	return p.get(ctx, uow, userId, repo.EmployerResetPasswordCodes, cache.EmployerResetPasswordCodesCache)
}

func (p *codeDataProvider) DeleteEmployerActivationCode(ctx context.Context, uow *uow.UnitOfWork, code *code.Code) error {
	return p.delete(ctx, uow, code, repo.EmployerActivationCodes, cache.EmployerActivationCodesCache)
}

func (p *codeDataProvider) DeleteEmployerResetPasswordCode(ctx context.Context, uow *uow.UnitOfWork, code *code.Code) error {
	return p.delete(ctx, uow, code, repo.EmployerResetPasswordCodes, cache.EmployerResetPasswordCodesCache)
}

func (p *codeDataProvider) get(
	ctx context.Context, uow *uow.UnitOfWork,
	userId int64,
	repoType impl.RepositoryType, cacheType impl.RepositoryType,
) (*code.Code, error) {
	cacheRepo := cache.NewCodeCacheRepository(p.redis, cacheType)
	code, err := cacheRepo.GetByUserId(ctx, userId)
	if err == nil && code != nil {
		return code, nil
	}

	dbRepo := repo.NewCodeRepository(uow, repoType)
	query := dal.NewQueryCodeDal(nil, &userId)
	code, err = dbRepo.QueryCode(ctx, query)
	if err != nil {
		return nil, err
	}
	if code == nil {
		return nil, nil
	}

	cacheRepo.SetByUserId(ctx, code)
	return code, nil
}

func (p *codeDataProvider) save(
	ctx context.Context, uow *uow.UnitOfWork,
	code *code.Code,
	repoType impl.RepositoryType, cacheType impl.RepositoryType,
) error {
	dbRepo := repo.NewCodeRepository(uow, repoType)

	if code.Id() == 0 {
		if err := dbRepo.CreateCode(ctx, code); err != nil {
			return err
		}
	} else {
		if err := dbRepo.UpdateCode(ctx, code); err != nil {
			return err
		}
	}

	cacheRepo := cache.NewCodeCacheRepository(p.redis, cacheType)
	if err := cacheRepo.SetByUserId(ctx, code); err != nil {
		return err
	}

	return nil
}

func (p *codeDataProvider) delete(
	ctx context.Context, uow *uow.UnitOfWork,
	code *code.Code,
	repoType impl.RepositoryType, cacheType impl.RepositoryType,
) error {
	cacheRepo := cache.NewCodeCacheRepository(p.redis, cacheType)
	if err := cacheRepo.DelByUserId(ctx, code.UserId()); err != nil {
		return err
	}

	dbRepo := repo.NewCodeRepository(uow, repoType)
	if err := dbRepo.DeleteCode(ctx, code); err != nil {
		return err
	}

	return nil
}
