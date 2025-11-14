package models

type V1RefreshTokenDal struct {
	Id        int64  `db:"id" json:"id"`
	UserId    int64  `db:"user_id" json:"user_id"`
	Token     string `db:"token" json:"token"`
	ExpiresAt int64  `db:"expires_at" json:"expires_at"`
	CreatedAt int64  `db:"created_at" json:"created_at"`
	UpdatedAt int64  `db:"updated_at" json:"updated_at"`
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
