package redisimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/ZaiiiRan/job_search_service/user-service/internal/domain/user/employer"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/repositories/models"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/transport/redis"
)

const (
	employerKeyPrefix  = "employer"
	employerListPrefix = "employer:list"
	employerTTL        = 10 * time.Minute
	employerListTTL    = 5 * time.Minute
)

type EmployerCacheRepository struct {
	redis *redis.RedisClient
}

func NewEmployerCacheRepository(redis *redis.RedisClient) interfaces.EmployerCacheRepository {
	return &EmployerCacheRepository{
		redis: redis,
	}
}

func (r *EmployerCacheRepository) SetEmployer(ctx context.Context, emp *employer.Employer) error {
	dal := models.V1EmployerDalFromDomain(emp)
	return set(ctx, r.redis, r.keyById(dal.Id), dal, employerTTL)
}

func (r *EmployerCacheRepository) GetEmployer(ctx context.Context, id int64) (*employer.Employer, error) {
	dal, err := get[models.V1EmployerDal](ctx, r.redis, r.keyById(id))
	if err != nil {
		return nil, err
	}
	if dal == nil {
		return nil, nil
	}
	return dal.ToDomain(), nil
}

func (r *EmployerCacheRepository) DeleteEmployer(ctx context.Context, id int64) error {
	return del(ctx, r.redis, r.keyById(id))
}

func (r *EmployerCacheRepository) SetEmployerList(ctx context.Context, query *models.QueryEmployersDal, employers []*employer.Employer) error {
	var dalList []models.V1EmployerDal
	for _, emp := range employers {
		dalList = append(dalList, models.V1EmployerDalFromDomain(emp))
	}

	key, err := r.keyByQuery(query)
	if err != nil {
		return err
	}
	return set(ctx, r.redis, key, dalList, employerListTTL)
}

func (r *EmployerCacheRepository) GetEmployerList(ctx context.Context, query *models.QueryEmployersDal) ([]*employer.Employer, error) {
	key, err := r.keyByQuery(query)
	if err != nil {
		return nil, err
	}

	val, err := get[[]models.V1EmployerDal](ctx, r.redis, key)
	if err != nil || val == nil {
		return nil, err
	}

	var res []*employer.Employer
	for _, dal := range *val {
		res = append(res, dal.ToDomain())
	}
	return res, nil
}

func (r *EmployerCacheRepository) InvalidateEmployerList(ctx context.Context) error {
	return invalidateByPrefix(ctx, r.redis, employerListPrefix)
}

func (r *EmployerCacheRepository) keyById(id int64) string {
	return fmt.Sprintf("%s:%d", employerKeyPrefix, id)
}

func (r *EmployerCacheRepository) keyByQuery(query *models.QueryEmployersDal) (string, error) {
	h, err := queryHash(query)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:query:%s", employerListPrefix, h), nil
}
