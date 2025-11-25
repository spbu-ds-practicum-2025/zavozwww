package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// RefreshToken представляет собой токен обновления для поддержания сессии пользователя.
type RefreshToken struct {
	ID        int64     `json:"-"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// NewRefreshToken создает новый токен обновления с заданными параметрами.
func NewRefreshToken(id int64, userID uuid.UUID, token string, expiresAt time.Time) *RefreshToken {
	return &RefreshToken{
		ID:        id,
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}
}

// IsExpired проверяет, истек ли срок действия токена обновления.
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// Remaining возвращает оставшееся время до истечения срока действия токена.
func (rt *RefreshToken) Remaining() time.Duration {
	return time.Until(rt.ExpiresAt)
}

// String возвращает строковое представление токена обновления.
func (rt *RefreshToken) String() string {
	return fmt.Sprintf("User: %d, until: %v, isExpired: %v", rt.UserID, rt.ExpiresAt, rt.IsExpired())
}
