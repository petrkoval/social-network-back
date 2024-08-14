package domain

import "github.com/golang-jwt/jwt/v5"

type Token struct {
	UserID       string `json:"user_id"       db:"user_id"`
	RefreshToken string `json:"refresh_token" db:"refresh_token"`
}

type TokenClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
