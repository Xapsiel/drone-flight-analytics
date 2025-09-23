package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Xapsiel/bpla_dashboard/internal/config"
	"github.com/Xapsiel/bpla_dashboard/internal/migration"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func NewPostgresDB(config config.DatabaseConfig) (*pgxpool.Pool, error) {
	PgUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", config.User, config.Password, config.Host, config.Port, config.Name, config.Sslmode)

	if err := migration.Migrate(PgUrl); err != nil {
		return nil, fmt.Errorf("migration error: %s", err.Error())
	}

	poolConfig, err := pgxpool.ParseConfig(PgUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pool config: %w", err)
	}

	poolConfig.MaxConns = config.MaxConnections
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create pool to database: %s", err.Error())
	}

	if err = pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %s", err.Error())
	}

	return pool, nil
}
