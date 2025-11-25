package entities_test

import (
	"authServ/internal/domain/entities"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserToString(t *testing.T) {
	user := &entities.User{
		Username: "testuser",
		Email:    "testuser@example.com",
	}

	expected := "testuser (testuser@example.com)"
	require.Equal(t, expected, user.String())
}

func TestUser_Password(t *testing.T) {
	user := &entities.User{}
	password := "my-secret-password-123"

	t.Run("Успешная установка пароля", func(t *testing.T) {
		err := user.SetPassword(password)
		require.NoError(t, err, "Установка пароля не должна вызывать ошибку")
		require.NotEmpty(t, user.PasswordHash, "Хеш пароля не должен быть пустым после установки")
	})

	t.Run("Ошибка при установке слишком длинного пароля", func(t *testing.T) {
		longPassword := strings.Repeat("a", 73)
		err := user.SetPassword(longPassword)
		require.Error(t, err, "Должна быть ошибка при установке слишком длинного пароля")
	})

	err := user.SetPassword(password)
	require.NoError(t, err)

	t.Run("Корректный пароль", func(t *testing.T) {
		match := user.ComparePassword(password)
		assert.True(t, match, "Сравнение с правильным паролем должно возвращать true")
	})

	t.Run("Некорректный пароль", func(t *testing.T) {
		match := user.ComparePassword("wrong-password")
		assert.False(t, match, "Сравнение с неправильным паролем должно возвращать false")
	})

	t.Run("Пустой пароль", func(t *testing.T) {
		match := user.ComparePassword("")
		assert.False(t, match, "Сравнение с пустым паролем должно возвращать false")
	})

	t.Run("Сравнение при некорректном хеше", func(t *testing.T) {
		userWithInvalidHash := &entities.User{PasswordHash: "not-a-valid-bcrypt-hash"}
		match := userWithInvalidHash.ComparePassword(password)
		assert.False(t, match, "Сравнение с невалидным хешем должно всегда возвращать false")
	})

	t.Run("Сравнение при пустом хеше", func(t *testing.T) {
		userWithEmptyHash := &entities.User{}
		match := userWithEmptyHash.ComparePassword(password)
		assert.False(t, match, "Сравнение при пустом хеше должно всегда возвращать false")
	})
}

func TestNewUser(t *testing.T) {
	t.Run("Успешное создание", func(t *testing.T) {
		username := "new_user"
		email := "new@example.com"
		password := "super-strong-password"

		user, err := entities.NewUser(username, email, password)

		require.NoError(t, err)
		require.NotNil(t, user, "Созданный пользователь не должен быть nil")

		assert.Equal(t, username, user.Username)
		assert.Equal(t, email, user.Email)
		assert.NotEmpty(t, user.PasswordHash, "Хеш пароля должен быть установлен")
		assert.NotEmpty(t, user.CreatedAt, "Дата создания должна быть установлена")

		// --- НОВЫЕ ПРОВЕРКИ ---
		assert.False(t, user.IsVerified, "Новый пользователь не должен быть верифицирован")
		assert.NotEmpty(t, user.VerificationCode, "Код верификации должен быть сгенерирован")
		assert.Len(t, user.VerificationCode, 5, "Код верификации должен состоять из 5 символов")
		_, errConv := strconv.Atoi(user.VerificationCode)
		assert.NoError(t, errConv, "Код верификации должен состоять только из цифр")
		// --- КОНЕЦ НОВЫХ ПРОВЕРОК ---

		assert.True(t, user.ComparePassword(password), "Пароль, установленный через конструктор, должен быть корректным")
	})

	t.Run("Ошибка при создании со слишком длинным паролем", func(t *testing.T) {
		username := "long_pass_user"
		email := "long@pass.com"
		longPassword := strings.Repeat("b", 73)

		_, err := entities.NewUser(username, email, longPassword)
		require.Error(t, err, "Конструктор должен возвращать ошибку при слишком длинном пароле")
	})
}
