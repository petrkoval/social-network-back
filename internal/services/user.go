package services

import (
	"context"
	"github.com/petrkoval/social-network-back/internal/domain"
	"github.com/rs/zerolog"
)

type UserStorage interface {
	Create(ctx context.Context, dto domain.CreateUserDTO) (*domain.User, error)
	FindByID(ctx context.Context, userID string) (*domain.User, error)
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	UpdateUsername(ctx context.Context, userID string, username string) (*domain.User, error)
	UpdatePassword(ctx context.Context, userID string, password string) (*domain.User, error)
}

type UserService struct {
	Storage UserStorage
	logger  *zerolog.Logger
}

func NewUserService(s UserStorage, l *zerolog.Logger) *UserService {
	return &UserService{
		Storage: s,
		logger:  l,
	}
}
