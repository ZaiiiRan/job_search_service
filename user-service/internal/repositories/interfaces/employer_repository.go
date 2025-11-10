package interfaces

import (
	"context"

	"github.com/ZaiiiRan/job_search_service/user-service/internal/domain/user/employer"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/repositories/models"
)

type EmployerRepository interface {
	Create(ctx context.Context, employer *employer.Employer) error
	Update(ctx context.Context, employer *employer.Employer) error
	Query(ctx context.Context, query *models.QueryEmployersDal) ([]*employer.Employer, error)
}
