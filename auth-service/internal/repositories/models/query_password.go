package models

type QueryPasswordDal struct {
	Id     *int64
	UserId *int64
}

func NewQueryPasswordDal(id, userId *int64) *QueryPasswordDal {
	return &QueryPasswordDal{
		Id:     id,
		UserId: userId,
	}
}
