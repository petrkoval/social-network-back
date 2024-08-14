package services

import "github.com/petrkoval/social-network-back/internal/domain"

type AuthResponse struct {
	AccessToken string      `json:"access_token"`
	User        domain.User `json:"user"`
}

type AuthService struct {
	tokenService TokenService
	userService  UserService
}

func NewAuthService(tokenService TokenService, userService UserService) *AuthService {
	return &AuthService{
		tokenService: tokenService,
		userService:  userService,
	}
}

func (s *AuthService) Register(dto domain.CreateUserDTO) (*AuthResponse, error) {

}

func (s *AuthService) Login(dto domain.CreateUserDTO) (*AuthResponse, error) {

}

func (s *AuthService) Logout(refreshToken string) error {

}

func (s *AuthService) Refresh(refreshToken string) (*AuthResponse, error) {

}
