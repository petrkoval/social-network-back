package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func NewPostgreSQLClient(c *config) (*pgxpool.Pool, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.User, c.Password, c.Host, c.Port, c.Database)

	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
