package storage

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenStorage struct {
	client *pgxpool.Pool
}

func NewTokenStorage(pool *pgxpool.Pool) *TokenStorage {
	return &TokenStorage{client: pool}
}

func (s *TokenStorage) Find(ctx context.Context, refreshToken string) (string, error) {

}

func (s *TokenStorage) Save(ctx context.Context, refreshToken string) error {

}

func (s *TokenStorage) Delete(ctx context.Context, refreshToken string) error {

}
