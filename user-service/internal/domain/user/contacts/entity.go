package contacts

import (
	"github.com/ZaiiiRan/job_search_service/user-service/internal/domain/user"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/errors/validationerror"
)

type Contacts struct {
	phoneNumber *string
	telegram    *string
}

func New(phoneNumber, telegram *string) (*Contacts, validationerror.ValidationError) {
	verr := make(validationerror.ValidationError)

	c := &Contacts{}
	if err := c.SetPhoneNumber(phoneNumber); err != nil {
		verr["contacts.phone_number"] = err.Error()
	}
	if err := c.SetTelegram(telegram); err != nil {
		verr["contacts.telegram"] = err.Error()
	}

	if len(verr) > 0 {
		return nil, verr
	}
	return c, nil
}

func FromStorage(phoneNumber, telegram *string) *Contacts {
	return &Contacts{
		phoneNumber: phoneNumber,
		telegram:    telegram,
	}
}

func (c *Contacts) PhoneNumber() *string {
	return c.phoneNumber
}

func (c *Contacts) Telegram() *string {
	return c.telegram
}

func (c *Contacts) SetPhoneNumber(phoneNumber *string) error {
	if phoneNumber == nil {
		return nil
	}
	if *phoneNumber == "" {
		c.phoneNumber = nil
		return nil
	}
	if err := user.ValidatePhoneNumber(*phoneNumber); err != nil {
		return err
	}
	c.phoneNumber = phoneNumber
	return nil
}

func (c *Contacts) SetTelegram(telegram *string) error {
	if telegram == nil {
		return nil
	}
	if *telegram == "" {
		c.telegram = nil
		return nil
	}
	if err := user.ValidateTelegram(*telegram); err != nil {
		return err
	}
	c.telegram = telegram
	return nil
}
