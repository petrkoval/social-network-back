package storage

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/petrkoval/social-network-back/internal/domain"
	"github.com/pkg/errors"
)

type ChannelStorage struct {
	client Client
}

func NewChannelStorage(client Client) *ChannelStorage {
	return &ChannelStorage{client: client}
}

func (c ChannelStorage) FindAll(ctx context.Context, limit, offset int) ([]*domain.Channel, error) {
	var (
		channels []*domain.Channel
		err      error
		query    = `SELECT * FROM channels LIMIT $1 OFFSET $2`
	)

	err = pgxscan.Select(ctx, c.client, &channels, query, limit, offset)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.Wrap(err, "ChannelStorage.FindAll")
	}

	return channels, nil
}

func (c ChannelStorage) FindByUserID(ctx context.Context, userID string) ([]*domain.Channel, error) {
	var (
		channels []*domain.Channel
		err      error
		query    = `SELECT * FROM channels WHERE user_id = $1`
	)

	err = pgxscan.Select(ctx, c.client, &channels, query, userID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.Wrap(err, "ChannelStorage.FindByUserID")
	}

	return channels, nil
}

func (c ChannelStorage) FindByID(ctx context.Context, id string) (*domain.Channel, error) {
	var (
		channel *domain.Channel
		err     error
		query   = `SELECT * FROM channels WHERE channel_id = $1`
	)

	err = pgxscan.Get(ctx, c.client, &channel, query, id)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, errors.Wrap(NotFoundChannelErr, "ChannelStorage.FindByID")
		default:
			return nil, errors.Wrap(err, "ChannelStorage.FindByID")
		}
	}

	return channel, nil
}

func (c ChannelStorage) Create(ctx context.Context, dto domain.CreateChannelDTO) (*domain.Channel, error) {
	var (
		channel *domain.Channel
		query   = `INSERT INTO channels (user_id, title, description) VALUES ($1, $2, $3) RETURNING *`
	)

	rows, err := c.client.Query(ctx, query, dto.UserID, dto.Title, dto.Description)
	if err != nil {
		return nil, errors.Wrap(err, "ChannelStorage.Create")
	}
	defer rows.Close()

	err = pgxscan.ScanOne(&channel, rows)
	if err != nil {
		return nil, errors.Wrap(err, "ChannelStorage.Create")
	}

	return channel, nil
}

func (c ChannelStorage) Update(ctx context.Context, id string, dto domain.UpdateChannelDTO) (*domain.Channel, error) {
	var (
		channel *domain.Channel
		query   = `UPDATE channels SET title = $1, description = $2 WHERE channel_id = $3 RETURNING *`
	)

	rows, err := c.client.Query(ctx, query, id, dto.Title, dto.Description)
	if err != nil {
		return nil, errors.Wrap(err, "ChannelStorage.Create")
	}
	defer rows.Close()

	err = pgxscan.ScanOne(&channel, rows)
	if err != nil {
		return nil, errors.Wrap(err, "ChannelStorage.Create")
	}

	return channel, nil
}

func (c ChannelStorage) Delete(ctx context.Context, id string) error {
	var (
		err   error
		query = `DELETE FROM channels WHERE channel_id = $1`
	)

	_, err = c.client.Exec(ctx, query, id)
	if err != nil {
		return errors.Wrap(err, "ChannelStorage.Delete")
	}

	return nil
}
