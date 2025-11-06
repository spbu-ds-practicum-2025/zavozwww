package repositories

import (
	"authServ/internal/domain/entities"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository предоставляет методы для работы с пользователями в базе данных PostgreSQL.
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository создает новый экземпляр UserRepository.
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// Проверка на соответствие интерфейсу во время компиляции.
var _ UsersRepository = (*UserRepository)(nil)

// GetUserByID получает пользователя по его ID.
func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	const op = "repositories.postgres.UserRepository.GetUserByID"

	query := `SELECT id, username, email, password_hash, created_at FROM users WHERE id = $1`

	row := r.db.QueryRow(ctx, query, id)

	var user entities.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

// SaveUser сохраняет нового пользователя в базе данных и возвращает его ID.
func (r *UserRepository) SaveUser(ctx context.Context, user *entities.User) (uuid.UUID, error) {
	const op = "repositories.postgres.UserRepository.SaveUser"

	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id`

	row := r.db.QueryRow(ctx, query, user.Username, user.Email, user.PasswordHash)

	var newID uuid.UUID
	err := row.Scan(&newID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_email_key":
				return uuid.Nil, fmt.Errorf("%s: %w", op, ErrEmailExists)
			case "users_username_key":
				return uuid.Nil, fmt.Errorf("%s: %w", op, ErrUsernameExists)
			default:
				return uuid.Nil, fmt.Errorf("%s: %w", op, ErrUserAlreadyExists)
			}
		}
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	return newID, nil
}

// UpdateUserName обновляет имя пользователя по его ID.
func (r *UserRepository) UpdateUserName(ctx context.Context, id uuid.UUID, newUsername string) error {
	const op = "repositories.postgres.UserRepository.UpdateUserName"
	query := `UPDATE users SET username = $1 WHERE id = $2`
	cmdTag, err := r.db.Exec(ctx, query, newUsername, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, ErrUsernameExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}
	return nil
}

// DeleteUserByID удаляет пользователя по его ID.
func (r *UserRepository) DeleteUserByID(ctx context.Context, id uuid.UUID) error {
	const op = "repositories.postgres.UserRepository.DeleteUserByID"
	query := `DELETE FROM users WHERE id = $1`
	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}
	return nil
}

