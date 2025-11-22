package jwt

import "github.com/golang-jwt/jwt/v5"

type EmployerClaims struct {
	Id          int64
	CompanyName string
	Email       string
	IsActive    bool
	IsDeleted   bool
	Version     int
	jwt.RegisteredClaims
}
