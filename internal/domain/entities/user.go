package entities

import (
	"crypto/rand"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User представляет собой пользователя в системе.
type User struct {
	ID                 uuid.UUID `json:"id"`
	Username           string    `json:"username"`
	Email              string    `json:"email"`
	PasswordHash       string    `json:"-"`
	CreatedAt          string    `json:"created_at"`
	VerificationCode   string    `json:"-"`
	IsVerified         bool      `json:"is_verified"`
	VerificationSentAt time.Time `json:"verification_sent_at"`
}

// ComparePassword сравнивает предоставленный пароль с хешем пароля пользователя.
func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// String возвращает строковое представление пользователя.
func (u *User) String() string {
	return u.Username + " (" + u.Email + ")"
}

// SetPassword устанавливает хеш пароля для пользователя.
func (u *User) SetPassword(password string) error {
	const op = "entities.User.SetPassword"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: %s", op, err.Error())
	}
	u.PasswordHash = string(hash)
	return nil
}

// NewUser создает нового пользователя, хешируя его пароль.
func NewUser(username, email, password string) (*User, error) {
	const op = "entities.NewUser"
	verificationCode, err := GenerateVerificationCode()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err.Error())
	}
	user := &User{
		ID:               uuid.New(),
		Username:         username,
		Email:            email,
		CreatedAt:        time.Now().Format(time.RFC3339),
		IsVerified:       false,
		VerificationCode: verificationCode,
	}

	if err := user.SetPassword(password); err != nil {
		return nil, fmt.Errorf("%s: %s", op, err.Error())
	}

	return user, nil
}

func GenerateVerificationCode() (string, error) {
	const codeLength = 5
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, codeLength)
	n, err := io.ReadAtLeast(rand.Reader, b, codeLength)
	if err != nil {
		return "", fmt.Errorf("failed to generate bytes for code: %w", err)
	}
	if n != codeLength {
		return "", fmt.Errorf("failed to generate random bytes for code: %w", err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b), nil
}
