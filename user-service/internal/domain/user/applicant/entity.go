package applicant

import (
	"time"

	"github.com/ZaiiiRan/job_search_service/user-service/internal/domain/user"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/domain/user/contacts"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/errors/validationerror"
)

type Applicant struct {
	id         int64
	firstName  string
	lastName   string
	patronymic *string
	birthDate  string
	city       string
	email      string
	contacts   *contacts.Contacts
	isActive   bool
	isDeleted  bool
	createdAt  time.Time
	updatedAt  time.Time
}

func New(
	firstName, lastName string,
	patronymic *string,
	birthDate, city, email string,
	phoneNumber, telegram *string,
	isActive, isDeleted bool,
) (*Applicant, validationerror.ValidationError) {
	verr := make(validationerror.ValidationError)

	a := &Applicant{}
	if err := a.SetFirstName(firstName); err != nil {
		verr["first_name"] = err.Error()
	}
	if err := a.SetLastName(lastName); err != nil {
		verr["last_name"] = err.Error()
	}
	if err := a.SetPatronymic(patronymic); err != nil {
		verr["patronymic"] = err.Error()
	}
	if err := a.SetBirthDate(birthDate); err != nil {
		verr["birth_date"] = err.Error()
	}
	if err := a.SetCity(city); err != nil {
		verr["city"] = err.Error()
	}
	if err := a.SetEmail(email); err != nil {
		verr["email"] = err.Error()
	}
	contacts, err := contacts.New(phoneNumber, telegram)
	if err != nil {
		verr.Merge(err)
	}

	if len(verr) > 0 {
		return nil, verr
	}

	a.contacts = contacts

	now := time.Now()
	a.createdAt = now
	a.updatedAt = now
	return a, nil
}

func FromStorage(
	id int64,
	firstName, lastName string,
	patronymic *string,
	birthDate, city, email string,
	phoneNumber, telegram *string,
	isActive, isDeleted bool,
	createdAt, updatedAt time.Time,
) *Applicant {
	return &Applicant{
		id:         id,
		firstName:  firstName,
		lastName:   lastName,
		patronymic: patronymic,
		birthDate:  birthDate,
		city:       city,
		email:      email,
		contacts:   contacts.FromStorage(phoneNumber, telegram),
		isActive:   isActive,
		isDeleted:  isDeleted,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
	}
}

func (a *Applicant) Id() int64            { return a.id }
func (a *Applicant) FirstName() string    { return a.firstName }
func (a *Applicant) LastName() string     { return a.lastName }
func (a *Applicant) Patronymic() *string  { return a.patronymic }
func (a *Applicant) BirthDate() string    { return a.birthDate }
func (a *Applicant) City() string         { return a.city }
func (a *Applicant) Email() string        { return a.email }
func (a *Applicant) PhoneNumber() *string { return a.contacts.PhoneNumber() }
func (a *Applicant) Telegram() *string    { return a.contacts.Telegram() }
func (a *Applicant) IsActive() bool       { return a.isActive }
func (a *Applicant) IsDeleted() bool      { return a.isDeleted }
func (a *Applicant) CreatedAt() time.Time { return a.createdAt }
func (a *Applicant) UpdatedAt() time.Time { return a.updatedAt }

func (a *Applicant) SetId(id int64) {
	if a.Id() == 0 {
		a.id = id
	}
}

func (a *Applicant) SetFirstName(firstName string) error {
	if err := user.ValidateName("first_name", firstName, true); err != nil {
		return err
	}
	a.firstName = firstName
	return nil
}

func (a *Applicant) SetLastName(lastName string) error {
	if err := user.ValidateName("last_name", lastName, true); err != nil {
		return err
	}
	a.lastName = lastName
	return nil
}

func (a *Applicant) SetPatronymic(patronymic *string) error {
	if patronymic == nil {
		return nil
	}
	if *patronymic == "" {
		a.patronymic = nil
		return nil
	}
	if err := user.ValidateName("patronymic", *patronymic, false); err != nil {
		return err
	}
	a.patronymic = patronymic
	return nil
}

func (a *Applicant) SetBirthDate(birthDate string) error {
	if err := user.ValidateBirthDate(birthDate); err != nil {
		return err
	}
	a.birthDate = birthDate
	return nil
}

func (a *Applicant) SetCity(city string) error {
	if err := user.ValidateName("city", city, true); err != nil {
		return err
	}
	a.city = city
	return nil
}

func (a *Applicant) SetEmail(email string) error {
	if err := user.ValidateEmail(email); err != nil {
		return err
	}
	a.email = email
	return nil
}

func (a *Applicant) SetPhoneNumber(phoneNumber *string) error {
	if err := a.contacts.SetPhoneNumber(phoneNumber); err != nil {
		return err
	}
	return nil
}

func (a *Applicant) SetTelegram(telegram *string) error {
	if err := a.contacts.SetTelegram(telegram); err != nil {
		return err
	}
	return nil
}

func (a *Applicant) SetIsActive(isActive bool) {
	a.isActive = isActive
}

func (a *Applicant) SetIsDeleted(isDeleted bool) {
	a.isDeleted = isDeleted
}

func (a *Applicant) SetUpdatedAt(updatedAt time.Time) {
	a.updatedAt = updatedAt
}
