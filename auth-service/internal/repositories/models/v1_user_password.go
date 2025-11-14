package models

import "time"

type V1UserPasswordDal struct {
	Id          int64     `db:"id" json:"id"`
	UserId      int64     `db:"user_id" json:"user_id"`
	PassordHash string    `db:"password_hash" json:"password_hash"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

func (p V1UserPasswordDal) IsNull() bool { return false }
func (p V1UserPasswordDal) Index(i int) any {
	switch i {
	case 0:
		return p.Id
	case 1:
		return p.UserId
	case 2:
		return p.PassordHash
	case 3:
		return p.CreatedAt
	case 4:
		return p.UpdatedAt
	default:
		return nil
	}
}
