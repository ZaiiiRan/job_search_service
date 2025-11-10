package models

import (
	"time"

	"github.com/ZaiiiRan/job_search_service/user-service/internal/domain/user/employer"
)

type V1EmployerDal struct {
	Id          int64     `db:"id" json:"id"`
	CompanyName string    `db:"company_name" json:"company_name"`
	City        string    `db:"city" json:"city"`
	Email       string    `db:"email" json:"email"`
	PhoneNumber *string   `db:"phone_number" json:"phone_number"`
	Telegram    *string   `db:"telegram" json:"telegram"`
	IsActive    bool      `db:"is_active" json:"is_active"`
	IsDeleted   bool      `db:"is_deleted" json:"is_deleted"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

func V1EmployerDalFromDomain(e *employer.Employer) V1EmployerDal {
	if e == nil {
		return V1EmployerDal{}
	}

	return V1EmployerDal{
		Id:          e.Id(),
		CompanyName: e.CompanyName(),
		City:        e.City(),
		Email:       e.Email(),
		PhoneNumber: e.PhoneNumber(),
		Telegram:    e.Telegram(),
		IsActive:    e.IsActive(),
		IsDeleted:   e.IsDeleted(),
		CreatedAt:   e.CreatedAt(),
		UpdatedAt:   e.UpdatedAt(),
	}
}

func (e V1EmployerDal) IsNull() bool { return false }
func (e V1EmployerDal) Index(i int) any {
	switch i {
	case 0:
		return e.Id
	case 1:
		return e.CompanyName
	case 2:
		return e.City
	case 3:
		return e.Email
	case 4:
		return e.PhoneNumber
	case 5:
		return e.Telegram
	case 6:
		return e.IsActive
	case 7:
		return e.IsDeleted
	case 8:
		return e.CreatedAt
	case 9:
		return e.UpdatedAt
	default:
		return nil
	}
}

func (e V1EmployerDal) ToDomain() *employer.Employer {
	return employer.FromStorage(
		e.Id,
		e.CompanyName, e.City, e.Email,
		e.PhoneNumber, e.Telegram,
		e.IsActive, e.IsDeleted,
		e.CreatedAt, e.UpdatedAt,
	)
}
