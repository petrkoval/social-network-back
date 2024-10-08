package storage

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/petrkoval/social-network-back/internal/domain"
	"github.com/pkg/errors"
)

type UserStorage struct {
	client Client
}

func NewUserStorage(c *pgxpool.Pool) *UserStorage {
	return &UserStorage{client: c}
}

func (s *UserStorage) Create(ctx context.Context, dto domain.CreateUserDTO) (*domain.AuthUser, error) {
	var (
		query  = `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING user_id, username;`
		entity = &domain.AuthUser{}
		rows   pgx.Rows
		err    error
	)

	rows, err = s.client.Query(ctx, query, dto.Username, dto.Password)
	if err != nil {
		return nil, errors.Wrap(err, "UserStorage.Create")
	}
	defer rows.Close()

	err = pgxscan.ScanOne(entity, rows)
	if err != nil {
		return nil, errors.Wrap(err, "UserStorage.Create")
	}

	return entity, nil
}

func (s *UserStorage) FindByID(ctx context.Context, userID string) (*domain.User, error) {
	var (
		query = `
			SELECT user_id,
				   username,
				   password,
				   created_at,
				   coalesce(account_description, '') as account_description
			FROM users
			WHERE user_id = $1;`
		entity = &domain.User{}
		err    error
	)

	err = pgxscan.Get(ctx, s.client, entity, query, userID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, errors.Wrap(NotFoundUserErr, "UserStorage.FindByID")
		default:
			return nil, errors.Wrap(err, "UserStorage.FindByID")
		}
	}

	return entity, nil
}

func (s *UserStorage) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var (
		query = `
			SELECT user_id,
				   username,
				   password,
				   created_at,
				   coalesce(account_description, '') as account_description
			FROM users
			WHERE username = $1;`
		entity = &domain.User{}
		err    error
	)

	err = pgxscan.Get(ctx, s.client, entity, query, username)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, errors.Wrap(NotFoundUserErr, "UserStorage.FindByUsername")
		default:
			return nil, errors.Wrap(err, "UserStorage.FindByUsername")
		}
	}

	return entity, nil
}

func (s *UserStorage) UpdateUsername(ctx context.Context, userID string, username string) (*domain.User, error) {
	var (
		query  = `UPDATE users SET username = $1 WHERE user_id = $2 RETURNING *;`
		entity = &domain.User{}
		rows   pgx.Rows
		err    error
	)

	rows, err = s.client.Query(ctx, query, username, userID)
	if err != nil {
		return nil, errors.Wrap(err, "UserStorage.UpdateUsername")
	}
	defer rows.Close()

	err = pgxscan.ScanOne(&entity, rows)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, errors.Wrap(NotFoundUserErr, "UserStorage.UpdateUsername")
		default:
			return nil, errors.Wrap(err, "UserStorage.UpdateUsername")
		}
	}

	return entity, nil
}

func (s *UserStorage) UpdatePassword(ctx context.Context, userID string, password string) (*domain.User, error) {
	var (
		query  = `UPDATE users SET password = $1 WHERE user_id = $2 RETURNING *;`
		entity = &domain.User{}
		rows   pgx.Rows
		err    error
	)

	rows, err = s.client.Query(ctx, query, password, userID)
	if err != nil {
		return nil, errors.Wrap(err, "UserStorage.UpdatePassword")
	}
	defer rows.Close()

	err = pgxscan.ScanOne(&entity, rows)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, errors.Wrap(NotFoundUserErr, "UserStorage.UpdatePassword")
		default:
			return nil, errors.Wrap(err, "UserStorage.UpdatePassword")
		}
	}

	return entity, nil
}
