package auth

import "github.com/petrkoval/social-network-back/internal/domain"

type Response struct {
	AccessToken string      `json:"access_token"`
	User        domain.User `json:"user"`
}
