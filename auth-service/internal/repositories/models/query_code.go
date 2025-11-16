package models

type QueryCodeDal struct {
	Id     *int64
	UserId *int64
}

func NewQueryCodeDal(id, userId *int64) *QueryCodeDal {
	return &QueryCodeDal{
		Id:     id,
		UserId: userId,
	}
}
