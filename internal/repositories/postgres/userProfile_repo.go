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

// UserProfileRepository предоставляет методы для работы с профилями пользователей.
type UserProfileRepository struct {
	db *pgxpool.Pool
}

// NewUserProfileRepository создает новый экземпляр UserProfileRepository.
func NewUserProfileRepository(db *pgxpool.Pool) *UserProfileRepository {
	return &UserProfileRepository{db: db}
}

// Проверка на соответствие интерфейсу во время компиляции.
var _ UserProfilesRepository = (*UserProfileRepository)(nil)

// SaveProfile сохраняет новый профиль пользователя.
func (r *UserProfileRepository) SaveProfile(ctx context.Context, profile *entities.UserProfile) error {
	const op = "repositories.postgres.UserProfileRepository.SaveProfile"
	query := `INSERT INTO user_profiles (user_id, first_name, last_name, age, info, city) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(ctx, query, profile.UserId, profile.FirstName, profile.LastName, profile.Age, profile.Info, profile.City)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, ErrProfileAlreadyExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// GetProfileByUserID получает профиль пользователя по его ID.
func (r *UserProfileRepository) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (*entities.UserProfile, error) {
	const op = "repositories.postgres.UserProfileRepository.GetProfileByUserID"
	query := `SELECT user_id, first_name, last_name, age, info, city FROM user_profiles WHERE user_id = $1`

	row := r.db.QueryRow(ctx, query, userID)
	var profile entities.UserProfile
	err := row.Scan(&profile.UserId, &profile.FirstName, &profile.LastName, &profile.Age, &profile.Info, &profile.City)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserProfileNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &profile, nil
}

// UpdateProfile обновляет существующий профиль пользователя.
func (r *UserProfileRepository) UpdateProfile(ctx context.Context, profile *entities.UserProfile) error {
	const op = "repositories.postgres.UserProfileRepository.UpdateProfile"
	query := `UPDATE user_profiles SET first_name = $1, last_name = $2, age = $3, info = $4, city = $5 WHERE user_id = $6`

	cmdTag, err := r.db.Exec(ctx, query, profile.FirstName, profile.LastName, profile.Age, profile.Info, profile.City, profile.UserId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, ErrUserProfileNotFound)
	}

	return nil
}

// DeleteProfile удаляет профиль пользователя по его ID.
func (r *UserProfileRepository) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	const op = "repositories.postgres.UserProfileRepository.DeleteProfile"
	query := `DELETE FROM user_profiles WHERE user_id = $1`

	cmdTag, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, ErrUserProfileNotFound)
	}

	return nil
}
