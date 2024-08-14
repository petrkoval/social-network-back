package storage

import "github.com/jackc/pgx/v5/pgxpool"

type UserStorage struct {
	client *pgxpool.Pool
}

func NewUserStorage(c *pgxpool.Pool) *UserStorage {
	return &UserStorage{client: c}
}
