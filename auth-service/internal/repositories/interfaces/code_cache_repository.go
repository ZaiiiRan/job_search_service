package interfaces

import (
	"context"

	"github.com/ZaiiiRan/job_search_service/auth-service/internal/domain/code"
)

type CodeCacheRepository interface {
	GetById(ctx context.Context, id int64) (*code.Code, error)
	SetById(ctx context.Context, password *code.Code) error
	DelById(ctx context.Context, id int64) error
	GetByUserId(ctx context.Context, userId int64) (*code.Code, error)
	SetByUserId(ctx context.Context, password *code.Code) error
	DelByUserId(ctx context.Context, userId int64) error
}
