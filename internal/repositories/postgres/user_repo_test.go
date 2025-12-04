package repositories_test

import (
	"authServ/internal/domain/entities"
	repositories "authServ/internal/repositories/postgres"
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupUserTestDB подготавливает базу данных для тестов пользователей
func setupUserTestDB(t *testing.T) (*pgxpool.Pool, func()) {
	ctx := context.Background()

	// 1. Подключение к системной БД для создания auth_user (если нужно)
	defaultConnStr := "user=postgres password=6324eB6324 host=localhost port=5432 dbname=postgres sslmode=disable"
	tempPool, err := pgxpool.New(ctx, defaultConnStr)
	require.NoError(t, err, "Не удалось подключиться к системной БД")
	_, _ = tempPool.Exec(ctx, "CREATE DATABASE auth_user")
	tempPool.Close()

	// 2. Подключение к целевой БД
	targetConnStr := os.Getenv("TEST_DATABASE_URL")
	if targetConnStr == "" {
		targetConnStr = "user=postgres password=6324eB6324 host=localhost port=5432 dbname=auth_user sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, targetConnStr)
	require.NoError(t, err, "Не удалось подключиться к БД auth_user")

	// 3. Создание таблицы users
	// ВАЖНО: Явно задаем имена CONSTRAINT, чтобы логика SaveUser (switch pgErr.ConstraintName) работала корректно
	_, err = pool.Exec(ctx, `
        CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            username TEXT NOT NULL,
            email TEXT NOT NULL,
            password_hash TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT NOW(),
            is_verified BOOLEAN DEFAULT FALSE,
            verification_code TEXT,
            verification_sent_at TIMESTAMP,
            CONSTRAINT users_username_key UNIQUE (username),
            CONSTRAINT users_email_key UNIQUE (email)
        )
    `)
	require.NoError(t, err, "Не удалось создать таблицу users")

	cleanup := func() {
		_, _ = pool.Exec(ctx, "TRUNCATE TABLE users CASCADE")
		pool.Close()
	}

	return pool, cleanup
}

// createTestUserEntity создает структуру пользователя для тестов
func createTestUserEntity() *entities.User {
	return &entities.User{
		Username:     "testuser_" + uuid.New().String()[:8],
		Email:        "test_" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "hashed_secret_password",
		IsVerified:   false,
	}
}

func TestUserRepository_SaveUser(t *testing.T) {
	pool, cleanup := setupUserTestDB(t)
	defer cleanup()

	repo := repositories.NewUserRepository(pool)
	ctx := context.Background()

	t.Run("Успешное создание пользователя", func(t *testing.T) {
		user := createTestUserEntity()

		id, err := repo.SaveUser(ctx, user)
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, id)

		// Проверяем, что пользователь реально в базе
		savedUser, err := repo.GetUserByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, user.Email, savedUser.Email)
		assert.Equal(t, user.Username, savedUser.Username)
	})

	t.Run("Ошибка: дубликат Email", func(t *testing.T) {
		user1 := createTestUserEntity()
		user1.Email = "duplicate@example.com"

		_, err := repo.SaveUser(ctx, user1)
		require.NoError(t, err)

		user2 := createTestUserEntity()
		user2.Email = "duplicate@example.com" // Тот же email

		_, err = repo.SaveUser(ctx, user2)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "email already exists")
	})

	t.Run("Ошибка: дубликат Username", func(t *testing.T) {
		user1 := createTestUserEntity()
		user1.Username = "unique_user"

		_, err := repo.SaveUser(ctx, user1)
		require.NoError(t, err)

		user2 := createTestUserEntity()
		user2.Username = "unique_user" // Тот же username

		_, err = repo.SaveUser(ctx, user2)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "username already exists")
	})
}

func TestUserRepository_GetUserByID(t *testing.T) {
	pool, cleanup := setupUserTestDB(t)
	defer cleanup()

	repo := repositories.NewUserRepository(pool)
	ctx := context.Background()

	t.Run("Пользователь найден", func(t *testing.T) {
		user := createTestUserEntity()
		id, err := repo.SaveUser(ctx, user)
		require.NoError(t, err)

		foundUser, err := repo.GetUserByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id, foundUser.ID)
		assert.Equal(t, user.Username, foundUser.Username)
		assert.NotEmpty(t, foundUser.CreatedAt) // Проверяем, что дата создания заполнилась
	})

	t.Run("Пользователь не найден", func(t *testing.T) {
		_, err := repo.GetUserByID(ctx, uuid.New())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})
}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	pool, cleanup := setupUserTestDB(t)
	defer cleanup()

	repo := repositories.NewUserRepository(pool)
	ctx := context.Background()

	t.Run("Пользователь найден по email", func(t *testing.T) {
		user := createTestUserEntity()
		id, err := repo.SaveUser(ctx, user)
		require.NoError(t, err)

		foundUser, err := repo.GetUserByEmail(ctx, user.Email)
		require.NoError(t, err)
		assert.Equal(t, id, foundUser.ID)
		assert.Equal(t, user.Email, foundUser.Email)
	})

	t.Run("Пользователь не найден по email", func(t *testing.T) {
		_, err := repo.GetUserByEmail(ctx, "nonexistent@example.com")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})
}

func TestUserRepository_UpdateUserName(t *testing.T) {
	pool, cleanup := setupUserTestDB(t)
	defer cleanup()

	repo := repositories.NewUserRepository(pool)
	ctx := context.Background()

	t.Run("Успешное обновление имени", func(t *testing.T) {
		user := createTestUserEntity()
		id, err := repo.SaveUser(ctx, user)
		require.NoError(t, err)

		newUsername := "updated_name_" + uuid.New().String()[:5]
		err = repo.UpdateUserName(ctx, id, newUsername)
		require.NoError(t, err)

		updatedUser, err := repo.GetUserByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, newUsername, updatedUser.Username)
	})

	t.Run("Ошибка: имя уже занято", func(t *testing.T) {
		// Создаем первого пользователя
		user1 := createTestUserEntity()
		user1.Username = "occupied_name"
		_, err := repo.SaveUser(ctx, user1)
		require.NoError(t, err)

		// Создаем второго пользователя
		user2 := createTestUserEntity()
		id2, err := repo.SaveUser(ctx, user2)
		require.NoError(t, err)

		// Пытаемся второму присвоить имя первого
		err = repo.UpdateUserName(ctx, id2, "occupied_name")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "username already exists")
	})

	t.Run("Ошибка: пользователь не найден", func(t *testing.T) {
		err := repo.UpdateUserName(ctx, uuid.New(), "some_name")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})
}

func TestUserRepository_UpdateUser(t *testing.T) {
	pool, cleanup := setupUserTestDB(t)
	defer cleanup()

	repo := repositories.NewUserRepository(pool)
	ctx := context.Background()

	t.Run("Полное обновление пользователя", func(t *testing.T) {
		user := createTestUserEntity()
		id, err := repo.SaveUser(ctx, user)
		require.NoError(t, err)

		// Подготавливаем обновленные данные
		// Важно: ID должен совпадать
		user.ID = id
		user.Username = "new_full_update_name"
		user.Email = "new_email@example.com"
		user.IsVerified = true
		user.VerificationCode = "123456"

		// Используем UTC для корректного сравнения
		now := time.Now().UTC()
		user.VerificationSentAt = now

		err = repo.UpdateUser(ctx, user)
		require.NoError(t, err)

		// Проверяем изменения
		updatedUser, err := repo.GetUserByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "new_full_update_name", updatedUser.Username)
		assert.Equal(t, "new_email@example.com", updatedUser.Email)
		assert.True(t, updatedUser.IsVerified)
		assert.Equal(t, "123456", updatedUser.VerificationCode)

		// Проверка времени с допуском
		require.NotNil(t, updatedUser.VerificationSentAt)
		assert.WithinDuration(t, now, updatedUser.VerificationSentAt, time.Second)
	})

	t.Run("Ошибка при обновлении несуществующего пользователя", func(t *testing.T) {
		user := createTestUserEntity()
		user.ID = uuid.New() // Несуществующий ID

		err := repo.UpdateUser(ctx, user)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})
}

func TestUserRepository_UpdateVerificationCode(t *testing.T) {
	pool, cleanup := setupUserTestDB(t)
	defer cleanup()

	repo := repositories.NewUserRepository(pool)
	ctx := context.Background()

	t.Run("Успешное обновление кода", func(t *testing.T) {
		user := createTestUserEntity()
		_, err := repo.SaveUser(ctx, user)
		require.NoError(t, err)

		newCode := "999888"
		sentAt := time.Now().UTC()

		err = repo.UpdateVerificationCode(ctx, user.Email, newCode, sentAt)
		require.NoError(t, err)

		updatedUser, err := repo.GetUserByEmail(ctx, user.Email)
		require.NoError(t, err)
		assert.Equal(t, newCode, updatedUser.VerificationCode)

		require.NotNil(t, updatedUser.VerificationSentAt)
		assert.WithinDuration(t, sentAt, updatedUser.VerificationSentAt, time.Second)
	})

	t.Run("Ошибка: email не найден", func(t *testing.T) {
		err := repo.UpdateVerificationCode(ctx, "wrong@email.com", "111", time.Now())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})
}

func TestUserRepository_DeleteUserByID(t *testing.T) {
	pool, cleanup := setupUserTestDB(t)
	defer cleanup()

	repo := repositories.NewUserRepository(pool)
	ctx := context.Background()

	t.Run("Успешное удаление", func(t *testing.T) {
		user := createTestUserEntity()
		id, err := repo.SaveUser(ctx, user)
		require.NoError(t, err)

		err = repo.DeleteUserByID(ctx, id)
		require.NoError(t, err)

		// Проверяем, что пользователя больше нет
		_, err = repo.GetUserByID(ctx, id)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("Ошибка удаления несуществующего пользователя", func(t *testing.T) {
		err := repo.DeleteUserByID(ctx, uuid.New())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})
}
