package services

type Storage interface {
	FindRefreshToken(refreshToken string) (string, error)
	SaveRefreshToken(refreshToken string) error
	DeleteRefreshToken(refreshToken string) error
}

type TokenService struct {
	storage Storage
}

func NewTokenService(s Storage) *TokenService {
	return &TokenService{storage: s}
}

func (s *TokenService) GenerateTokens(userID string) (string, error) {

}

func (s *TokenService) VerifyAccessToken(accessToken string) (string, error) {

}

func (s *TokenService) VerifyRefreshToken(refreshToken string) (string, error) {

}
