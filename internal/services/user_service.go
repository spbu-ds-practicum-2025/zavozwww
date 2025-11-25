package services

import (
	"authServ/internal/domain/entities"
	repositories "authServ/internal/repositories/postgres"
	email "authServ/pkg/emailSender"
	"authServ/pkg/logger"
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

// RefreshInput определяет структуру для запроса на обновление токенов.
type RefreshInput struct {
	RefreshToken string `json:"refresh_token"`
}

// ResendVerificationEmailInput определяет структуру для запроса на повторную отправку письма.
type ResendVerificationEmailInput struct {
	Email string `json:"email"`
}

// tokenClaims определяет структуру данных (полезной нагрузки) для JWT токена.
type tokenClaims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
}

// Tokens представляет собой структуру для хранения пар токенов.
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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
	RegisterUser(ctx context.Context, input RegisterUser) error
	LoginUser(ctx context.Context, input LoginUser) (*Tokens, error)
	CheckTruthEmail(ctx context.Context, input CheckTruthEmail) (Tokens, error)
	ResendVerificationEmail(ctx context.Context, input ResendVerificationEmailInput) error
	RefreshTokens(ctx context.Context, input RefreshInput) (Tokens, error)
	Logout(ctx context.Context, input RefreshInput) error
	ParseAccessToken(ctx context.Context, accessToken string) (uuid.UUID, error)
}

// RegisterUser регистрирует нового пользователя, сохраняет его в БД и отправляет письмо для верификации.
func (s *userService) RegisterUser(ctx context.Context, input RegisterUser) error {
	const op = "services.userService.RegisterUser"

	user, err := entities.NewUser(input.Username, input.Email, input.Password)

	if err != nil {
		return fmt.Errorf("%s: failed to create user entity: %w", op, err)
	}

	user.VerificationSentAt = time.Now()

	if _, err := s.userRepo.SaveUser(ctx, user); err != nil {
		return fmt.Errorf("%s: failed to save user: %w", op, err)
	}

	msg := email.Message{
		To:           []string{user.Email},
		Subject:      "Подтверждение регистрации",
		TemplateName: "verification.html",
		TemplateData: map[string]interface{}{
			"Code":    user.VerificationCode,
			"Subject": "Подтверждение регистрации",
			"Year":    time.Now().Year(),
		},
	}

	if err := s.emailSender.Send(ctx, msg); err != nil {
		return fmt.Errorf("%s: failed to send verification email: %w", op, err)
	}

	s.log.Info("user registered and verification email sent", "email", user.Email)
	return nil
}

// LoginUser аутентифицирует пользователя и возвращает токены.
func (s *userService) LoginUser(ctx context.Context, input LoginUser) (*Tokens, error) {
	const op = "services.userService.LoginUser"
	user, err := s.userRepo.GetUserByEmail(ctx, input.Email)

	if err != nil {
		return nil, fmt.Errorf("%s: invalid email or password", op)
	}

	if !user.IsVerified {
		return nil, fmt.Errorf("%s: email not verified", op)
	}

	if !user.ComparePassword(input.Password) {
		return nil, fmt.Errorf("%s: invalid email or password", op)
	}

	accessToken, err := s.generateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to generate access token: %w", op, err)
	}

	refreshToken, err := s.generateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to generate refresh token: %w", op, err)
	}

	s.log.Info("user logged in successfully", "email", user.Email)

	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// CheckTruthEmail проверяет код верификации, активирует пользователя и выдает токены.
func (s *userService) CheckTruthEmail(ctx context.Context, input CheckTruthEmail) (Tokens, error) {
	const op = "services.userService.CheckTruthEmail"
	user, err := s.userRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return Tokens{}, fmt.Errorf("%s: user with email %s not found", op, input.Email)
	}

	if user.IsVerified {
		return Tokens{}, fmt.Errorf("%s: user is already verified", op)
	}

	if user.VerificationCode != input.Code {
		return Tokens{}, fmt.Errorf("%s: invalid verification code", op)
	}

	user.IsVerified = true
	user.VerificationCode = ""
	user.VerificationSentAt = time.Time{} // Очищаем время после успешной верификации

	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return Tokens{}, fmt.Errorf("%s: failed to update user status: %w", op, err)
	}

	// 5. Генерируем и возвращаем пару токенов для автоматического входа.
	accessToken, err := s.generateAccessToken(user.ID, user.Email)
	if err != nil {
		return Tokens{}, fmt.Errorf("%s: failed to generate access token: %w", op, err)
	}

	refreshToken, err := s.generateRefreshToken(ctx, user.ID)
	if err != nil {
		return Tokens{}, fmt.Errorf("%s: failed to generate refresh token: %w", op, err)
	}

	s.log.Info("user successfully verified", "email", user.Email)

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ResendVerificationEmail проверяет, можно ли отправить письмо повторно, и если да, отправляет его.
func (s *userService) ResendVerificationEmail(ctx context.Context, input ResendVerificationEmailInput) error {
	const op = "services.userService.ResendVerificationEmail"
	const waitDuration = 60 * time.Second

	user, err := s.userRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return fmt.Errorf("%s: failed to process request", op)
	}

	if user.IsVerified {
		return fmt.Errorf("%s: user is already verified", op)
	}

	if time.Since(user.VerificationSentAt) < waitDuration {
		remaining := waitDuration - time.Since(user.VerificationSentAt)
		return fmt.Errorf("%s: please wait %.0f seconds before requesting a new code", op, remaining.Seconds())
	}

	newVerificationCode, err := entities.GenerateVerificationCode()
	if err != nil {
		return fmt.Errorf("%s: failed to generate verification code: %w", op, err)
	}
	now := time.Now()

	if err := s.userRepo.UpdateVerificationCode(ctx, user.Email, newVerificationCode, now); err != nil {
		return fmt.Errorf("%s: failed to update verification info: %w", op, err)
	}

	msg := email.Message{
		To:           []string{user.Email},
		Subject:      "Новый код подтверждения",
		TemplateName: "verification.html",
		TemplateData: map[string]interface{}{
			"Code":    newVerificationCode,
			"Subject": "Новый код подтверждения",
			"Year":    time.Now().Year(),
		},
	}

	if err := s.emailSender.Send(ctx, msg); err != nil {
		return fmt.Errorf("%s: failed to send verification email: %w", op, err)
	}

	s.log.Info("resent verification email", "email", user.Email)
	return nil
}

// RefreshTokens проверяет refresh токен и выдает новую пару токенов.
func (s *userService) RefreshTokens(ctx context.Context, input RefreshInput) (Tokens, error) {
	const op = "services.userService.RefreshTokens"
	token, err := s.tokenRepo.GetToken(ctx, input.RefreshToken)

	if err != nil {
		return Tokens{}, fmt.Errorf("%s: invalid or expired refresh token", op)
	}

	if time.Now().After(token.ExpiresAt) {
		_ = s.tokenRepo.DeleteToken(ctx, input.RefreshToken)
		return Tokens{}, fmt.Errorf("%s: refresh token expired", op)
	}

	user, err := s.userRepo.GetUserByID(ctx, token.UserID)
	if err != nil {
		return Tokens{}, fmt.Errorf("%s: user not found for token: %w", op, err)
	}

	if err := s.tokenRepo.DeleteToken(ctx, input.RefreshToken); err != nil {
		return Tokens{}, fmt.Errorf("%s: failed to delete old refresh token: %w", op, err)
	}

	newAccessToken, err := s.generateAccessToken(user.ID, user.Email)
	if err != nil {
		return Tokens{}, fmt.Errorf("%s: failed to generate access token: %w", op, err)
	}

	newRefreshToken, err := s.generateRefreshToken(ctx, user.ID)
	if err != nil {
		return Tokens{}, fmt.Errorf("%s: failed to generate refresh token: %w", op, err)
	}

	s.log.Info("tokens refreshed successfully", "user_id", user.ID)

	return Tokens{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// Logout завершает сессию пользователя, удаляя его refresh токен.
func (s *userService) Logout(ctx context.Context, input RefreshInput) error {
	const op = "services.userService.Logout"
	if err := s.tokenRepo.DeleteToken(ctx, input.RefreshToken); err != nil {
		return fmt.Errorf("%s: failed to delete refresh token: %w", op, err)
	}
	s.log.Info("user logged out successfully")
	return nil
}

// generateAccessToken создает новый JWT access токен для пользователя.
func (s *userService) generateAccessToken(userID uuid.UUID, email string) (string, error) {
	const op = "services.userService.generateAccessToken"
	claims := tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,
		Email:  email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", fmt.Errorf("%s: failed to sign access token: %w", op, err)
	}

	return signedToken, nil
}

// generateRefreshToken создает новый refresh токен, сохраняет его в БД и возвращает.
func (s *userService) generateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	const op = "services.userService.generateRefreshToken"
	refreshToken := entities.RefreshToken{
		UserID:    userID,
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	if err := s.tokenRepo.SaveToken(ctx, &refreshToken); err != nil {
		return "", fmt.Errorf("%s: failed to save refresh token: %w", op, err)
	}

	return refreshToken.Token, nil
}

// ParseAccessToken проверяет access токен и возвращает ID пользователя из него.
func (s *userService) ParseAccessToken(ctx context.Context, accessToken string) (uuid.UUID, error) {
	const op = "services.userService.ParseAccessToken"

	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	if claims, ok := token.Claims.(*tokenClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return uuid.Nil, fmt.Errorf("%s: invalid token", op)
}

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
