package entities

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User представляет собой пользователя в системе.
type User struct {
	ID               uuid.UUID `json:"id"`
	Username         string    `json:"username"`
	Email            string    `json:"email"`
	PasswordHash     string    `json:"-"`
	CreatedAt        string    `json:"created_at"`
	VerificationCode string    `json:"-"`
	IsVerified       bool      `json:"is_verified"`
}

// String возвращает строковое представление пользователя.
func (u *User) String() string {
	return u.Username + " (" + u.Email + ")"
}

// ComparePassword сравнивает предоставленный пароль с хешем пароля пользователя.
func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
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
	verificationCode := fmt.Sprintf("%05d", rand.Intn(100000))
	user := &User{
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

//TODO
// func (u *User) Validate() error {

// }
