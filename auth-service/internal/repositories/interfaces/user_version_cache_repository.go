package interfaces

import (
	"context"

	userversion "github.com/ZaiiiRan/job_search_service/auth-service/internal/domain/user_version"
)

type UserVersionCacheRepository interface {
	GetByUserId(ctx context.Context, userId int64) (*userversion.UserVersion, error)
	SetByUserId(ctx context.Context, uv *userversion.UserVersion) error
	DelByUserId(ctx context.Context, userId int64) error
}
