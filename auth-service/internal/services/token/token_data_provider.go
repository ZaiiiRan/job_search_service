package tokenservice

import (
	"context"

	"github.com/ZaiiiRan/job_search_service/auth-service/internal/domain/token"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/impl"
	repo "github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/impl/postgres"
	cache "github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/impl/redis"
	dal "github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/models"
	postgresunitofwork "github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/transport/redis"
)

type tokenDataProvider struct {
	redis *redis.RedisClient
}

func newTokenDataProvider(redis *redis.RedisClient) *tokenDataProvider {
	return &tokenDataProvider{
		redis: redis,
	}
}

func (p *tokenDataProvider) SaveApplicantToken(ctx context.Context, uow *postgresunitofwork.UnitOfWork, t *token.Token) error {
	return p.save(ctx, uow, t, repo.ApplicantRefreshTokenRepository, cache.ApplicantRefreshTokenCache)
}

func (p *tokenDataProvider) SaveEmployerToken(ctx context.Context, uow *postgresunitofwork.UnitOfWork, t *token.Token) error {
	return p.save(ctx, uow, t, repo.EmployerRefreshTokenRepository, cache.EmployerRefreshTokenCache)
}

func (p *tokenDataProvider) GetApplicantToken(
	ctx context.Context, uow *postgresunitofwork.UnitOfWork,
	token string,
) (*token.Token, error) {
	return p.get(ctx, uow, token, repo.ApplicantRefreshTokenRepository, cache.ApplicantRefreshTokenCache)
}

func (p *tokenDataProvider) GetEmployerToken(
	ctx context.Context, uow *postgresunitofwork.UnitOfWork,
	token string,
) (*token.Token, error) {
	return p.get(ctx, uow, token, repo.EmployerRefreshTokenRepository, cache.EmployerRefreshTokenCache)
}

func (p *tokenDataProvider) DeleteApplicantToken(
	ctx context.Context, uow *postgresunitofwork.UnitOfWork,
	token string,
) error {
	return p.delete(ctx, uow, token, repo.ApplicantRefreshTokenRepository, cache.ApplicantRefreshTokenCache)
}

func (p *tokenDataProvider) DeleteApplicantTokenFromCache(ctx context.Context, token string) error {
	cacheRepo := cache.NewTokenCacheRepository(p.redis, cache.ApplicantRefreshTokenCache)
	return cacheRepo.Del(ctx, token)
}

func (p *tokenDataProvider) DeleteEmployerToken(
	ctx context.Context, uow *postgresunitofwork.UnitOfWork,
	token string,
) error {
	return p.delete(ctx, uow, token, repo.EmployerRefreshTokenRepository, cache.EmployerRefreshTokenCache)
}

func (p *tokenDataProvider) DeleteEmployerTokenFromCache(ctx context.Context, token string) error {
	cacheRepo := cache.NewTokenCacheRepository(p.redis, cache.EmployerRefreshTokenCache)
	return cacheRepo.Del(ctx, token)
}

func (p *tokenDataProvider) get(
	ctx context.Context, uow *postgresunitofwork.UnitOfWork,
	token string,
	repoType impl.RepositoryType, cacheType impl.RepositoryType,
) (*token.Token, error) {
	cacheRepo := cache.NewTokenCacheRepository(p.redis, cacheType)
	t, err := cacheRepo.Get(ctx, token)
	if err == nil && t != nil {
		return t, nil
	}

	dbRepo := repo.NewTokenRepository(uow, repoType)
	query := dal.NewQueryTokenDal(nil, nil, &token)
	t, err = dbRepo.QueryToken(ctx, query)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, nil
	}

	cacheRepo.Set(ctx, t)
	return t, nil
}

func (p *tokenDataProvider) delete(
	ctx context.Context, uow *postgresunitofwork.UnitOfWork,
	token string,
	repoType impl.RepositoryType, cacheType impl.RepositoryType,
) error {
	dbRepo := repo.NewTokenRepository(uow, repoType)
	err := dbRepo.DeleteToken(ctx, token)
	if err != nil {
		return err
	}

	cacheRepo := cache.NewTokenCacheRepository(p.redis, cacheType)
	cacheRepo.Del(ctx, token)
	return nil
}

func (p *tokenDataProvider) save(
	ctx context.Context, uow *postgresunitofwork.UnitOfWork,
	t *token.Token,
	repoType impl.RepositoryType, cacheType impl.RepositoryType,
) error {
	dbRepo := repo.NewTokenRepository(uow, repoType)

	if t.Id() == 0 {
		if err := dbRepo.CreateToken(ctx, t); err != nil {
			return err
		}
	} else {
		if err := dbRepo.UpdateToken(ctx, t); err != nil {
			return err
		}
	}

	cacheRepo := cache.NewTokenCacheRepository(p.redis, cacheType)
	cacheRepo.Set(ctx, t)

	return nil
}
