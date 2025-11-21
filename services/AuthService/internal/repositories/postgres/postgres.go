package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Config для подключения к базе данных.
type PgConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLMode  string
}

// NewPostgresStorage создает новый пул соединений с PostgreSQL.
func NewPostgresStorage(ctx context.Context, cfg PgConfig) (*pgxpool.Pool, error) {
	const op = "repositories.postgres.NewPostgresStorage"

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create connection pool: %w", op, err)
	}

	maxAttempts := 5
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err = pool.Ping(ctx)
		if err == nil {
			return pool, nil
		}
		fmt.Printf("Attempt %d: failed to connect to database: %v. Retrying in 5 seconds...\n", attempt, err)
		time.Sleep(5 * time.Second)
	}

	return nil, fmt.Errorf("%s: failed to connect to database after %d attempts: %w", op, maxAttempts, err)
}
