package models

type QueryTokenDal struct {
	Id      *int64
	UserId  *int64
	Token   *string
	Version *int
}

func NewQueryTokenDal(id, userId *int64, token *string, version *int) *QueryTokenDal {
	return &QueryTokenDal{
		Id:      id,
		UserId:  userId,
		Token:   token,
		Version: version,
	}
}
