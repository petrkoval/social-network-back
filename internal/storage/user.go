package storage

import (
	"context"
	"errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/petrkoval/social-network-back/internal/domain"
)

var (
	insertErr   = errors.New("insert query error")
	notFoundErr = errors.New("no user found")
)

type UserStorage struct {
	client *pgxpool.Pool
}

func NewUserStorage(c *pgxpool.Pool) *UserStorage {
	return &UserStorage{client: c}
}

func (s *UserStorage) Create(ctx context.Context, dto domain.CreateUserDTO) (*domain.User, error) {
	var (
		query  = `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING *;`
		entity = &domain.User{}
		err    error
	)

	err = pgxscan.Get(ctx, s.client, entity, query, dto.Username, dto.Password)
	if err != nil {
		return nil, insertErr
	}

	return entity, nil
}

func (s *UserStorage) Find(ctx context.Context, userID string) (*domain.User, error) {
	var (
		query  = `SELECT * FROM users WHERE user_id = $1;`
		entity = &domain.User{}
		err    error
	)

	err = pgxscan.Get(ctx, s.client, entity, query, userID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, notFoundErr
		default:
			return nil, err
		}
	}

	return entity, nil
}

func (s *UserStorage) UpdateUsername(ctx context.Context, userID string, username string) (*domain.User, error) {
	var (
		query  = `UPDATE users SET username = $1 WHERE user_id = $2 RETURNING *;`
		entity = &domain.User{}
		err    error
	)

	err = pgxscan.Get(ctx, s.client, entity, query, username, userID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, notFoundErr
		default:
			return nil, err
		}
	}

	return entity, nil
}

func (s *UserStorage) UpdatePassword(ctx context.Context, userID string, password string) (*domain.User, error) {
	var (
		query  = `UPDATE users SET password = $1 WHERE user_id = $2 RETURNING *;`
		entity = &domain.User{}
		err    error
	)

	err = pgxscan.Get(ctx, s.client, entity, query, password, userID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, notFoundErr
		default:
			return nil, err
		}
	}

	return entity, nil
}
