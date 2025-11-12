package employer

import (
	"time"

	"github.com/ZaiiiRan/job_search_service/user-service/internal/domain/user"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/domain/user/contacts"
	"github.com/ZaiiiRan/job_search_service/common/pkg/errors/validationerror"
)

type Employer struct {
	id          int64
	companyName string
	city        string
	email       string
	contacts    *contacts.Contacts
	isActive    bool
	isDeleted   bool
	createdAt   time.Time
	updatedAt   time.Time
}

func New(
	companyName, city, email string,
	phoneNumber, telegram *string,
	isActive, isDeleted bool,
) (*Employer, validationerror.ValidationError) {
	verr := make(validationerror.ValidationError)

	e := &Employer{}
	if err := e.SetCompanyName(companyName); err != nil {
		verr["company_name"] = err.Error()
	}
	if err := e.SetCity(city); err != nil {
		verr["city"] = err.Error()
	}
	if err := e.SetEmail(email); err != nil {
		verr["email"] = err.Error()
	}
	contacts, err := contacts.New(phoneNumber, telegram)
	if err != nil {
		verr.Merge(err)
	}

	if len(verr) > 0 {
		return nil, verr
	}

	e.contacts = contacts

	now := time.Now()
	e.createdAt = now
	e.updatedAt = now
	return e, nil
}

func FromStorage(
	id int64,
	companyName, city, email string,
	phoneNumber, telegram *string,
	isActive, isDeleted bool,
	createdAt, updatedAt time.Time,
) *Employer {
	return &Employer{
		id:          id,
		companyName: companyName,
		city:        city,
		email:       email,
		contacts:    contacts.FromStorage(phoneNumber, telegram),
		isActive:    isActive,
		isDeleted:   isDeleted,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (e *Employer) Id() int64            { return e.id }
func (e *Employer) CompanyName() string  { return e.companyName }
func (e *Employer) City() string         { return e.city }
func (e *Employer) Email() string        { return e.email }
func (e *Employer) PhoneNumber() *string { return e.contacts.PhoneNumber() }
func (e *Employer) Telegram() *string    { return e.contacts.Telegram() }
func (e *Employer) IsActive() bool       { return e.isActive }
func (e *Employer) IsDeleted() bool      { return e.isDeleted }
func (e *Employer) CreatedAt() time.Time { return e.createdAt }
func (e *Employer) UpdatedAt() time.Time { return e.updatedAt }

func (e *Employer) SetId(id int64) {
	if e.Id() == 0 {
		e.id = id
	}
}

func (e *Employer) SetCompanyName(companyName string) error {
	if err := user.ValidateCompanyName(companyName); err != nil {
		return err
	}
	e.companyName = companyName
	return nil
}

func (e *Employer) SetCity(city string) error {
	if err := user.ValidateName("city", city, true); err != nil {
		return err
	}
	e.city = city
	return nil
}

func (e *Employer) SetEmail(email string) error {
	if err := user.ValidateEmail(email); err != nil {
		return err
	}
	e.email = email
	return nil
}

func (e *Employer) SetPhoneNumber(phoneNumber *string) error {
	if err := e.contacts.SetPhoneNumber(phoneNumber); err != nil {
		return err
	}
	return nil
}

func (e *Employer) SetTelegram(telegram *string) error {
	if err := e.contacts.SetTelegram(telegram); err != nil {
		return err
	}
	return nil
}

func (e *Employer) SetIsActive(isActive bool) {
	e.isActive = isActive
}

func (e *Employer) SetIsDeleted(isDeleted bool) {
	e.isDeleted = isDeleted
}

func (e *Employer) SetUpdatedAt(updatedAt time.Time) {
	e.updatedAt = updatedAt
}
