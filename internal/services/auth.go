package services

import (
	"context"
	"errors"
	"github.com/petrkoval/social-network-back/internal/domain"
	"github.com/petrkoval/social-network-back/internal/storage"
)

type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"-"`
	User         *domain.User `json:"user"`
}

type AuthService struct {
	tokens *TokenService
	users  *UserService
}

func NewAuthService(tokenService *TokenService, userService *UserService) *AuthService {
	return &AuthService{
		tokens: tokenService,
		users:  userService,
	}
}

func (s *AuthService) Register(ctx context.Context, dto domain.CreateUserDTO) (*AuthResponse, error) {
	var (
		err    error
		entity *domain.User
	)

	_, err = s.users.Storage.FindByUsername(ctx, dto.Username)
	if err != nil && !errors.Is(err, storage.NotFoundUserErr) {
		return nil, err
	} else if err == nil {
		return nil, UserExistsErr
	}

	entity, err = s.users.Storage.Create(ctx, dto)
	if err != nil {
		return nil, err
	}

	return s.generateAndSaveTokens(ctx, entity)
}

func (s *AuthService) Login(ctx context.Context, dto domain.CreateUserDTO) (*AuthResponse, error) {
	var (
		err    error
		entity *domain.User
	)

	entity, err = s.users.Storage.FindByUsername(ctx, dto.Username)
	if err != nil {
		return nil, err
	}

	if !(dto.Password == entity.Password) {
		return nil, WrongPasswordErr
	}

	return s.generateAndSaveTokens(ctx, entity)
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.tokens.Storage.Delete(ctx, refreshToken)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	var (
		err    error
		entity *domain.User
	)

	entity, err = s.tokens.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	_, err = s.tokens.Storage.Find(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return s.generateAndSaveTokens(ctx, entity)
}

func (s *AuthService) generateAndSaveTokens(ctx context.Context, user *domain.User) (*AuthResponse, error) {
	accessToken, refreshToken, err := s.tokens.GenerateTokens(*user)
	if err != nil {
		return nil, err
	}

	err = s.tokens.Storage.Save(ctx, domain.Token{
		UserID:       user.ID,
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}
