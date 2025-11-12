package employerservice

import (
	"context"

	"github.com/ZaiiiRan/job_search_service/user-service/internal/domain/user/employer"
	repo "github.com/ZaiiiRan/job_search_service/user-service/internal/repositories/impl/postgres"
	cache "github.com/ZaiiiRan/job_search_service/user-service/internal/repositories/impl/redis"
	dal "github.com/ZaiiiRan/job_search_service/user-service/internal/repositories/models"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/transport/postgres"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/transport/redis"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/utils"
)

type employerDataProvider struct {
	pg    *postgres.PostgresClient
	redis *redis.RedisClient
}

func newEmployerDataProvider(pg *postgres.PostgresClient, redis *redis.RedisClient) *employerDataProvider {
	return &employerDataProvider{
		pg:    pg,
		redis: redis,
	}
}

func (p *employerDataProvider) GetByEmail(ctx context.Context, email string) (*employer.Employer, error) {
	cacheRepo := cache.NewEmployerCacheRepository(p.redis)
	e, err := cacheRepo.GetEmployerByEmail(ctx, email)
	if err == nil && e != nil {
		return e, nil
	}

	pgConn, err := p.pg.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	defer pgConn.Release()

	dbRepo := repo.NewEmployerRepository(pgConn)
	query := dal.NewQueryEmployersDal(nil, []string{email}, nil, nil, nil, nil, utils.BoolPtr(false), nil, nil, nil, nil, 1, 1)
	list, err := dbRepo.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	cacheRepo.SetEmployer(ctx, list[0])
	cacheRepo.SetEmployerByEmail(ctx, list[0])

	return list[0], err
}

func (p *employerDataProvider) GetById(ctx context.Context, id int64) (*employer.Employer, error) {
	cacheRepo := cache.NewEmployerCacheRepository(p.redis)
	e, err := cacheRepo.GetEmployer(ctx, id)
	if err == nil && e != nil {
		return e, nil
	}

	pgConn, err := p.pg.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	defer pgConn.Release()

	dbRepo := repo.NewEmployerRepository(pgConn)
	query := dal.NewQueryEmployersDal([]int64{id}, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, 1, 1)
	list, err := dbRepo.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	cacheRepo.SetEmployer(ctx, list[0])
	return list[0], err
}

func (p *employerDataProvider) Save(ctx context.Context, e *employer.Employer) error {
	pgConn, err := p.pg.GetConn(ctx)
	if err != nil {
		return err
	}
	defer pgConn.Release()

	dbRepo := repo.NewEmployerRepository(pgConn)
	if e.Id() == 0 {
		if err := dbRepo.Create(ctx, e); err != nil {
			return err
		}
	} else {
		if err := dbRepo.Update(ctx, e); err != nil {
			return err
		}
	}

	cacheRepo := cache.NewEmployerCacheRepository(p.redis)
	cacheRepo.InvalidateEmployerList(ctx)
	cacheRepo.SetEmployer(ctx, e)

	return nil
}

func (p *employerDataProvider) QueryList(ctx context.Context, query *dal.QueryEmployersDal) ([]*employer.Employer, error) {
	cacheRepo := cache.NewEmployerCacheRepository(p.redis)
	list, err := cacheRepo.GetEmployerList(ctx, query)
	if err == nil && len(list) > 0 {
		return list, nil
	}

	pgConn, err := p.pg.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	defer pgConn.Release()

	dbRepo := repo.NewEmployerRepository(pgConn)
	list, err = dbRepo.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		cacheRepo.SetEmployerList(ctx, query, list)
	}

	return list, err
}
