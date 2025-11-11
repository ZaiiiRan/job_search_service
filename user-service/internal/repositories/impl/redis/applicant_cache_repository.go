package redisimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/ZaiiiRan/job_search_service/user-service/internal/domain/user/applicant"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/repositories/models"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/transport/redis"
)

const (
	applicantKeyPrefix  = "applicant"
	applicantListPrefix = "applicant:list"
	applicantTTL        = 10 * time.Minute
	applicantListTTL    = 5 * time.Minute
)

type ApplicantCacheRepository struct {
	redis *redis.RedisClient
}

func NewApplicantCacheRepository(redis *redis.RedisClient) interfaces.ApplicantCacheRepository {
	return &ApplicantCacheRepository{
		redis: redis,
	}
}

func (r *ApplicantCacheRepository) SetApplicant(ctx context.Context, applicant *applicant.Applicant) error {
	dal := models.V1ApplicantDalFromDomain(applicant)
	return set(ctx, r.redis, r.keyById(dal.Id), dal, applicantTTL)
}

func (r *ApplicantCacheRepository) GetApplicant(ctx context.Context, id int64) (*applicant.Applicant, error) {
	dal, err := get[models.V1ApplicantDal](ctx, r.redis, r.keyById(id))
	if err != nil {
		return nil, err
	}
	if dal == nil {
		return nil, nil
	}
	return dal.ToDomain(), nil
}

func (r *ApplicantCacheRepository) DeleteApplicant(ctx context.Context, id int64) error {
	return del(ctx, r.redis, r.keyById(id))
}

func (r *ApplicantCacheRepository) SetApplicantList(ctx context.Context, query *models.QueryApplicantsDal, applicants []*applicant.Applicant) error {
	var dal []models.V1ApplicantDal
	for _, applicant := range applicants {
		dal = append(dal, models.V1ApplicantDalFromDomain(applicant))
	}

	key, err := r.keyByQuery(query)
	if err != nil {
		return err
	}
	return set(ctx, r.redis, key, dal, applicantListTTL)
}

func (r *ApplicantCacheRepository) GetApplicantList(ctx context.Context, query *models.QueryApplicantsDal) ([]*applicant.Applicant, error) {
	key, err := r.keyByQuery(query)
	if err != nil {
		return nil, err
	}
	val, err := get[[]models.V1ApplicantDal](ctx, r.redis, key)
	if err != nil || val == nil {
		return nil, err
	}

	var res []*applicant.Applicant
	for _, dal := range *val {
		res = append(res, dal.ToDomain())
	}
	return res, nil
}

func (r *ApplicantCacheRepository) InvalidateApplicantList(ctx context.Context) error {
	return invalidateByPrefix(ctx, r.redis, applicantListPrefix)
}

func (r *ApplicantCacheRepository) keyById(id int64) string {
	return fmt.Sprintf("%s:%d", applicantKeyPrefix, id)
}

func (r *ApplicantCacheRepository) keyByQuery(query *models.QueryApplicantsDal) (string, error) {
	h, err := queryHash(query)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:query:%s", applicantListPrefix, h), nil
}
