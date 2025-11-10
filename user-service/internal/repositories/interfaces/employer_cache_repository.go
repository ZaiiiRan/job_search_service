package interfaces

import (
	"context"

	"github.com/ZaiiiRan/job_search_service/user-service/internal/domain/user/employer"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/repositories/models"
)

type EmployerCacheRepository interface {
	SetEmployer(ctx context.Context, employer *employer.Employer) error
	GetEmployer(ctx context.Context, id int64) (*employer.Employer, error)
	DeleteEmployer(ctx context.Context, id int64) error
	SetEmployerList(ctx context.Context, query *models.QueryEmployersDal, employers []*employer.Employer) error
	GetEmployerList(ctx context.Context, query *models.QueryEmployersDal) ([]*employer.Employer, error)
	InvalidateEmployerList(ctx context.Context) error
}
