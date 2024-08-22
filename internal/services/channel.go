package services

import (
	"context"
	"github.com/petrkoval/social-network-back/internal/config"
	"github.com/petrkoval/social-network-back/internal/domain"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"strconv"
)

type ChannelStorage interface {
	FindAll(ctx context.Context, limit, offset int) ([]*domain.Channel, error)
	FindByUserID(ctx context.Context, userID string) ([]*domain.Channel, error)
	FindByID(ctx context.Context, id string) (*domain.Channel, error)
	Create(ctx context.Context, dto domain.CreateChannelDTO) (*domain.Channel, error)
	Update(ctx context.Context, id string, dto domain.UpdateChannelDTO) (*domain.Channel, error)
	Delete(ctx context.Context, id string) error
}

type ChannelService struct {
	ChannelStorage
	logger *zerolog.Logger
	cfg    *config.TokensConfig
}

func NewChannelService(s ChannelStorage, l *zerolog.Logger, c *config.TokensConfig) *ChannelService {
	return &ChannelService{
		ChannelStorage: s,
		logger:         l,
		cfg:            c,
	}
}

func (s *ChannelService) FindAll(ctx context.Context, limit, offset string) ([]*domain.Channel, error) {
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return nil, errors.Wrap(QueryParamParsingErr, "ChannelService.FindAll")
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return nil, errors.Wrap(QueryParamParsingErr, "ChannelService.FindAll")
	}

	return s.ChannelStorage.FindAll(ctx, limitInt, offsetInt)
}
