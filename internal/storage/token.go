package storage

import "github.com/jackc/pgx/v5/pgxpool"

type TokenStorage struct {
	client *pgxpool.Pool
}

func NewTokenStorage(pool *pgxpool.Pool) *TokenStorage {
	return &TokenStorage{client: pool}
}

func (s *TokenStorage) FindRefreshToken(refreshToken string) (string, error) {

}

func (s *TokenStorage) SaveRefreshToken(refreshToken string) error {

}

func (s *TokenStorage) DeleteRefreshToken(refreshToken string) error {

}
