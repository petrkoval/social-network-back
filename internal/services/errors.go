package services

import "errors"

var (
	JwtSigningErr              = errors.New("error while signing jwt")
	UnexpectedSigningMethodErr = errors.New("unexpected signing method")
	TokenExpiredErr            = errors.New("token is expired")
	InvalidTokenErr            = errors.New("invalid token")

	UserExistsErr = errors.New("user already exists")
)
