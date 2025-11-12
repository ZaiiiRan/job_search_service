package interfaces

import (
	"context"

	"github.com/ZaiiiRan/job_search_service/user-service/internal/domain/user/applicant"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/repositories/models"
)

type ApplicantCacheRepository interface {
	SetApplicant(ctx context.Context, applicant *applicant.Applicant) error
	GetApplicant(ctx context.Context, id int64) (*applicant.Applicant, error)
	DeleteApplicant(ctx context.Context, id int64) error
	SetApplicantByEmail(ctx context.Context, applicant *applicant.Applicant) error
	GetApplicantByEmail(ctx context.Context, email string) (*applicant.Applicant, error)
	DeleteApplicantByEmail(ctx context.Context, email string) error
	SetApplicantList(ctx context.Context, query *models.QueryApplicantsDal, applicants []*applicant.Applicant) error
	GetApplicantList(ctx context.Context, query *models.QueryApplicantsDal) ([]*applicant.Applicant, error)
	InvalidateApplicantList(ctx context.Context) error
}
