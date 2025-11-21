package entities

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const (
	ErrOfValidationFirstName = "first name must be between 2 and 50 characters"
	ErrOfValidationLastName  = "last name must be between 2 and 50 characters"
	ErrOfValidationAge       = "age must be greater than 6"
	ErrOfValidationInfo      = "info must be less than 5000 characters"
	ErrOfValidationCity      = "city must be less than 100 characters"
	ErrOfValidationUserID    = "user ID cannot be nil"
)

// UserProfile представляет профиль пользователя.
type UserProfile struct {
	UserId    uuid.UUID `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Age       int       `json:"age"`
	Info      string    `json:"info"`
	City     string    `json:"city"`
}

// NewUserProfile создает и валидирует новый профиль пользователя.
func NewUserProfile(userID uuid.UUID, firstName, lastName string, age int, info, city string) (*UserProfile, error) {
	profile := &UserProfile{
		UserId:    userID,
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
		Info:      info,
		City:      city,
	}

	if err := profile.Validate(); err != nil {
		return nil, err
	}

	return profile, nil
}

// Validate проверяет корректность данных профиля.
func (p *UserProfile) Validate() error {
	if p.UserId == uuid.Nil {
		return fmt.Errorf(ErrOfValidationUserID)
	}
	if len(p.FirstName) < 2 || len(p.FirstName) > 50 {
		return fmt.Errorf(ErrOfValidationFirstName)
	}
	if len(p.LastName) < 2 || len(p.LastName) > 50 {
		return fmt.Errorf(ErrOfValidationLastName)
	}
	if p.Age < 6 {
		return fmt.Errorf(ErrOfValidationAge)
	}
	if len(p.Info) > 5000 {
		return fmt.Errorf(ErrOfValidationInfo)
	}
	if len(p.City) > 100 {
		return fmt.Errorf(ErrOfValidationCity)
	}
	return nil
}

// FullName возвращает полное имя пользователя.
func (p *UserProfile) FullName() string {
	return fmt.Sprintf("%s %s", strings.TrimSpace(p.FirstName), strings.TrimSpace(p.LastName))
}

// String возвращает строковое представление профиля.
func (p *UserProfile) String() string {
	return fmt.Sprintf("Профиль пользователя %s (ID: %s)", p.FullName(), p.UserId.String())
}
