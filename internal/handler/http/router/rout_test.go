package rout_test

import (
	rout "authServ/internal/handler/http/router"
	"authServ/internal/services"
	"authServ/pkg/logger"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockUserService - мок для UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) RegisterUser(ctx context.Context, input services.RegisterUser) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *MockUserService) LoginUser(ctx context.Context, input services.LoginUser) (*services.Tokens, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.Tokens), args.Error(1)
}

func (m *MockUserService) CheckTruthEmail(ctx context.Context, input services.CheckTruthEmail) (services.Tokens, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(services.Tokens), args.Error(1)
}

func (m *MockUserService) ResendVerificationEmail(ctx context.Context, input services.ResendVerificationEmailInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *MockUserService) RefreshTokens(ctx context.Context, input services.RefreshInput) (services.Tokens, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(services.Tokens), args.Error(1)
}

func (m *MockUserService) Logout(ctx context.Context, input services.RefreshInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *MockUserService) ParseAccessToken(ctx context.Context, tokenString string) (uuid.UUID, error) {
	args := m.Called(ctx, tokenString)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

// MockProfileService - мок для ProfileService
type MockProfileService struct {
	mock.Mock
}

func (m *MockProfileService) SaveProfile(ctx context.Context, input services.ProfileReq, userID uuid.UUID) error {
	args := m.Called(ctx, input, userID)
	return args.Error(0)
}

// createTestLogger создает простой тестовый логгер без зависимостей от конфигурационного файла
func createTestLogger() *logger.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slogLogger := slog.New(handler)

	return &logger.Logger{
		Logger: *slogLogger,
	}
}

// createTestJWT создает тестовый JWT токен
func createTestJWT(userID uuid.UUID, secret string) string {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"email":   "test@example.com",
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

// === ТЕСТЫ ДЛЯ REGISTER ===

func TestHandler_Register(t *testing.T) {
	t.Run("Успешная регистрация", func(t *testing.T) {
		mockUserService := new(MockUserService)
		mockProfileService := new(MockProfileService)
		log := createTestLogger()
		handler := rout.NewHandler(mockUserService, mockProfileService, log)

		input := services.RegisterUser{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}

		mockUserService.On("RegisterUser", mock.Anything, input).Return(nil)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/filmbuddy/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := handler.InitRoutes()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "verification email sent")
		mockUserService.AssertExpectations(t)
	})

	t.Run("Невалидное тело запроса", func(t *testing.T) {
		mockUserService := new(MockUserService)
		mockProfileService := new(MockProfileService)
		log := createTestLogger()
		handler := rout.NewHandler(mockUserService, mockProfileService, log)

		req := httptest.NewRequest(http.MethodPost, "/filmbuddy/register", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := handler.InitRoutes()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid request body")
	})

	t.Run("Ошибка регистрации (дубликат email)", func(t *testing.T) {
		mockUserService := new(MockUserService)
		mockProfileService := new(MockProfileService)
		log := createTestLogger()
		handler := rout.NewHandler(mockUserService, mockProfileService, log)

		input := services.RegisterUser{
			Username: "testuser",
			Email:    "duplicate@example.com",
			Password: "password123",
		}

		mockUserService.On("RegisterUser", mock.Anything, input).Return(errors.New("email already exists"))

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/filmbuddy/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := handler.InitRoutes()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "email already exists")
	})
}

// === ТЕСТЫ ДЛЯ LOGIN ===

func TestHandler_Login(t *testing.T) {
	t.Run("Успешный вход", func(t *testing.T) {
		mockUserService := new(MockUserService)
		mockProfileService := new(MockProfileService)
		log := createTestLogger()
		handler := rout.NewHandler(mockUserService, mockProfileService, log)

		input := services.LoginUser{
			Email:    "test@example.com",
			Password: "password123",
		}

		tokens := &services.Tokens{
			AccessToken:  "access_token",
			RefreshToken: "refresh_token",
		}

		mockUserService.On("LoginUser", mock.Anything, input).Return(tokens, nil)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/filmbuddy/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := handler.InitRoutes()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "access_token")
		assert.Contains(t, w.Body.String(), "refresh_token")
	})

	t.Run("Неверные учетные данные", func(t *testing.T) {
		mockUserService := new(MockUserService)
		mockProfileService := new(MockProfileService)
		log := createTestLogger()
		handler := rout.NewHandler(mockUserService, mockProfileService, log)

		input := services.LoginUser{
			Email:    "wrong@example.com",
			Password: "wrongpassword",
		}

		mockUserService.On("LoginUser", mock.Anything, input).Return(nil, errors.New("invalid email or password"))

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/filmbuddy/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := handler.InitRoutes()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "invalid email or password")
	})

	t.Run("Невалидное тело запроса", func(t *testing.T) {
		mockUserService := new(MockUserService)
		mockProfileService := new(MockProfileService)
		log := createTestLogger()
		handler := rout.NewHandler(mockUserService, mockProfileService, log)

		req := httptest.NewRequest(http.MethodPost, "/filmbuddy/login", bytes.NewReader([]byte("{invalid}")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := handler.InitRoutes()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// === ТЕСТЫ ДЛЯ VERIFY EMAIL ===

func TestHandler_VerifyEmail(t *testing.T) {
	t.Run("Успешная верификация", func(t *testing.T) {
		mockUserService := new(MockUserService)
		mockProfileService := new(MockProfileService)
		log := createTestLogger()
		handler := rout.NewHandler(mockUserService, mockProfileService, log)

		input := services.CheckTruthEmail{
			Email: "test@example.com",
			Code:  "123456",
		}

		tokens := services.Tokens{
			AccessToken:  "new_access_token",
			RefreshToken: "new_refresh_token",
		}

		mockUserService.On("CheckTruthEmail", mock.Anything, input).Return(tokens, nil)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/filmbuddy/verify", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := handler.InitRoutes()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "new_access_token")
	})

	t.Run("Неверный код верификации", func(t *testing.T) {
		mockUserService := new(MockUserService)
		mockProfileService := new(MockProfileService)
		log := createTestLogger()
		handler := rout.NewHandler(mockUserService, mockProfileService, log)

		input := services.CheckTruthEmail{
			Email: "test@example.com",
			Code:  "wrong_code",
		}

		mockUserService.On("CheckTruthEmail", mock.Anything, input).Return(services.Tokens{}, errors.New("invalid verification code"))

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/filmbuddy/verify", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := handler.InitRoutes()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid verification code")
	})
}

// === ТЕСТЫ ДЛЯ REFRESH TOKENS ===

func TestHandler_RefreshTokens(t *testing.T) {
	t.Run("Успешное обновление токенов", func(t *testing.T) {
		mockUserService := new(MockUserService)
		mockProfileService := new(MockProfileService)
		log := createTestLogger()
		handler := rout.NewHandler(mockUserService, mockProfileService, log)

		input := services.RefreshInput{
			RefreshToken: "old_refresh_token",
		}

		newTokens := services.Tokens{
			AccessToken:  "new_access_token",
			RefreshToken: "new_refresh_token",
		}

		mockUserService.On("RefreshTokens", mock.Anything, input).Return(newTokens, nil)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/filmbuddy/refresh", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := handler.InitRoutes()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "new_access_token")
	})

	t.Run("Невалидный refresh token", func(t *testing.T) {
		mockUserService := new(MockUserService)
		mockProfileService := new(MockProfileService)
		log := createTestLogger()
		handler := rout.NewHandler(mockUserService, mockProfileService, log)

		input := services.RefreshInput{
			RefreshToken: "invalid_token",
		}

		mockUserService.On("RefreshTokens", mock.Anything, input).Return(services.Tokens{}, errors.New("invalid or expired refresh token"))

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/filmbuddy/refresh", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := handler.InitRoutes()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "invalid or expired refresh token")
	})
}

// === ТЕСТЫ ДЛЯ LOGOUT ===

func TestHandler_Logout(t *testing.T) {
	t.Run("Успешный logout", func(t *testing.T) {
		mockUserService := new(MockUserService)
		mockProfileService := new(MockProfileService)
		log := createTestLogger()
		handler := rout.NewHandler(mockUserService, mockProfileService, log)

		input := services.RefreshInput{
			RefreshToken: "some_refresh_token",
		}

		mockUserService.On("Logout", mock.Anything, input).Return(nil)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/filmbuddy/logout", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := handler.InitRoutes()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "logged out successfully")
	})

	t.Run("Ошибка при logout", func(t *testing.T) {
		mockUserService := new(MockUserService)
		mockProfileService := new(MockProfileService)
		log := createTestLogger()
		handler := rout.NewHandler(mockUserService, mockProfileService, log)

		input := services.RefreshInput{
			RefreshToken: "invalid_token",
		}

		mockUserService.On("Logout", mock.Anything, input).Return(errors.New("token not found"))

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/filmbuddy/logout", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := handler.InitRoutes()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "token not found")
	})
}

// === ТЕСТЫ ДЛЯ RESEND VERIFICATION EMAIL ===

func TestHandler_ResendVerificationEmail(t *testing.T) {
	t.Run("Успешная повторная отправка", func(t *testing.T) {
		mockUserService := new(MockUserService)
		mockProfileService := new(MockProfileService)
		log := createTestLogger()
		handler := rout.NewHandler(mockUserService, mockProfileService, log)

		input := services.ResendVerificationEmailInput{
			Email: "test@example.com",
		}

		mockUserService.On("ResendVerificationEmail", mock.Anything, input).Return(nil)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/filmbuddy/resend-verification", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := handler.InitRoutes()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "a new verification email has been sent")
	})

	t.Run("Слишком частые запросы (rate limit)", func(t *testing.T) {
		mockUserService := new(MockUserService)
		mockProfileService := new(MockProfileService)
		log := createTestLogger()
		handler := rout.NewHandler(mockUserService, mockProfileService, log)

		input := services.ResendVerificationEmailInput{
			Email: "test@example.com",
		}

		mockUserService.On("ResendVerificationEmail", mock.Anything, input).Return(errors.New("please wait before requesting another verification email"))

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/filmbuddy/resend-verification", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := handler.InitRoutes()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)
		assert.Contains(t, w.Body.String(), "please wait")
	})
}

// === ТЕСТЫ ДЛЯ SAVE PROFILE ===

func TestHandler_SaveProfile(t *testing.T) {
	t.Run("Успешное сохранение профиля", func(t *testing.T) {
		mockUserService := new(MockUserService)
		mockProfileService := new(MockProfileService)
		log := createTestLogger()
		handler := rout.NewHandler(mockUserService, mockProfileService, log)

		userID := uuid.New()
		token := createTestJWT(userID, "secret")

		mockUserService.On("ParseAccessToken", mock.Anything, token).Return(userID, nil)

		profileInput := rout.ProfileReq{
			FirstName: "John",
			LastName:  "Doe",
			Age:       30,
			Info:      "Developer",
			City:      "Moscow",
		}

		mockProfileService.On("SaveProfile", mock.Anything, mock.MatchedBy(func(req services.ProfileReq) bool {
			return req.FirstName == "John" && req.LastName == "Doe"
		}), userID).Return(nil)

		body, _ := json.Marshal(profileInput)
		req := httptest.NewRequest(http.MethodPost, "/filmbuddy/profile", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router := handler.InitRoutes()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "profile saved successfully")
		mockProfileService.AssertExpectations(t)
	})

	t.Run("Отсутствует Authorization header", func(t *testing.T) {
		mockUserService := new(MockUserService)
		mockProfileService := new(MockProfileService)
		log := createTestLogger()
		handler := rout.NewHandler(mockUserService, mockProfileService, log)

		profileInput := rout.ProfileReq{
			FirstName: "John",
			LastName:  "Doe",
		}

		body, _ := json.Marshal(profileInput)
		req := httptest.NewRequest(http.MethodPost, "/filmbuddy/profile", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router := handler.InitRoutes()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "authorization header is required")
	})

	t.Run("Неверный формат Authorization header", func(t *testing.T) {
		mockUserService := new(MockUserService)
		mockProfileService := new(MockProfileService)
		log := createTestLogger()
		handler := rout.NewHandler(mockUserService, mockProfileService, log)

		profileInput := rout.ProfileReq{
			FirstName: "John",
		}

		body, _ := json.Marshal(profileInput)
		req := httptest.NewRequest(http.MethodPost, "/filmbuddy/profile", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "InvalidFormat")
		w := httptest.NewRecorder()

		router := handler.InitRoutes()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "invalid authorization header format")
	})

	t.Run("Невалидный access token", func(t *testing.T) {
		mockUserService := new(MockUserService)
		mockProfileService := new(MockProfileService)
		log := createTestLogger()
		handler := rout.NewHandler(mockUserService, mockProfileService, log)

		mockUserService.On("ParseAccessToken", mock.Anything, "invalid_token").Return(uuid.Nil, errors.New("invalid token"))

		profileInput := rout.ProfileReq{
			FirstName: "John",
		}

		body, _ := json.Marshal(profileInput)
		req := httptest.NewRequest(http.MethodPost, "/filmbuddy/profile", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer invalid_token")
		w := httptest.NewRecorder()

		router := handler.InitRoutes()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "invalid token")
	})

	t.Run("Ошибка при сохранении профиля", func(t *testing.T) {
		mockUserService := new(MockUserService)
		mockProfileService := new(MockProfileService)
		log := createTestLogger()
		handler := rout.NewHandler(mockUserService, mockProfileService, log)

		userID := uuid.New()
		token := createTestJWT(userID, "secret")

		mockUserService.On("ParseAccessToken", mock.Anything, token).Return(userID, nil)
		mockProfileService.On("SaveProfile", mock.Anything, mock.Anything, userID).Return(errors.New("database error"))

		profileInput := rout.ProfileReq{
			FirstName: "John",
		}

		body, _ := json.Marshal(profileInput)
		req := httptest.NewRequest(http.MethodPost, "/filmbuddy/profile", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router := handler.InitRoutes()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "failed to save profile")
	})
}
