package password

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	id           int64
	userId       int64
	passwordHash string
	createdAt    time.Time
	updatedAt    time.Time
}

func New(userId int64, password string) (*Password, error) {
	p := &Password{}

	if err := p.SetPassword(password); err != nil {
		return nil, err
	}

	now := time.Now()
	p.createdAt = now
	p.updatedAt = now

	return p, nil
}

func FromStorage(
	id, userId int64,
	passwordHash string,
	createdAt, updatedAt time.Time,
) *Password {
	return &Password{
		id:           id,
		userId:       userId,
		passwordHash: passwordHash,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

func (p *Password) Id() int64            { return p.id }
func (p *Password) UserId() int64        { return p.userId }
func (p *Password) PasswordHash() string { return p.passwordHash }
func (p *Password) CreatedAt() time.Time { return p.createdAt }
func (p *Password) UpdatedAt() time.Time { return p.updatedAt }

func (p *Password) SetId(id int64) {
	if p.Id() == 0 {
		p.id = id
	}
}

func (p *Password) SetPassword(password string) error {
	if err := ValidatePassword(password); err != nil {
		return err
	}

	if p.Id() != 0 {
		if time.Since(p.updatedAt) < 24*time.Hour {
			return NewPasswordValidationError("password can be changed only once per 24 hours")
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.passwordHash = string(hash)
	p.updatedAt = time.Now()
	return nil
}

func (p *Password) Check(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(p.passwordHash), []byte(password)); err != nil {
		return false
	} else {
		return true
	}
}
