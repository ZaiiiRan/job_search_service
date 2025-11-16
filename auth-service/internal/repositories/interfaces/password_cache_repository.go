package interfaces

import (
	"context"

	"github.com/ZaiiiRan/job_search_service/auth-service/internal/domain/password"
)

type PasswordCacheRepository interface {
	GetById(ctx context.Context, id int64) (*password.Password, error)
	SetById(ctx context.Context, password *password.Password) error
	DelById(ctx context.Context, id int64) error
	GetByUserId(ctx context.Context, userId int64) (*password.Password, error)
	SetByUserId(ctx context.Context, password *password.Password) error
	DelByUserId(ctx context.Context, userId int64) error
}
