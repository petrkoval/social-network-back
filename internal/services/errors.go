package services

import "errors"

var (
	jwtSigningErr              = errors.New("error while signing jwt")
	unexpectedSigningMethodErr = errors.New("unexpected signing method")
	tokenExpiredErr            = errors.New("token is expired")
	invalidTokenErr            = errors.New("invalid token")
)
