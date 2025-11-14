package models

import "time"

type V1CodeDal struct {
	Id        int64     `db:"id" json:"id"`
	UserId    int64     `db:"user_id" json:"user_id"`
	Code      string    `db:"code" json:"code"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (p V1CodeDal) IsNull() bool { return false }
func (p V1CodeDal) Index(i int) any {
	switch i {
	case 0:
		return p.Id
	case 1:
		return p.UserId
	case 2:
		return p.Code
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
