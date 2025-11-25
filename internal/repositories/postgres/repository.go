package repositories

import (
	"authServ/internal/domain/entities"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Общие, но более специфичные ошибки для всех репозиториев.
var (
	// Ошибки "не найдено"
	ErrUserNotFound        = errors.New("user not found")
	ErrUserProfileNotFound = errors.New("user profile not found")
	ErrTokenNotFound       = errors.New("token not found")

	// Ошибки "уже существует" (общие)
	ErrUserAlreadyExists    = errors.New("user with such email or username already exists")
	ErrProfileAlreadyExists = errors.New("user profile already exists")
	ErrTokenAlreadyExists   = errors.New("token already exists")

	// Ошибки "уже существует" (конкретные для пользователя)
	ErrEmailExists    = errors.New("user with this email already exists")
	ErrUsernameExists = errors.New("user with this username already exists")
)

// UsersRepository определяет методы для работы с пользователями.
type UsersRepository interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	SaveUser(ctx context.Context, user *entities.User) (uuid.UUID, error)
	UpdateUserName(ctx context.Context, id uuid.UUID, newUsername string) error
	DeleteUserByID(ctx context.Context, id uuid.UUID) error
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) error
	UpdateVerificationCode(ctx context.Context, email, code string, sentAt time.Time) error
}

// UserProfilesRepository определяет методы для работы с профилями.
type UserProfilesRepository interface {
	SaveProfile(ctx context.Context, profile *entities.UserProfile) error
	GetProfileByUserID(ctx context.Context, userID uuid.UUID) (*entities.UserProfile, error)
	UpdateProfile(ctx context.Context, profile *entities.UserProfile) error
	DeleteProfile(ctx context.Context, userID uuid.UUID) error
}

// RefreshTokensRepository определяет методы для работы с токенами.
type RefreshTokensRepository interface {
	SaveToken(ctx context.Context, token *entities.RefreshToken) error
	GetToken(ctx context.Context, tokenString string) (*entities.RefreshToken, error)
	DeleteToken(ctx context.Context, tokenString string) error
}
