package jwt

import "github.com/golang-jwt/jwt/v5"

type ApplicantClaims struct {
	Id         int64
	FirstName  string
	LastName   string
	Patronymic *string
	Email      string
	IsActive   bool
	IsDeleted  bool
	jwt.RegisteredClaims
}
