package models

import (
	"time"

	"github.com/ZaiiiRan/job_search_service/user-service/internal/domain/user/applicant"
)

type V1ApplicantDal struct {
	Id          int64     `db:"id" json:"id"`
	FirstName   string    `db:"first_name" json:"first_name"`
	LastName    string    `db:"last_name" json:"last_name"`
	Patronymic  *string   `db:"patronymic" json:"patronymic"`
	BirthDate   string    `db:"birth_date" json:"birth_date"`
	City        string    `db:"city" json:"city"`
	Email       string    `db:"email" json:"email"`
	PhoneNumber *string   `db:"phone_number" json:"phone_number"`
	Telegram    *string   `db:"telegram" json:"telegram"`
	IsActive    bool      `db:"is_active" json:"is_active"`
	IsDeleted   bool      `db:"is_deleted" json:"is_deleted"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

func V1ApplicantDalFromDomain(a *applicant.Applicant) V1ApplicantDal {
	if a == nil {
		return V1ApplicantDal{}
	}

	return V1ApplicantDal{
		Id:          a.Id(),
		FirstName:   a.FirstName(),
		LastName:    a.LastName(),
		Patronymic:  a.Patronymic(),
		BirthDate:   a.BirthDate(),
		City:        a.City(),
		Email:       a.Email(),
		PhoneNumber: a.PhoneNumber(),
		Telegram:    a.Telegram(),
		IsActive:    a.IsActive(),
		IsDeleted:   a.IsDeleted(),
		CreatedAt:   a.CreatedAt(),
		UpdatedAt:   a.UpdatedAt(),
	}
}

func (a V1ApplicantDal) IsNull() bool { return false }
func (a V1ApplicantDal) Index(i int) any {
	switch i {
	case 0:
		return a.Id
	case 1:
		return a.FirstName
	case 2:
		return a.LastName
	case 3:
		return a.Patronymic
	case 4:
		return a.BirthDate
	case 5:
		return a.City
	case 6:
		return a.Email
	case 7:
		return a.PhoneNumber
	case 8:
		return a.Telegram
	case 9:
		return a.IsActive
	case 10:
		return a.IsDeleted
	case 11:
		return a.CreatedAt
	case 12:
		return a.UpdatedAt
	default:
		return nil
	}
}

func (a V1ApplicantDal) ToDomain() *applicant.Applicant {
	return applicant.FromStorage(
		a.Id,
		a.FirstName, a.LastName,
		a.Patronymic,
		a.BirthDate, a.City, a.Email,
		a.PhoneNumber, a.Telegram,
		a.IsActive, a.IsDeleted,
		a.CreatedAt, a.UpdatedAt,
	)
}
