package redisimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/ZaiiiRan/job_search_service/auth-service/internal/domain/password"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/impl"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/models"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/transport/redis"
)

const (
	ApplicantPasswordCacheRepository impl.RepositoryType = "password:applicant"
	EmployerPasswordCacheRepository  impl.RepositoryType = "password:employer"
)

const (
	passwordTTL       = 5 * time.Minute
)

type PasswordCacheRepository struct {
	redis          *redis.RedisClient
	repositoryType impl.RepositoryType
}

func NewPasswordCacheRepository(redis *redis.RedisClient, repositoryType impl.RepositoryType) interfaces.PasswordCacheRepository {
	return &PasswordCacheRepository{
		redis:          redis,
		repositoryType: repositoryType,
	}
}

func (r *PasswordCacheRepository) GetById(ctx context.Context, id int64) (*password.Password, error) {
	dal, err := get[models.V1UserPasswordDal](ctx, r.redis, r.keyById(id))
	if err != nil {
		return nil, err
	}
	if dal == nil {
		return nil, nil
	}
	return dal.ToDomain(), nil
}

func (r *PasswordCacheRepository) SetById(ctx context.Context, password *password.Password) error {
	dal := models.V1UserPasswordDalFromDomain(password)
	return set(ctx, r.redis, r.keyById(dal.Id), dal, passwordTTL)
}

func (r *PasswordCacheRepository) DelById(ctx context.Context, id int64) error {
	return del(ctx, r.redis, r.keyById(id))
}

func (r *PasswordCacheRepository) GetByUserId(ctx context.Context, userId int64) (*password.Password, error) {
	dal, err := get[models.V1UserPasswordDal](ctx, r.redis, r.keyByUserId(userId))
	if err != nil {
		return nil, err
	}
	if dal == nil {
		return nil, nil
	}
	return dal.ToDomain(), nil
}

func (r *PasswordCacheRepository) SetByUserId(ctx context.Context, password *password.Password) error {
	dal := models.V1UserPasswordDalFromDomain(password)
	return set(ctx, r.redis, r.keyByUserId(dal.UserId), dal, passwordTTL)
}

func (r *PasswordCacheRepository) DelByUserId(ctx context.Context, userId int64) error {
	return del(ctx, r.redis, r.keyByUserId(userId))
}

func (r *PasswordCacheRepository) keyById(id int64) string {
	return fmt.Sprintf("%s:%d", r.repositoryType, id)
}

func (r *PasswordCacheRepository) keyByUserId(userId int64) string {
	return fmt.Sprintf("%s:user_id:%d", r.repositoryType, userId)
}
