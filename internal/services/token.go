package services

import (
	"context"
	"github.com/petrkoval/social-network-back/internal/domain"
)

type TokenStorage interface {
	Find(ctx context.Context, refreshToken string) (*domain.Token, error)
	Save(ctx context.Context, token domain.Token) error
	Delete(ctx context.Context, refreshToken string) error
}

type TokenService struct {
	storage TokenStorage
}

func NewTokenService(s TokenStorage) *TokenService {
	return &TokenService{storage: s}
}

func (s *TokenService) GenerateTokens(userID string) (string, error) {

}

func (s *TokenService) VerifyAccessToken(accessToken string) (string, error) {

}

func (s *TokenService) VerifyRefreshToken(refreshToken string) (string, error) {

}
