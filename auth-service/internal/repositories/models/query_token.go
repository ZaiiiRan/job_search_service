package models

type QueryTokenDal struct {
	Id     *int64
	UserId *int64
	Token  *string
}

func NewQueryTokenDal(id, userId *int64, token *string) *QueryTokenDal {
	return &QueryTokenDal{
		Id:     id,
		UserId: userId,
		Token:  token,
	}
}
