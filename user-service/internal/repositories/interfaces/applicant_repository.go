package interfaces

import (
	"context"

	"github.com/ZaiiiRan/job_search_service/user-service/internal/domain/user/applicant"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/repositories/models"
)

type ApplicantRepository interface {
	Create(ctx context.Context, applicant *applicant.Applicant) error
	Update(ctx context.Context, applicant *applicant.Applicant) error
	Query(ctx context.Context, query *models.QueryApplicantsDal) ([]*applicant.Applicant, error)
}
