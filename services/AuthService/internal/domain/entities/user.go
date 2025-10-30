// В этом пакете отображаются основные сущности нашего сервиса регитрации и аутентификации пользователей.
package entities

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User представляет собой пользователя в системе.
type User struct {
	ID             int64
	Username       string `json:"username"`
	Email          string `json:"email"`
	PasswordHash   string `json:"password_hash"`
	DateOfRegister string `json:"date_of_register"`
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
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

// NewUser создает нового пользователя, хешируя его пароль.
func NewUser(username, email, password string) (*User, error) {
	user := &User{
		Username:       username,
		Email:          email,
		DateOfRegister: time.Now().Format(time.RFC3339),
	}

	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	return user, nil
}

//TODO
// func (u *User) Validate() error {

// }
