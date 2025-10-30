package entities

import (
	"fmt"
	"strings"
)

// Ограничения для полей UserProfile.
const (
	maxLenOfFirstName = 50
	minLenOfFirstName = 2
	maxLenOfLastName  = 50
	minLenOfLastName  = 2
	maxLenOfInfo      = 5000
	maxLenOfCity      = 100
	minAge            = 6
)

// Ошибки валидации для UserProfile.
const (
	ErrOfValidationFirstName = "first name should be between 2 and 50 characters"
	ErrOfValidationLastName  = "last name should be between 2 and 50 characters"
	ErrOfValidationAge       = "age should be at least 6"
	ErrOfValidationInfo      = "info should be at most 5000 characters"
	ErrOfValidationCity      = "city exceeds 100 characters"
)

// UserProfile представляет собой профиль пользователя с дополнительной информацией.
type UserProfile struct {
	UserId    int64  `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
	Info      string `json:"info"`
	City      string `json:"city"`
}

// NewUserProfile создает новый профиль пользователя с заданными параметрами.
func NewUserProfile(userId int64, firstName, lastName string, age int, info, city string) (*UserProfile, error) {
	profile := &UserProfile{
		UserId:    userId,
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

// Validate проверяет корректность полей UserProfile.
func (up *UserProfile) Validate() error {
	if len(up.FirstName) < minLenOfFirstName || len(up.FirstName) > maxLenOfFirstName {
		return fmt.Errorf(ErrOfValidationFirstName)
	}
	if len(up.LastName) < minLenOfLastName || len(up.LastName) > maxLenOfLastName {
		return fmt.Errorf(ErrOfValidationLastName)
	}
	if up.Age < minAge {
		return fmt.Errorf(ErrOfValidationAge)
	}
	if len(up.Info) > maxLenOfInfo {
		return fmt.Errorf(ErrOfValidationInfo)
	}
	if len(up.City) > maxLenOfCity {
		return fmt.Errorf(ErrOfValidationCity)
	}
	return nil
}

// FullName возвращает полное имя пользователя.
func (p *UserProfile) FullName() string {
	return strings.TrimSpace(fmt.Sprintf("%s %s", p.FirstName, p.LastName))
}

// String возвращает строковое представление профиля пользователя для логирования.
func (p *UserProfile) String() string {
	return fmt.Sprintf("Профиль пользователя %s (ID: %d)", p.FullName(), p.UserId)
}
