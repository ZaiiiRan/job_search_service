package models

type QueryUserVersionDal struct {
	Id     *int64
	UserId *int64
}

func NewQueryUserVersionDal(id, userId *int64) *QueryUserVersionDal {
	return &QueryUserVersionDal{
		Id:     id,
		UserId: userId,
	}
}
