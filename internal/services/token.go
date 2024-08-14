package services

import "context"

type TokenStorage interface {
	Find(ctx context.Context, refreshToken string) (string, error)
	Save(ctx context.Context, refreshToken string) error
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
