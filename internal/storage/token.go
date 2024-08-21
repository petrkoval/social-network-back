package storage

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/petrkoval/social-network-back/internal/domain"
	"github.com/pkg/errors"
)

type TokenStorage struct {
	client Client
}

func NewTokenStorage(pool *pgxpool.Pool) *TokenStorage {
	return &TokenStorage{client: pool}
}

func (s *TokenStorage) FindByToken(ctx context.Context, refreshToken string) (*domain.Token, error) {
	var (
		query  = `SELECT * FROM tokens WHERE refresh_token = $1;`
		entity = &domain.Token{}
		err    error
	)

	err = pgxscan.Get(ctx, s.client, entity, query, refreshToken)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, errors.Wrap(NotFoundTokenErr, "TokenStorage.FindByToken")
		default:
			return nil, errors.Wrap(err, "TokenStorage.FindByToken")
		}
	}

	return entity, nil
}

func (s *TokenStorage) FindByUserID(ctx context.Context, userID string) (*domain.Token, error) {
	var (
		query  = `SELECT * FROM tokens WHERE user_id = $1;`
		entity = &domain.Token{}
		err    error
	)

	err = pgxscan.Get(ctx, s.client, entity, query, userID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, errors.Wrap(NotFoundTokenErr, "TokenStorage.FindByToken")
		default:
			return nil, errors.Wrap(err, "TokenStorage.FindByToken")
		}
	}

	return entity, nil
}

func (s *TokenStorage) Save(ctx context.Context, token domain.Token) error {
	var err error

	_, err = s.FindByUserID(ctx, token.UserID)
	if err != nil {
		switch {
		case errors.Is(err, NotFoundTokenErr):
			err = s.create(ctx, token)
			if err != nil {
				return errors.Wrap(err, "TokenStorage.Save")
			}
		default:
			return errors.Wrap(err, "TokenStorage.Save")
		}
	}

	err = s.update(ctx, token)
	if err != nil {
		return errors.Wrap(err, "TokenStorage.Save")
	}

	return nil
}

func (s *TokenStorage) Delete(ctx context.Context, refreshToken string) error {
	var (
		query = `DELETE FROM tokens WHERE refresh_token = $1;`
		err   error
	)

	_, err = s.client.Exec(ctx, query, refreshToken)
	if err != nil {
		return errors.Wrap(err, "TokenStorage.Delete")
	}

	return nil
}

func (s *TokenStorage) update(ctx context.Context, token domain.Token) error {
	var (
		query = `UPDATE tokens SET refresh_token = $1 WHERE user_id = $2;`
		err   error
	)

	_, err = s.client.Exec(ctx, query, token.RefreshToken, token.UserID)
	if err != nil {
		return errors.Wrap(err, "TokenStorage.update")
	}

	return nil
}

func (s *TokenStorage) create(ctx context.Context, token domain.Token) error {
	var (
		query = `INSERT INTO tokens (user_id, refresh_token) VALUES ($1, $2);`
		err   error
	)

	_, err = s.client.Exec(ctx, query, token.UserID, token.RefreshToken)
	if err != nil {
		return errors.Wrap(err, "TokenStorage.create")
	}

	return nil
}
