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

	return make([]*domain.Channel, 0), nil
}

func (c channelStorage) FindByUserID(ctx context.Context, serID string) ([]*domain.Channel, error) {
	//TODO implement me
	panic("implement me")
}

func (c channelStorage) FindByID(ctx context.Context, id string) (*domain.Channel, error) {
	//TODO implement me
	panic("implement me")
}

func (c channelStorage) Create(ctx context.Context, dto domain.CreateChannelDTO) (*domain.Channel, error) {
	//TODO implement me
	panic("implement me")
}

func (c channelStorage) Update(ctx context.Context, id string, dto domain.UpdateChannelDTO) (*domain.Channel, error) {
	//TODO implement me
	panic("implement me")
}

func (c channelStorage) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
