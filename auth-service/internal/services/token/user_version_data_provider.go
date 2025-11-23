package tokenservice

import (
	"context"

	userversion "github.com/ZaiiiRan/job_search_service/auth-service/internal/domain/user_version"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/impl"
	repo "github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/impl/postgres"
	cache "github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/impl/redis"
	dal "github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/models"
	postgresunitofwork "github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/transport/redis"
)

type userVersionDataProvider struct {
	redis *redis.RedisClient
}

func newUserVersionDataProvider(redis *redis.RedisClient) *userVersionDataProvider {
	return &userVersionDataProvider{
		redis: redis,
	}
}

func (p *userVersionDataProvider) GetApplicantVersion(ctx context.Context, uow *postgresunitofwork.UnitOfWork, userId int64) (*userversion.UserVersion, error) {
	return p.get(ctx, uow, userId, repo.ApplicantVersion, cache.ApplicantVersionCache)
}

func (p *userVersionDataProvider) GetEmployerVersion(ctx context.Context, uow *postgresunitofwork.UnitOfWork, userId int64) (*userversion.UserVersion, error) {
	return p.get(ctx, uow, userId, repo.EmployerVersion, cache.EmployerVersionCache)
}

func (p *userVersionDataProvider) SaveApplicantVersion(ctx context.Context, uow *postgresunitofwork.UnitOfWork, uv *userversion.UserVersion) error {
	return p.save(ctx, uow, uv, repo.ApplicantVersion, cache.ApplicantVersionCache)
}

func (p *userVersionDataProvider) SaveEmployerVersion(ctx context.Context, uow *postgresunitofwork.UnitOfWork, uv *userversion.UserVersion) error {
	return p.save(ctx, uow, uv, repo.EmployerVersion, cache.EmployerVersionCache)
}

func (p *userVersionDataProvider) DeleteApplicantVersion(ctx context.Context, uow *postgresunitofwork.UnitOfWork, uv *userversion.UserVersion) error {
	return p.delete(ctx, uow, uv, repo.ApplicantVersion, cache.ApplicantVersionCache)
}

func (p *userVersionDataProvider) DeleteEmployerVersion(ctx context.Context, uow *postgresunitofwork.UnitOfWork, uv *userversion.UserVersion) error {
	return p.delete(ctx, uow, uv, repo.EmployerVersion, cache.EmployerVersionCache)
}

func (p *userVersionDataProvider) get(
	ctx context.Context, uow *postgresunitofwork.UnitOfWork,
	userId int64,
	repoType impl.RepositoryType, cacheType impl.RepositoryType,
) (*userversion.UserVersion, error) {
	cacheRepo := cache.NewUserVersionCacheRepository(p.redis, cacheType)
	uv, err := cacheRepo.GetByUserId(ctx, userId)
	if err == nil && uv != nil {
		return uv, nil
	}

	dbRepo := repo.NewUserVersionRepository(uow, repoType)
	query := dal.NewQueryUserVersionDal(nil, &userId)
	uv, err = dbRepo.QueryUserVersion(ctx, query)
	if err != nil {
		return nil, err
	}
	if uv == nil {
		return nil, nil
	}

	cacheRepo.SetByUserId(ctx, uv)
	return uv, nil
}

func (p *userVersionDataProvider) save(
	ctx context.Context, uow *postgresunitofwork.UnitOfWork,
	uv *userversion.UserVersion,
	repoType impl.RepositoryType, cacheType impl.RepositoryType,
) error {
	dbRepo := repo.NewUserVersionRepository(uow, repoType)

	if uv.Id() == 0 {
		if err := dbRepo.CreateUserVersion(ctx, uv); err != nil {
			return err
		}
	} else {
		if err := dbRepo.UpdateUserVersion(ctx, uv); err != nil {
			return err
		}
	}

	cacheRepo := cache.NewUserVersionCacheRepository(p.redis, cacheType)
	if err := cacheRepo.SetByUserId(ctx, uv); err != nil {
		return err
	}

	return nil
}

func (p *userVersionDataProvider) delete(
	ctx context.Context, uow *postgresunitofwork.UnitOfWork,
	uv *userversion.UserVersion,
	repoType impl.RepositoryType, cacheType impl.RepositoryType,
) error {
	cacheRepo := cache.NewUserVersionCacheRepository(p.redis, cacheType)
	if err := cacheRepo.DelByUserId(ctx, uv.UserId()); err != nil {
		return err
	}

	dbRepo := repo.NewUserVersionRepository(uow, repoType)
	if err := dbRepo.DeleteUserVersion(ctx, uv); err != nil {
		return err
	}

	return nil
}
