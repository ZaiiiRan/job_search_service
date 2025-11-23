package redisimpl

import (
	"context"
	"fmt"
	"time"

	userversion "github.com/ZaiiiRan/job_search_service/auth-service/internal/domain/user_version"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/impl"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/models"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/transport/redis"
)

const (
	ApplicantVersionCache impl.RepositoryType = "user_version:applicant"
	EmployerVersionCache  impl.RepositoryType = "user_version:employer"
	userVersionTTL                            = time.Minute * 15
)

type UserVersionCacheRepository struct {
	redis          *redis.RedisClient
	repositoryType impl.RepositoryType
}

func NewUserVersionCacheRepository(redis *redis.RedisClient, repositoryType impl.RepositoryType) interfaces.UserVersionCacheRepository {
	return &UserVersionCacheRepository{
		redis:          redis,
		repositoryType: repositoryType,
	}
}

func (r *UserVersionCacheRepository) SetByUserId(ctx context.Context, uv *userversion.UserVersion) error {
	dal := models.V1UserVersionDalFromDomain(uv)
	return set(ctx, r.redis, r.keyByUserId(dal.UserId), dal, userVersionTTL)
}

func (r *UserVersionCacheRepository) GetByUserId(ctx context.Context, userId int64) (*userversion.UserVersion, error) {
	dal, err := get[models.V1UserVersionDal](ctx, r.redis, r.keyByUserId(userId))
	if err != nil {
		return nil, err
	}
	if dal == nil {
		return nil, nil
	}
	return dal.ToDomain(), nil
}

func (r *UserVersionCacheRepository) DelByUserId(ctx context.Context, userId int64) error {
	return del(ctx, r.redis, r.keyByUserId(userId))
}

func (r *UserVersionCacheRepository) keyByUserId(userId int64) string {
	return fmt.Sprintf("%s:user_id:%d", r.repositoryType, userId)
}
