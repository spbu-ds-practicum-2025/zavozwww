package repositories

import (
	"authServ/internal/domain/entities"
	"context"
	"errors"
	"fmt"
	"time" // <-- Добавлен импорт time

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

	// ИСПРАВЛЕНИЕ: Преобразуем created_at в текст, чтобы соответствовать полю CreatedAt (string) в структуре.
	// Остальные поля, включая verification_sent_at, не трогаем.
	query := `SELECT id, username, email, password_hash, created_at::text, is_verified, verification_code, verification_sent_at 
               FROM users 
               WHERE id = $1`

	row := r.db.QueryRow(ctx, query, id)

	var user entities.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.IsVerified,
		&user.VerificationCode,
		&user.VerificationSentAt,
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
	user.IsVerified = false
	query := `INSERT INTO users (username, email, password_hash, is_verified, verification_code, verification_sent_at) 
              VALUES ($1, $2, $3, $4, $5, $6) 
              RETURNING id`

	row := r.db.QueryRow(ctx, query, user.Username, user.Email, user.PasswordHash, user.IsVerified, user.VerificationCode, user.VerificationSentAt)

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

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	const op = "repositories.postgres.UserRepository.GetUserByEmail"

	query := `SELECT id, username, email, password_hash, created_at::text, is_verified, verification_code, verification_sent_at 
              FROM users 
              WHERE email = $1`

	row := r.db.QueryRow(ctx, query, email)
	var user entities.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.IsVerified,
		&user.VerificationCode,
		&user.VerificationSentAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Добавим лог, чтобы было видно, что пользователь не найден на этом этапе
			fmt.Printf("--- ОТЛАДКА: Пользователь с email '%s' не найден основным запросом (pgx.ErrNoRows).\n", email)
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

// ИСПРАВЛЕННЫЙ РЕСИВЕР: Теперь это метод UserRepository, как и должно быть.
func (r *UserRepository) UpdateUser(ctx context.Context, user *entities.User) error {
	const op = "repositories.postgres.UserRepository.UpdateUser"
	// ИСПРАВЛЕНИЕ: Добавлено поле verification_sent_at
	query := `UPDATE users SET 
                username = $1, 
                email = $2, 
                password_hash = $3, 
                is_verified = $4, 
                verification_code = $5,
                verification_sent_at = $6
              WHERE id = $7`

	cmdTag, err := r.db.Exec(ctx, query, user.Username, user.Email, user.PasswordHash, user.IsVerified, user.VerificationCode, user.VerificationSentAt, user.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}
	return nil
}

// UpdateVerificationCode обновляет код верификации и время его отправки.
func (r *UserRepository) UpdateVerificationCode(ctx context.Context, email, code string, sentAt time.Time) error {
	const op = "repositories.postgres.UserRepository.UpdateVerificationCode"
	query := `UPDATE users SET verification_code = $1, verification_sent_at = $2 WHERE email = $3`
	cmdTag, err := r.db.Exec(ctx, query, code, sentAt, email)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}
	return nil
}
