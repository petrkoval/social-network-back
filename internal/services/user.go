package services

import (
	"context"
	"github.com/petrkoval/social-network-back/internal/domain"
)

type UserStorage interface {
	Create(ctx context.Context, dto domain.CreateUserDTO) (*domain.User, error)
	Find(ctx context.Context, userID string) (*domain.User, error)
	UpdateUsername(ctx context.Context, userID string, username string) (*domain.User, error)
	UpdatePassword(ctx context.Context, userID string, password string) (*domain.User, error)
}

type UserService struct {
	storage UserStorage
}

func NewUserService(s UserStorage) *UserService {
	return &UserService{storage: s}
}
