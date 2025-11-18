package code

import "time"

const (
	maxGenerationsLeft = 3
	codeTTL            = 10 * time.Minute
)

type Code struct {
	id              int64
	userId          int64
	code            string
	generationsLeft int
	expiresAt       time.Time
	createdAt       time.Time
	updatedAt       time.Time
}

func New(userId int64) (*Code, error) {
	c := &Code{}
	c.generationsLeft = maxGenerationsLeft
	c.userId = userId

	if err := c.GenerateCode(); err != nil {
		return nil, err
	}

	now := time.Now()
	c.createdAt = now
	c.updatedAt = now

	return c, nil
}

func FromStorage(
	id, userId int64,
	code string,
	generationsLeft int,
	expiresAt, createdAt, updatedAt time.Time,
) *Code {
	return &Code{
		id:        id,
		userId:    userId,
		code:      code,
		expiresAt: expiresAt,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (c *Code) Id() int64            { return c.id }
func (c *Code) UserId() int64        { return c.userId }
func (c *Code) Code() string         { return c.code }
func (c *Code) GenerationsLeft() int { return c.generationsLeft }
func (c *Code) ExpiresAt() time.Time { return c.expiresAt }
func (c *Code) CreatedAt() time.Time { return c.createdAt }
func (c *Code) UpdatedAt() time.Time { return c.updatedAt }

func (c *Code) SetId(id int64) {
	if c.Id() == 0 {
		c.id = id
	}
}

func (c *Code) GenerateCode() error {
	if c.Id() != 0 {
		if c.generationsLeft <= 0 && time.Since(c.updatedAt) < 5*time.Minute {
			return NewCodeValidationError("the number of code resends has been exhausted")
		} else if c.generationsLeft <= 0 {
			c.generationsLeft = maxGenerationsLeft
		}
	}
	c.generationsLeft--

	code, err := generateSixDigitCode()
	if err != nil {
		return err
	}
	c.code = code
	c.expiresAt = time.Now().Add(codeTTL)
	c.updatedAt = time.Now()
	return nil
}

func (c *Code) CheckCode(rawCode string) (bool, error) {
	if time.Now().After(c.expiresAt) {
		return false, NewCodeValidationError("code has been expired")
	}
	if c.code == rawCode {
		return true, nil
	}
	return false, nil
}