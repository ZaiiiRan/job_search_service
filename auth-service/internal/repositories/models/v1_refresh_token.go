package models

import (
	"time"

	"github.com/ZaiiiRan/job_search_service/auth-service/internal/domain/token"
)

type V1RefreshTokenDal struct {
	Id        int64     `db:"id" json:"id"`
	UserId    int64     `db:"user_id" json:"user_id"`
	Token     string    `db:"token" json:"token"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func V1RefreshTokenDalFromDomain(t *token.Token) V1RefreshTokenDal {
	if t == nil {
		return V1RefreshTokenDal{}
	}

	return V1RefreshTokenDal{
		Id:        t.Id(),
		UserId:    t.UserId(),
		Token:     t.Token(),
		ExpiresAt: t.ExpiresAt(),
		CreatedAt: t.CreatedAt(),
		UpdatedAt: t.UpdatedAt(),
	}
}

func (p V1RefreshTokenDal) IsNull() bool { return false }
func (p V1RefreshTokenDal) Index(i int) any {
	switch i {
	case 0:
		return p.Id
	case 1:
		return p.UserId
	case 2:
		return p.Token
	case 3:
		return p.ExpiresAt
	case 4:
		return p.CreatedAt
	case 5:
		return p.UpdatedAt
	default:
		return nil
	}
}

func (p V1RefreshTokenDal) ToDomain() *token.Token {
	return token.FromStorage(
		p.Id, p.UserId,
		p.Token, "refresh",
		p.ExpiresAt, p.CreatedAt, p.UpdatedAt,
	)
}
