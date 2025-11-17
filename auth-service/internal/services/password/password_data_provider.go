package passwordservice

import (
	"context"

	"github.com/ZaiiiRan/job_search_service/auth-service/internal/domain/password"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/impl"
	repo "github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/impl/postgres"
	cache "github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/impl/redis"
	dal "github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/models"
	uow "github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/transport/redis"
)

type passwordDataProvider struct {
	redis *redis.RedisClient
}

func newPasswordDataProvider(redis *redis.RedisClient) *passwordDataProvider {
	return &passwordDataProvider{
		redis: redis,
	}
}

func (p *passwordDataProvider) SaveApplicantPassword(ctx context.Context, uow *uow.UnitOfWork, password *password.Password) error {
	return p.save(ctx, uow, password, repo.ApplicantPasswordRepository, cache.ApplicantPasswordCache)
}

func (p *passwordDataProvider) SaveEmployerPassword(ctx context.Context, uow *uow.UnitOfWork, password *password.Password) error {
	return p.save(ctx, uow, password, repo.EmployerPasswordRepository, cache.EmployerPasswordCache)
}

func (p *passwordDataProvider) GetApplicantPasswordByUserId(ctx context.Context, uow *uow.UnitOfWork, userId int64) (*password.Password, error) {
	return p.get(ctx, uow, userId, repo.ApplicantPasswordRepository, cache.ApplicantPasswordCache)
}

func (p *passwordDataProvider) GetEmployerPasswordByUserId(ctx context.Context, uow *uow.UnitOfWork, userId int64) (*password.Password, error) {
	return p.get(ctx, uow, userId, repo.EmployerPasswordRepository, cache.EmployerPasswordCache)
}

func (p *passwordDataProvider) get(
	ctx context.Context, uow *uow.UnitOfWork,
	userId int64,
	repoType impl.RepositoryType, cacheType impl.RepositoryType,
) (*password.Password, error) {
	cacheRepo := cache.NewPasswordCacheRepository(p.redis, cacheType)
	password, err := cacheRepo.GetByUserId(ctx, userId)
	if err == nil && password != nil {
		return password, nil
	}

	dbRepo := repo.NewPasswordRepository(uow, repoType)
	query := dal.NewQueryPasswordDal(nil, &userId)
	password, err = dbRepo.QueryPassword(ctx, query)
	if err != nil {
		return nil, err
	}
	if password == nil {
		return nil, nil
	}

	cacheRepo.SetByUserId(ctx, password)
	return password, nil
}

func (p *passwordDataProvider) save(
	ctx context.Context, uow *uow.UnitOfWork,
	password *password.Password,
	repoType impl.RepositoryType, cacheType impl.RepositoryType,
) error {
	dbRepo := repo.NewPasswordRepository(uow, repoType)

	if password.Id() == 0 {
		if err := dbRepo.CreatePassword(ctx, password); err != nil {
			return err
		}
	} else {
		if err := dbRepo.UpdatePassword(ctx, password); err != nil {
			return err
		}
	}

	cacheRepo := cache.NewPasswordCacheRepository(p.redis, cacheType)
	if err := cacheRepo.SetByUserId(ctx, password); err != nil {
		return err
	}

	return nil
}
