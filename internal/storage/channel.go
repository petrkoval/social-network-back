package storage

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/petrkoval/social-network-back/internal/domain"
	"github.com/petrkoval/social-network-back/internal/services"
	"github.com/pkg/errors"
)

type channelStorage struct {
	client Client
}

func NewChannelStorage(client Client) services.ChannelStorage {
	return &channelStorage{client: client}
}

func (c channelStorage) FindAll(ctx context.Context, limit, offset int) ([]*domain.Channel, error) {
	var (
		channels []*domain.Channel
		err      error
		query    = `SELECT * FROM channels LIMIT $1 OFFSET $2`
	)

	err = pgxscan.Select(ctx, c.client, &channels, query, limit, offset)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.Wrap(err, "channelStorage.FindAll")
	}

	return channels, nil
}

func (c channelStorage) FindByUserID(ctx context.Context, userID string) ([]*domain.Channel, error) {
	var (
		channels []*domain.Channel
		err      error
		query    = `SELECT * FROM channels WHERE user_id = $1`
	)

	err = pgxscan.Select(ctx, c.client, &channels, query, userID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.Wrap(err, "channelStorage.FindByUserID")
	}

	return channels, nil
}

func (c channelStorage) FindByID(ctx context.Context, id string) (*domain.Channel, error) {
	var (
		channel *domain.Channel
		err     error
		query   = `SELECT * FROM channels WHERE channel_id = $1`
	)

	err = pgxscan.Get(ctx, c.client, &channel, query, id)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, errors.Wrap(NotFoundChannelErr, "channelStorage.FindByID")
		default:
			return nil, errors.Wrap(err, "channelStorage.FindByID")
		}
	}

	return channel, nil
}

func (c channelStorage) Create(ctx context.Context, dto domain.CreateChannelDTO) (*domain.Channel, error) {
	var (
		channel *domain.Channel
		query   = `INSERT INTO channels (user_id, title, description) VALUES ($1, $2, $3) RETURNING *`
	)

	rows, err := c.client.Query(ctx, query, dto.UserID, dto.Title, dto.Description)
	if err != nil {
		return nil, errors.Wrap(err, "channelStorage.Create")
	}
	defer rows.Close()

	err = pgxscan.ScanOne(&channel, rows)
	if err != nil {
		return nil, errors.Wrap(err, "channelStorage.Create")
	}

	return channel, nil
}

func (c channelStorage) Update(ctx context.Context, id string, dto domain.UpdateChannelDTO) (*domain.Channel, error) {
	var (
		channel *domain.Channel
		query   = `UPDATE channels SET title = $1, description = $2 WHERE channel_id = $3 RETURNING *`
	)

	rows, err := c.client.Query(ctx, query, id, dto.Title, dto.Description)
	if err != nil {
		return nil, errors.Wrap(err, "channelStorage.Create")
	}
	defer rows.Close()

	err = pgxscan.ScanOne(&channel, rows)
	if err != nil {
		return nil, errors.Wrap(err, "channelStorage.Create")
	}

	return channel, nil
}

func (c channelStorage) Delete(ctx context.Context, id string) error {
	var (
		err   error
		query = `DELETE FROM channels WHERE channel_id = $1`
	)

	_, err = c.client.Exec(ctx, query, id)
	if err != nil {
		return errors.Wrap(err, "channelStorage.Delete")
	}

	return nil
}
