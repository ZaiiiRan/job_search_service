package redisimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/ZaiiiRan/job_search_service/auth-service/internal/domain/token"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/impl"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/models"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/transport/redis"
)

const (
	ApplicantRefreshTokenCache impl.RepositoryType = "refresh:applicant"
	EmployerRefreshTokenCache  impl.RepositoryType = "refresh:employer"
)

type TokenCacheRepository struct {
	redis          *redis.RedisClient
	repositoryType impl.RepositoryType
}

func NewTokenCacheRepository(redis *redis.RedisClient, repositoryType impl.RepositoryType) interfaces.TokenCacheRepository {
	return &TokenCacheRepository{
		redis:          redis,
		repositoryType: repositoryType,
	}
}

func (r *TokenCacheRepository) Get(ctx context.Context, userId int64, token string) (*token.Token, error) {
	dal, err := get[models.V1RefreshTokenDal](ctx, r.redis, r.keyToken(userId, token))
	if err != nil {
		return nil, err
	}
	if dal == nil {
		return nil, nil
	}
	return dal.ToDomain(), nil
}

func (r *TokenCacheRepository) Set(ctx context.Context, token *token.Token) error {
	dal := models.V1RefreshTokenDalFromDomain(token)
	return set(ctx, r.redis, r.keyToken(dal.UserId, dal.Token), dal, time.Until(dal.ExpiresAt))
}

func (r *TokenCacheRepository) Del(ctx context.Context, userId int64, token string) error {
	return del(ctx, r.redis, r.keyToken(userId, token))
}

func (r *TokenCacheRepository) DelByUserId(ctx context.Context, userId int64) error {
	return invalidateByPrefix(ctx, r.redis, r.keyUserId(userId))
}

func (r *TokenCacheRepository) keyToken(userId int64, token string) string {
	return fmt.Sprintf("%s:%d:%s", r.repositoryType, userId, token)
}

func (r *TokenCacheRepository) keyUserId(userId int64) string {
	return fmt.Sprintf("%s:%d", r.repositoryType, userId)
}
