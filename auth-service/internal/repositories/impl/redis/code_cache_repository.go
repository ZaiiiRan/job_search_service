package redisimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/ZaiiiRan/job_search_service/auth-service/internal/domain/code"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/impl"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/models"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/transport/redis"
)

const (
	ApplicantActivationCodesCache    impl.RepositoryType = "code:activation:applicant"
	EmployerActivationCodesCache     impl.RepositoryType = "code:activation:employer"
	ApplicantResetPasswordCodesCache impl.RepositoryType = "code:reset:applicant"
	EmployerResetPasswordCodesCache  impl.RepositoryType = "code:reset:employer"
)

type CodeCacheRepository struct {
	redis          *redis.RedisClient
	repositoryType impl.RepositoryType
}

func NewCodeCacheRepository(redis *redis.RedisClient, repositoryType impl.RepositoryType) interfaces.CodeCacheRepository {
	return &CodeCacheRepository{
		redis:          redis,
		repositoryType: repositoryType,
	}
}

func (r *CodeCacheRepository) GetById(ctx context.Context, id int64) (*code.Code, error) {
	dal, err := get[models.V1CodeDal](ctx, r.redis, r.keyById(id))
	if err != nil {
		return nil, err
	}
	if dal == nil {
		return nil, nil
	}
	return dal.ToDomain(), nil
}

func (r *CodeCacheRepository) SetById(ctx context.Context, code *code.Code) error {
	dal := models.V1CodeDalFromDomain(code)
	return set(ctx, r.redis, r.keyById(dal.Id), dal, time.Until(dal.ExpiresAt))
}

func (r *CodeCacheRepository) DelById(ctx context.Context, id int64) error {
	return del(ctx, r.redis, r.keyById(id))
}

func (r *CodeCacheRepository) GetByUserId(ctx context.Context, userId int64) (*code.Code, error) {
	dal, err := get[models.V1CodeDal](ctx, r.redis, r.keyByUserId(userId))
	if err != nil {
		return nil, err
	}
	if dal == nil {
		return nil, nil
	}
	return dal.ToDomain(), nil
}

func (r *CodeCacheRepository) SetByUserId(ctx context.Context, code *code.Code) error {
	dal := models.V1CodeDalFromDomain(code)
	return set(ctx, r.redis, r.keyByUserId(dal.UserId), dal, time.Until(dal.ExpiresAt))
}

func (r *CodeCacheRepository) DelByUserId(ctx context.Context, userId int64) error {
	return del(ctx, r.redis, r.keyByUserId(userId))
}

func (r *CodeCacheRepository) keyById(id int64) string {
	return fmt.Sprintf("%s:id:%d", r.repositoryType, id)
}

func (r *CodeCacheRepository) keyByUserId(userId int64) string {
	return fmt.Sprintf("%s:user_id:%d", r.repositoryType, userId)
}
