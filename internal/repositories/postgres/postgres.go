package repositories

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PgConfig содержит параметры подключения к базе данных PostgreSQL.
type PgConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLMode  string
}

// NewPostgresStorage создает новый пул соединений с PostgreSQL с повторными попытками подключения.
func NewPostgresStorage(ctx context.Context, cfg PgConfig) (*pgxpool.Pool, error) {
	const op = "repositories.postgres.NewPostgresStorage"

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create connection pool: %w", op, err)
	}

	const maxAttempts = 5
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err = pool.Ping(ctx); err == nil {
			return pool, nil
		}
		fmt.Printf("[%s] Attempt %d: failed to connect to database: %v. Retrying in 5 seconds...\n", op, attempt, err)
		time.Sleep(5 * time.Second)
	}

	return nil, fmt.Errorf("%s: failed to connect to database after %d attempts: %w", op, maxAttempts, err)
}

// GetConfig заполняет PgConfig из переменных окружения и валидирует их.
func (p *PgConfig) GetConfig() error {
	const op = "repositories.postgres.PgConfig.GetConfig"

	var err error
	if p.User, err = getEnv("DB_USER"); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if p.Password, err = getEnv("POSTGRES_PASSWORD"); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if p.Host, err = getEnv("POSTGRES_HOST"); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if p.Port, err = getEnv("POSTGRES_PORT"); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if p.DBName, err = getEnv("POSTGRES_DB"); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if p.SSLMode, err = getEnv("POSTGRES_SSLMODE"); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func getEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("%s is not set", key)
	}
	return val, nil
}
