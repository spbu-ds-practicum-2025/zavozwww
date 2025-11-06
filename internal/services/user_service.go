package services

import (
	"authServ/internal/domain/entities"
	repositories "authServ/internal/repositories/postgres"
	email "authServ/pkg/emailSender"
	"authServ/pkg/logger"
	"context"
	"time"

	"github.com/google/uuid"
)

// RegisterUser используется для регистрации нового пользователя.
type RegisterUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginUser используется для аутентификации пользователя.
type LoginUser struct {
	Email    string                `json:"email"`
	Password string                `json:"password"`
	RefToken entities.RefreshToken `json:"ref_token"`
}

// UpdateUserName используется для обновления имени пользователя.
type UpdateUserName struct {
	UserID      string `json:"-"`
	NewUsername string `json:"new_username"`
}

// CheckTruthEmail используется для проверки подлинности электронной почты пользователя.
type CheckTruthEmail struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

// Tokens представляет собой структуру для хранения пар токенов.
type Tokens struct {
	AccessToken  string
	RefreshToken string
}

// UserService определяет методы для работы с пользователями и токенами.
type userService struct {
	userRepo       repositories.UsersRepository
	tokenRepo      repositories.RefreshTokensRepository
	emailSender    email.EmailSender
	log            *logger.Logger
	secretKey      string
	accessTokenTTL time.Duration
}

// UserService определяет методы для работы с пользователями.
type UserService interface {
	RegisterUser(ctx context.Context, input RegisterUser) (uuid.UUID, error)
	LoginUser(input LoginUser) (Tokens, error)
	CheckTruthEmail(ctx context.Context, input CheckTruthEmail) (Tokens, error)
}

// RegisterUser регистрирует нового пользователя.
func (s *userService) RegisterUser(ctx context.Context, input RegisterUser) (uuid.UUID, error) {
	return uuid.Nil, nil
}

// LoginUser аутентифицирует пользователя и возвращает токены.
func (s *userService) LoginUser(input LoginUser) (Tokens, error) {
	return Tokens{}, nil
}

// CheckTruthEmail отправляет код подтверждения на указанный email.
func (s *userService) CheckTruthEmail(ctx context.Context, input CheckTruthEmail) (Tokens, error) {
	return Tokens{}, nil
}

// NewUserService создает новый экземпляр UserService.
func NewUserService(
	userRepo repositories.UsersRepository,
	tokenRepo repositories.RefreshTokensRepository,
	emailSender email.EmailSender,
	log *logger.Logger,
	secretKey string,
	tokenTTL time.Duration,
) UserService {
	return &userService{
		userRepo:       userRepo,
		tokenRepo:      tokenRepo,
		emailSender:    emailSender,
		log:            log,
		secretKey:      secretKey,
		accessTokenTTL: tokenTTL,
	}
}
