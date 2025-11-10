package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultMinConns            = 0
	defaultMaxConns            = 4
	defaultMaxConnLifetime     = time.Hour
	defaultMaxConnIdleTime     = 30 * time.Minute
	defaultHealthCheckInterval = time.Minute
	defaultConnectTimeout      = 5 * time.Second
)

func Connect(ctx context.Context, databaseUrl string) (*pgxpool.Pool, error) {
	cfg, err := config(databaseUrl)
	if err != nil {
		return nil, err
	}

	connPool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	err = connPool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return connPool, nil
}

func config(url string) (*pgxpool.Config, error) {
	dbConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckInterval
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	return dbConfig, nil
}
