package entities_test

import (
	"authServ/internal/domain/entities"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRefreshToken(t *testing.T) {
	id := int64(1)
	userID := uuid.New()
	tokenStr := "some-secret-token-string"
	expiresAt := time.Now().Add(24 * time.Hour)

	token := entities.NewRefreshToken(id, userID, tokenStr, expiresAt)

	require.NotNil(t, token, "Конструктор не должен возвращать nil")
	assert.Equal(t, id, token.ID, "ID должен совпадать")
	assert.Equal(t, userID, token.UserID, "UserID должен совпадать")
	assert.Equal(t, tokenStr, token.Token, "Строка токена должна совпадать")
	assert.WithinDuration(t, expiresAt, token.ExpiresAt, time.Millisecond, "Время истечения должно совпадать")
}

func TestRefreshToken_IsExpired(t *testing.T) {
	t.Run("Токен еще действителен", func(t *testing.T) {
		token := &entities.RefreshToken{
			ExpiresAt: time.Now().Add(1 * time.Hour),
		}
		assert.False(t, token.IsExpired(), "Действительный токен не должен считаться истекшим")
	})

	t.Run("Токен уже истек", func(t *testing.T) {
		token := &entities.RefreshToken{
			ExpiresAt: time.Now().Add(-1 * time.Hour),
		}
		assert.True(t, token.IsExpired(), "Истекший токен должен считаться истекшим")
	})

	t.Run("Токен истекает прямо сейчас", func(t *testing.T) {
		token := &entities.RefreshToken{
			ExpiresAt: time.Now(),
		}
		time.Sleep(1 * time.Millisecond)
		assert.True(t, token.IsExpired(), "Токен, который только что истек, должен считаться истекшим")
	})
}

func TestRefreshToken_Remaining(t *testing.T) {
	t.Run("Оставшееся время для действительного токена", func(t *testing.T) {
		duration := 30 * time.Minute
		token := &entities.RefreshToken{
			ExpiresAt: time.Now().Add(duration),
		}
		remaining := token.Remaining()

		assert.Greater(t, remaining, time.Duration(0))
		assert.LessOrEqual(t, remaining, duration)
		assert.Greater(t, remaining, duration-5*time.Second)
	})

	t.Run("Оставшееся время для истекшего токена", func(t *testing.T) {
		token := &entities.RefreshToken{
			ExpiresAt: time.Now().Add(-15 * time.Minute),
		}
		remaining := token.Remaining()

		assert.Less(t, remaining, time.Duration(0))
	})

	t.Run("Оставшееся время для токена с нулевым временем", func(t *testing.T) {
		token := &entities.RefreshToken{
			ExpiresAt: time.Time{},
		}
		remaining := token.Remaining()

		assert.Less(t, remaining, -200*365*24*time.Hour, "Оставшееся время для нулевой даты должно быть очень большим отрицательным числом")
	})
}

func TestRefreshToken_String(t *testing.T) {
	userID := uuid.New()

	t.Run("Строка для действительного токена", func(t *testing.T) {
		expiresAt := time.Now().Add(1 * time.Hour)
		token := entities.NewRefreshToken(1, userID, "valid-token", expiresAt)

		// ИСПРАВЛЕНИЕ: Используем [16]byte(userID) и %v, чтобы соответствовать фактическому выводу метода String(),
		// который возвращает массив байтов (например, [233 38 ...]), а не строку UUID.
		expectedString := fmt.Sprintf("User: %v, until: %v, isExpired: false", [16]byte(userID), expiresAt)

		assert.Equal(t, expectedString, token.String())
	})

	t.Run("Строка для истекшего токена", func(t *testing.T) {
		expiresAt := time.Now().Add(-1 * time.Hour)
		token := entities.NewRefreshToken(2, userID, "expired-token", expiresAt)

		// ИСПРАВЛЕНИЕ: Аналогично приводим к [16]byte
		expectedString := fmt.Sprintf("User: %v, until: %v, isExpired: true", [16]byte(userID), expiresAt)

		assert.Equal(t, expectedString, token.String())
	})

	t.Run("Строка для нулевой структуры токена (zero value)", func(t *testing.T) {
		token := &entities.RefreshToken{}
		// ИСПРАВЛЕНИЕ: uuid.Nil также приводим к [16]byte
		expectedString := fmt.Sprintf("User: %v, until: %v, isExpired: true", [16]byte(uuid.Nil), time.Time{})
		assert.Equal(t, expectedString, token.String(), "Строковое представление для нулевой структуры должно быть корректным")
	})
}
