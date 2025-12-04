package services_test

import (
	"authServ/internal/domain/entities"
	"authServ/internal/services"
	email "authServ/pkg/emailSender"
	"authServ/pkg/logger"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockUsersRepository - мок репозитория пользователей
type MockUsersRepository struct {
	mock.Mock
}

func (m *MockUsersRepository) SaveUser(ctx context.Context, user *entities.User) (uuid.UUID, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockUsersRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUsersRepository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUsersRepository) UpdateUserName(ctx context.Context, id uuid.UUID, username string) error {
	args := m.Called(ctx, id, username)
	return args.Error(0)
}

func (m *MockUsersRepository) UpdateUser(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUsersRepository) UpdateVerificationCode(ctx context.Context, email, code string, sentAt time.Time) error {
	args := m.Called(ctx, email, code, sentAt)
	return args.Error(0)
}

func (m *MockUsersRepository) DeleteUserByID(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockRefreshTokensRepository - мок репозитория токенов
type MockRefreshTokensRepository struct {
	mock.Mock
}

func (m *MockRefreshTokensRepository) SaveToken(ctx context.Context, token *entities.RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockRefreshTokensRepository) GetToken(ctx context.Context, tokenStr string) (*entities.RefreshToken, error) {
	args := m.Called(ctx, tokenStr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokensRepository) DeleteToken(ctx context.Context, tokenStr string) error {
	args := m.Called(ctx, tokenStr)
	return args.Error(0)
}

// MockEmailSender - мок отправителя email
type MockEmailSender struct {
	mock.Mock
}

// Send имитирует отправку email сообщения
func (m *MockEmailSender) Send(ctx context.Context, message email.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

// Вспомогательная функция для создания тестового пользователя
func createTestUser(email string, verified bool) *entities.User {
	user, _ := entities.NewUser("testuser", email, "password123")
	user.ID = uuid.New()
	user.IsVerified = verified
	now := time.Now().Add(-2 * time.Minute)
	user.VerificationSentAt = now
	return user
}

func TestUserService_RegisterUser(t *testing.T) {
	t.Run("Успешная регистрация пользователя", func(t *testing.T) {
		mockUserRepo := new(MockUsersRepository)
		mockTokenRepo := new(MockRefreshTokensRepository)
		mockEmailSender := new(MockEmailSender)
		log, _ := logger.NewLogger()

		service := services.NewUserService(mockUserRepo, mockTokenRepo, mockEmailSender, log, "secret", 15*time.Minute)
		ctx := context.Background()

		input := services.RegisterUser{
			Username: "newuser",
			Email:    "newuser@example.com",
			Password: "password123",
		}

		userID := uuid.New()
		mockUserRepo.On("SaveUser", ctx, mock.AnythingOfType("*entities.User")).Return(userID, nil)
		mockEmailSender.On("Send", ctx, mock.Anything).Return(nil)

		err := service.RegisterUser(ctx, input)

		require.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
		mockEmailSender.AssertExpectations(t)
	})

	t.Run("Ошибка при сохранении пользователя (дубликат email)", func(t *testing.T) {
		mockUserRepo := new(MockUsersRepository)
		mockTokenRepo := new(MockRefreshTokensRepository)
		mockEmailSender := new(MockEmailSender)
		log, err := logger.NewLogger()

		service := services.NewUserService(mockUserRepo, mockTokenRepo, mockEmailSender, log, "secret", 15*time.Minute)
		ctx := context.Background()

		input := services.RegisterUser{
			Username: "testuser",
			Email:    "duplicate@example.com",
			Password: "password123",
		}

		mockUserRepo.On("SaveUser", ctx, mock.Anything).Return(uuid.Nil, errors.New("email already exists"))

		err = service.RegisterUser(ctx, input)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to save user")
		mockEmailSender.AssertNotCalled(t, "Send", mock.Anything, mock.Anything)
	})

	t.Run("Ошибка при отправке email", func(t *testing.T) {
		mockUserRepo := new(MockUsersRepository)
		mockTokenRepo := new(MockRefreshTokensRepository)
		mockEmailSender := new(MockEmailSender)
		log, err := logger.NewLogger()

		service := services.NewUserService(mockUserRepo, mockTokenRepo, mockEmailSender, log, "secret", 15*time.Minute)
		ctx := context.Background()

		input := services.RegisterUser{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}

		userID := uuid.New()
		mockUserRepo.On("SaveUser", ctx, mock.Anything).Return(userID, nil)
		mockEmailSender.On("Send", ctx, mock.Anything).Return(errors.New("SMTP connection failed"))

		err = service.RegisterUser(ctx, input)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to send verification email")
	})
}

func TestUserService_LoginUser(t *testing.T) {
	t.Run("Успешный вход", func(t *testing.T) {
		mockUserRepo := new(MockUsersRepository)
		mockTokenRepo := new(MockRefreshTokensRepository)
		mockEmailSender := new(MockEmailSender)
		log, err := logger.NewLogger()

		service := services.NewUserService(mockUserRepo, mockTokenRepo, mockEmailSender, log, "secret", 15*time.Minute)
		ctx := context.Background()

		user := createTestUser("test@example.com", true)
		err = user.SetPassword("password123")
		require.NoError(t, err)
		input := services.LoginUser{
			Email:    "test@example.com",
			Password: "password123",
		}

		mockUserRepo.On("GetUserByEmail", ctx, input.Email).Return(user, nil)
		mockTokenRepo.On("SaveToken", ctx, mock.AnythingOfType("*entities.RefreshToken")).Return(nil)

		tokens, err := service.LoginUser(ctx, input)

		require.NoError(t, err)
		assert.NotEmpty(t, tokens.AccessToken)
		assert.NotEmpty(t, tokens.RefreshToken)
	})

	t.Run("Неверный email", func(t *testing.T) {
		mockUserRepo := new(MockUsersRepository)
		mockTokenRepo := new(MockRefreshTokensRepository)
		mockEmailSender := new(MockEmailSender)
		log, err := logger.NewLogger()

		service := services.NewUserService(mockUserRepo, mockTokenRepo, mockEmailSender, log, "secret", 15*time.Minute)
		ctx := context.Background()

		input := services.LoginUser{
			Email:    "wrong@example.com",
			Password: "password123",
		}

		mockUserRepo.On("GetUserByEmail", ctx, input.Email).Return(nil, errors.New("user not found"))

		tokens, err := service.LoginUser(ctx, input)

		require.Error(t, err)
		assert.Nil(t, tokens)
		assert.Contains(t, err.Error(), "invalid email or password")
	})

	t.Run("Email не подтвержден", func(t *testing.T) {
		mockUserRepo := new(MockUsersRepository)
		mockTokenRepo := new(MockRefreshTokensRepository)
		mockEmailSender := new(MockEmailSender)
		log, err := logger.NewLogger()

		service := services.NewUserService(mockUserRepo, mockTokenRepo, mockEmailSender, log, "secret", 15*time.Minute)
		ctx := context.Background()

		user := createTestUser("test@example.com", false)
		input := services.LoginUser{
			Email:    "test@example.com",
			Password: "password123",
		}

		mockUserRepo.On("GetUserByEmail", ctx, input.Email).Return(user, nil)

		tokens, err := service.LoginUser(ctx, input)

		require.Error(t, err)
		assert.Nil(t, tokens)
		assert.Contains(t, err.Error(), "email not verified")
	})

	t.Run("Неверный пароль", func(t *testing.T) {
		mockUserRepo := new(MockUsersRepository)
		mockTokenRepo := new(MockRefreshTokensRepository)
		mockEmailSender := new(MockEmailSender)
		log, _ := logger.NewLogger()

		service := services.NewUserService(mockUserRepo, mockTokenRepo, mockEmailSender, log, "secret", 15*time.Minute)
		ctx := context.Background()

		user := createTestUser("test@example.com", true)
		err := user.SetPassword("correctpassword")
		require.NoError(t, err, "Ошибка при установке пароля в тестовом пользователе")
		input := services.LoginUser{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		mockUserRepo.On("GetUserByEmail", ctx, input.Email).Return(user, nil)

		tokens, err := service.LoginUser(ctx, input)

		require.Error(t, err)
		assert.Nil(t, tokens)
		assert.Contains(t, err.Error(), "invalid email or password")
	})
}

// === ТЕСТЫ ДЛЯ CheckTruthEmail ===

func TestUserService_CheckTruthEmail(t *testing.T) {
	t.Run("Успешная верификация", func(t *testing.T) {
		mockUserRepo := new(MockUsersRepository)
		mockTokenRepo := new(MockRefreshTokensRepository)
		mockEmailSender := new(MockEmailSender)
		log, _ := logger.NewLogger()

		service := services.NewUserService(mockUserRepo, mockTokenRepo, mockEmailSender, log, "secret", 15*time.Minute)
		ctx := context.Background()

		user := createTestUser("test@example.com", false)
		user.VerificationCode = "123456"

		input := services.CheckTruthEmail{
			Email: "test@example.com",
			Code:  "123456",
		}

		mockUserRepo.On("GetUserByEmail", ctx, input.Email).Return(user, nil)
		mockUserRepo.On("UpdateUser", ctx, mock.MatchedBy(func(u *entities.User) bool {
			return u.IsVerified && u.VerificationCode == ""
		})).Return(nil)
		mockTokenRepo.On("SaveToken", ctx, mock.Anything).Return(nil)

		tokens, err := service.CheckTruthEmail(ctx, input)

		require.NoError(t, err)
		assert.NotEmpty(t, tokens.AccessToken)
		assert.NotEmpty(t, tokens.RefreshToken)
	})

	t.Run("Пользователь не найден", func(t *testing.T) {
		mockUserRepo := new(MockUsersRepository)
		mockTokenRepo := new(MockRefreshTokensRepository)
		mockEmailSender := new(MockEmailSender)
		log, _ := logger.NewLogger()

		service := services.NewUserService(mockUserRepo, mockTokenRepo, mockEmailSender, log, "secret", 15*time.Minute)
		ctx := context.Background()

		input := services.CheckTruthEmail{
			Email: "notfound@example.com",
			Code:  "123456",
		}

		mockUserRepo.On("GetUserByEmail", ctx, input.Email).Return(nil, errors.New("user not found"))

		tokens, err := service.CheckTruthEmail(ctx, input)

		require.Error(t, err)
		assert.Empty(t, tokens.AccessToken)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Пользователь уже верифицирован", func(t *testing.T) {
		mockUserRepo := new(MockUsersRepository)
		mockTokenRepo := new(MockRefreshTokensRepository)
		mockEmailSender := new(MockEmailSender)
		log, _ := logger.NewLogger()

		service := services.NewUserService(mockUserRepo, mockTokenRepo, mockEmailSender, log, "secret", 15*time.Minute)
		ctx := context.Background()

		user := createTestUser("test@example.com", true)

		input := services.CheckTruthEmail{
			Email: "test@example.com",
			Code:  "123456",
		}

		mockUserRepo.On("GetUserByEmail", ctx, input.Email).Return(user, nil)

		tokens, err := service.CheckTruthEmail(ctx, input)

		require.Error(t, err)
		assert.Empty(t, tokens.AccessToken)
		assert.Contains(t, err.Error(), "already verified")
	})

	t.Run("Неверный код верификации", func(t *testing.T) {
		mockUserRepo := new(MockUsersRepository)
		mockTokenRepo := new(MockRefreshTokensRepository)
		mockEmailSender := new(MockEmailSender)
		log, _ := logger.NewLogger()

		service := services.NewUserService(mockUserRepo, mockTokenRepo, mockEmailSender, log, "secret", 15*time.Minute)
		ctx := context.Background()

		user := createTestUser("test@example.com", false)
		user.VerificationCode = "123456"

		input := services.CheckTruthEmail{
			Email: "test@example.com",
			Code:  "wrong-code",
		}

		mockUserRepo.On("GetUserByEmail", ctx, input.Email).Return(user, nil)

		tokens, err := service.CheckTruthEmail(ctx, input)

		require.Error(t, err)
		assert.Empty(t, tokens.AccessToken)
		assert.Contains(t, err.Error(), "invalid verification code")
	})
}

// === ТЕСТЫ ДЛЯ ResendVerificationEmail ===

func TestUserService_ResendVerificationEmail(t *testing.T) {
	t.Run("Успешная повторная отправка", func(t *testing.T) {
		mockUserRepo := new(MockUsersRepository)
		mockTokenRepo := new(MockRefreshTokensRepository)
		mockEmailSender := new(MockEmailSender)
		log, _ := logger.NewLogger()

		service := services.NewUserService(mockUserRepo, mockTokenRepo, mockEmailSender, log, "secret", 15*time.Minute)
		ctx := context.Background()

		user := createTestUser("test@example.com", false)
		pastTime := time.Now().Add(-2 * time.Minute)
		user.VerificationSentAt = pastTime

		input := services.ResendVerificationEmailInput{
			Email: "test@example.com",
		}

		mockUserRepo.On("GetUserByEmail", ctx, input.Email).Return(user, nil)
		mockUserRepo.On("UpdateVerificationCode", ctx, input.Email, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(nil)
		mockEmailSender.On("Send", ctx, mock.Anything).Return(nil)

		err := service.ResendVerificationEmail(ctx, input)

		require.NoError(t, err)
		mockEmailSender.AssertCalled(t, "Send", ctx, mock.Anything)
	})

	t.Run("Слишком рано для повторной отправки", func(t *testing.T) {
		mockUserRepo := new(MockUsersRepository)
		mockTokenRepo := new(MockRefreshTokensRepository)
		mockEmailSender := new(MockEmailSender)
		log, _ := logger.NewLogger()

		service := services.NewUserService(mockUserRepo, mockTokenRepo, mockEmailSender, log, "secret", 15*time.Minute)
		ctx := context.Background()

		user := createTestUser("test@example.com", false)
		recentTime := time.Now().Add(-30 * time.Second)
		user.VerificationSentAt = recentTime

		input := services.ResendVerificationEmailInput{
			Email: "test@example.com",
		}

		mockUserRepo.On("GetUserByEmail", ctx, input.Email).Return(user, nil)

		err := service.ResendVerificationEmail(ctx, input)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "please wait")
		mockEmailSender.AssertNotCalled(t, "Send", mock.Anything, mock.Anything)
	})

	t.Run("Пользователь уже верифицирован", func(t *testing.T) {
		mockUserRepo := new(MockUsersRepository)
		mockTokenRepo := new(MockRefreshTokensRepository)
		mockEmailSender := new(MockEmailSender)
		log, _ := logger.NewLogger()

		service := services.NewUserService(mockUserRepo, mockTokenRepo, mockEmailSender, log, "secret", 15*time.Minute)
		ctx := context.Background()

		user := createTestUser("test@example.com", true)

		input := services.ResendVerificationEmailInput{
			Email: "test@example.com",
		}

		mockUserRepo.On("GetUserByEmail", ctx, input.Email).Return(user, nil)

		err := service.ResendVerificationEmail(ctx, input)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "already verified")
	})
}

func TestUserService_RefreshTokens(t *testing.T) {
	t.Run("Успешное обновление токенов", func(t *testing.T) {
		mockUserRepo := new(MockUsersRepository)
		mockTokenRepo := new(MockRefreshTokensRepository)
		mockEmailSender := new(MockEmailSender)
		log, _ := logger.NewLogger()

		service := services.NewUserService(mockUserRepo, mockTokenRepo, mockEmailSender, log, "secret", 15*time.Minute)
		ctx := context.Background()

		userID := uuid.New()
		user := createTestUser("test@example.com", true)
		user.ID = userID

		refreshToken := &entities.RefreshToken{
			UserID:    userID,
			Token:     "valid-refresh-token",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}

		input := services.RefreshInput{
			RefreshToken: "valid-refresh-token",
		}

		mockTokenRepo.On("GetToken", ctx, input.RefreshToken).Return(refreshToken, nil)
		mockUserRepo.On("GetUserByID", ctx, userID).Return(user, nil)
		mockTokenRepo.On("DeleteToken", ctx, input.RefreshToken).Return(nil)
		mockTokenRepo.On("SaveToken", ctx, mock.Anything).Return(nil)

		tokens, err := service.RefreshTokens(ctx, input)

		require.NoError(t, err)
		assert.NotEmpty(t, tokens.AccessToken)
		assert.NotEmpty(t, tokens.RefreshToken)
	})

	t.Run("Недействительный refresh token", func(t *testing.T) {
		mockUserRepo := new(MockUsersRepository)
		mockTokenRepo := new(MockRefreshTokensRepository)
		mockEmailSender := new(MockEmailSender)
		log, _ := logger.NewLogger()

		service := services.NewUserService(mockUserRepo, mockTokenRepo, mockEmailSender, log, "secret", 15*time.Minute)
		ctx := context.Background()

		input := services.RefreshInput{
			RefreshToken: "invalid-token",
		}

		mockTokenRepo.On("GetToken", ctx, input.RefreshToken).Return(nil, errors.New("token not found"))

		tokens, err := service.RefreshTokens(ctx, input)

		require.Error(t, err)
		assert.Empty(t, tokens.AccessToken)
		assert.Contains(t, err.Error(), "invalid or expired refresh token")
	})

	t.Run("Истек срок действия refresh token", func(t *testing.T) {
		mockUserRepo := new(MockUsersRepository)
		mockTokenRepo := new(MockRefreshTokensRepository)
		mockEmailSender := new(MockEmailSender)
		log, _ := logger.NewLogger()

		service := services.NewUserService(mockUserRepo, mockTokenRepo, mockEmailSender, log, "secret", 15*time.Minute)
		ctx := context.Background()

		refreshToken := &entities.RefreshToken{
			UserID:    uuid.New(),
			Token:     "expired-token",
			ExpiresAt: time.Now().Add(-1 * time.Hour),
		}

		input := services.RefreshInput{
			RefreshToken: "expired-token",
		}

		mockTokenRepo.On("GetToken", ctx, input.RefreshToken).Return(refreshToken, nil)
		mockTokenRepo.On("DeleteToken", ctx, input.RefreshToken).Return(nil)

		tokens, err := service.RefreshTokens(ctx, input)

		require.Error(t, err)
		assert.Empty(t, tokens.AccessToken)
		assert.Contains(t, err.Error(), "refresh token expired")
	})
}

// === ТЕСТЫ ДЛЯ Logout ===

func TestUserService_Logout(t *testing.T) {
	t.Run("Успешный logout", func(t *testing.T) {
		mockUserRepo := new(MockUsersRepository)
		mockTokenRepo := new(MockRefreshTokensRepository)
		mockEmailSender := new(MockEmailSender)
		log, _ := logger.NewLogger()

		service := services.NewUserService(mockUserRepo, mockTokenRepo, mockEmailSender, log, "secret", 15*time.Minute)
		ctx := context.Background()

		input := services.RefreshInput{
			RefreshToken: "some-token",
		}

		mockTokenRepo.On("DeleteToken", ctx, input.RefreshToken).Return(nil)

		err := service.Logout(ctx, input)

		require.NoError(t, err)
		mockTokenRepo.AssertCalled(t, "DeleteToken", ctx, input.RefreshToken)
	})

	t.Run("Ошибка при удалении токена", func(t *testing.T) {
		mockUserRepo := new(MockUsersRepository)
		mockTokenRepo := new(MockRefreshTokensRepository)
		mockEmailSender := new(MockEmailSender)
		log, _ := logger.NewLogger()

		service := services.NewUserService(mockUserRepo, mockTokenRepo, mockEmailSender, log, "secret", 15*time.Minute)
		ctx := context.Background()

		input := services.RefreshInput{
			RefreshToken: "some-token",
		}

		mockTokenRepo.On("DeleteToken", ctx, input.RefreshToken).Return(errors.New("token not found"))

		err := service.Logout(ctx, input)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete refresh token")
	})
}

// === ТЕСТЫ ДЛЯ ParseAccessToken ===

func TestUserService_ParseAccessToken(t *testing.T) {
	t.Run("Валидный access token", func(t *testing.T) {
		mockUserRepo := new(MockUsersRepository)
		mockTokenRepo := new(MockRefreshTokensRepository)
		mockEmailSender := new(MockEmailSender)
		log, _ := logger.NewLogger()

		secretKey := "test-secret-key"
		service := services.NewUserService(mockUserRepo, mockTokenRepo, mockEmailSender, log, secretKey, 15*time.Minute)
		ctx := context.Background()

		userID := uuid.New()
		claims := jwt.MapClaims{
			"user_id": userID.String(),
			"email":   "test@example.com",
			"exp":     time.Now().Add(15 * time.Minute).Unix(),
			"iat":     time.Now().Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(secretKey))

		parsedUserID, err := service.ParseAccessToken(ctx, tokenString)

		require.NoError(t, err)
		assert.Equal(t, userID, parsedUserID)
	})

	t.Run("Невалидный токен", func(t *testing.T) {
		mockUserRepo := new(MockUsersRepository)
		mockTokenRepo := new(MockRefreshTokensRepository)
		mockEmailSender := new(MockEmailSender)
		log, _ := logger.NewLogger()

		service := services.NewUserService(mockUserRepo, mockTokenRepo, mockEmailSender, log, "secret", 15*time.Minute)
		ctx := context.Background()

		userID, err := service.ParseAccessToken(ctx, "invalid.token.here")

		require.Error(t, err)
		assert.Equal(t, uuid.Nil, userID)
	})

	t.Run("Истекший токен", func(t *testing.T) {
		mockUserRepo := new(MockUsersRepository)
		mockTokenRepo := new(MockRefreshTokensRepository)
		mockEmailSender := new(MockEmailSender)
		log, _ := logger.NewLogger()

		secretKey := "test-secret-key"
		service := services.NewUserService(mockUserRepo, mockTokenRepo, mockEmailSender, log, secretKey, 15*time.Minute)
		ctx := context.Background()

		claims := jwt.MapClaims{
			"user_id": uuid.New().String(),
			"email":   "test@example.com",
			"exp":     time.Now().Add(-1 * time.Hour).Unix(),
			"iat":     time.Now().Add(-2 * time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(secretKey))

		userID, err := service.ParseAccessToken(ctx, tokenString)

		require.Error(t, err)
		assert.Equal(t, uuid.Nil, userID)
	})
}
