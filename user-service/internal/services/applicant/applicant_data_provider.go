package applicantservice

import (
	"context"

	"github.com/ZaiiiRan/job_search_service/user-service/internal/domain/user/applicant"
	repo "github.com/ZaiiiRan/job_search_service/user-service/internal/repositories/impl/postgres"
	cache "github.com/ZaiiiRan/job_search_service/user-service/internal/repositories/impl/redis"
	dal "github.com/ZaiiiRan/job_search_service/user-service/internal/repositories/models"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/transport/postgres"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/transport/redis"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/utils"
)

type applicantDataProvider struct {
	pg    *postgres.PostgresClient
	redis *redis.RedisClient
}

func newApplicantDataProvider(pg *postgres.PostgresClient, redis *redis.RedisClient) *applicantDataProvider {
	return &applicantDataProvider{pg: pg, redis: redis}
}

func (p *applicantDataProvider) GetByEmail(ctx context.Context, email string) (*applicant.Applicant, error) {
	cacheRepo := cache.NewApplicantCacheRepository(p.redis)
	a, err := cacheRepo.GetApplicantByEmail(ctx, email)
	if err == nil && a != nil {
		return a, nil
	}

	pgConn, err := p.pg.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	defer pgConn.Release()

	dbRepo := repo.NewApplicantRepository(pgConn)
	query := dal.NewQueryApplicantsDal(nil, []string{email}, nil, nil, utils.BoolPtr(false), nil, nil, nil, nil, 1, 1)
	list, err := dbRepo.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	cacheRepo.SetApplicant(ctx, list[0])
	cacheRepo.SetApplicantByEmail(ctx, list[0])

	return list[0], err
}

func (p *applicantDataProvider) GetById(ctx context.Context, id int64) (*applicant.Applicant, error) {
	cacheRepo := cache.NewApplicantCacheRepository(p.redis)
	a, err := cacheRepo.GetApplicant(ctx, id)
	if err == nil && a != nil {
		return a, nil
	}

	pgConn, err := p.pg.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	defer pgConn.Release()

	dbRepo := repo.NewApplicantRepository(pgConn)
	query := dal.NewQueryApplicantsDal([]int64{id}, nil, nil, nil, nil, nil, nil, nil, nil, 1, 1)
	list, err := dbRepo.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}

	cacheRepo.SetApplicant(ctx, list[0])
	return list[0], nil
}

func (p *applicantDataProvider) Save(ctx context.Context, a *applicant.Applicant) error {
	pgConn, err := p.pg.GetConn(ctx)
	if err != nil {
		return err
	}
	defer pgConn.Release()

	dbRepo := repo.NewApplicantRepository(pgConn)

	if a.Id() == 0 {
		if err := dbRepo.Create(ctx, a); err != nil {
			return err
		}
	} else {
		if err := dbRepo.Update(ctx, a); err != nil {
			return err
		}
	}

	cacheRepo := cache.NewApplicantCacheRepository(p.redis)
	cacheRepo.InvalidateApplicantList(ctx)
	cacheRepo.SetApplicant(ctx, a)
	cacheRepo.SetApplicantByEmail(ctx, a)

	return nil
}

func (p *applicantDataProvider) QueryList(ctx context.Context, query *dal.QueryApplicantsDal) ([]*applicant.Applicant, error) {
	cacheRepo := cache.NewApplicantCacheRepository(p.redis)
	list, err := cacheRepo.GetApplicantList(ctx, query)
	if err == nil && len(list) > 0 {
		return list, nil
	}

	pgConn, err := p.pg.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	defer pgConn.Release()

	dbRepo := repo.NewApplicantRepository(pgConn)
	list, err = dbRepo.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		cacheRepo.SetApplicantList(ctx, query, list)
	}

	return list, err
}
