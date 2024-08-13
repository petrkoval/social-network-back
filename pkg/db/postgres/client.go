package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/petrkoval/social-network-back/internal/config"
	"github.com/rs/zerolog"
)

func NewPostgreSQLClient(c *config.DBConfig, l *zerolog.Logger) (*pgxpool.Pool, error) {
	connUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", c.User, c.Password, c.Host, c.Port, c.Database)

	cfg, err := pgxpool.ParseConfig(connUrl)
	if err != nil {
		return nil, err
	}

	cfg.ConnConfig.Tracer = &zerologPgxTracer{logger: l}

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
