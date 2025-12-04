package repositories

import (
	"authServ/internal/domain/entities"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepository struct {
	db *pgxpool.Pool
}

// NewRefreshTokenRepository создает новый экземпляр RefreshTokenRepository.
func NewRefreshTokenRepository(db *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// Проверка на соответствие интерфейсу во время компиляции.
var _ RefreshTokensRepository = (*RefreshTokenRepository)(nil)

// SaveToken сохраняет новый токен обновления.
func (r *RefreshTokenRepository) SaveToken(ctx context.Context, token *entities.RefreshToken) error {
	const op = "repositories.postgres.RefreshTokenRepository.SaveToken"
	query := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, token.UserID, token.Token, token.ExpiresAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, ErrTokenAlreadyExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// GetToken получает токен по его строковому представлению.
func (r *RefreshTokenRepository) GetToken(ctx context.Context, tokenString string) (*entities.RefreshToken, error) {
	const op = "repositories.postgres.RefreshTokenRepository.GetToken"
	query := `SELECT id, user_id, token, expires_at FROM refresh_tokens WHERE token = $1`
	row := r.db.QueryRow(ctx, query, tokenString)

	var token entities.RefreshToken
	err := row.Scan(&token.ID, &token.UserID, &token.Token, &token.ExpiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrTokenNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &token, nil
}

// DeleteToken удаляет токен по его строковому представлению.
func (r *RefreshTokenRepository) DeleteToken(ctx context.Context, tokenString string) error {
	const op = "repositories.postgres.RefreshTokenRepository.DeleteToken"
	query := `DELETE FROM refresh_tokens WHERE token = $1`
	cmdTag, err := r.db.Exec(ctx, query, tokenString)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, ErrTokenNotFound)
	}
	return nil
}
