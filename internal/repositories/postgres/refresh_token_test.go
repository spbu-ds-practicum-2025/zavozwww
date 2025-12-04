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

// setupRefreshTokenTestDB подключается к БД и создает таблицу для токенов
func setupRefreshTokenTestDB(t *testing.T) (*pgxpool.Pool, func()) {
	ctx := context.Background()

	// Подключаемся к стандартной базе для создания auth_user
	defaultConnStr := "user=postgres password=6324eB6324 host=localhost port=5432 dbname=postgres sslmode=disable"

	tempPool, err := pgxpool.New(ctx, defaultConnStr)
	require.NoError(t, err, "Не удалось подключиться к системной БД postgres")

	_, _ = tempPool.Exec(ctx, "CREATE DATABASE auth_user")
	tempPool.Close()

	// Подключаемся к auth_user
	targetConnStr := os.Getenv("TEST_DATABASE_URL")
	if targetConnStr == "" {
		targetConnStr = "user=postgres password=6324eB6324 host=localhost port=5432 dbname=auth_user sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, targetConnStr)
	require.NoError(t, err, "Не удалось подключиться к БД auth_user")

	// Создаем таблицу refresh_tokens
	_, err = pool.Exec(ctx, `
        CREATE TABLE IF NOT EXISTS refresh_tokens (
            id SERIAL PRIMARY KEY,
            user_id UUID NOT NULL,
            token TEXT NOT NULL UNIQUE,
            expires_at TIMESTAMP NOT NULL,
            created_at TIMESTAMP DEFAULT NOW()
        )
    `)
	require.NoError(t, err, "Не удалось создать таблицу refresh_tokens")

	cleanup := func() {
		_, _ = pool.Exec(ctx, "TRUNCATE TABLE refresh_tokens CASCADE")
		pool.Close()
	}

	return pool, cleanup
}

// createValidToken создает валидный токен для тестов
func createValidToken(userID uuid.UUID) *entities.RefreshToken {
	return &entities.RefreshToken{
		UserID:    userID,
		Token:     "token-" + uuid.New().String(),
		ExpiresAt: time.Now().Add(24 * time.Hour).UTC(),
	}
}

func TestRefreshTokenRepository_SaveToken(t *testing.T) {
	pool, cleanup := setupRefreshTokenTestDB(t)
	defer cleanup()

	repo := repositories.NewRefreshTokenRepository(pool)
	ctx := context.Background()

	t.Run("Успешное сохранение токена", func(t *testing.T) {
		token := createValidToken(uuid.New())

		err := repo.SaveToken(ctx, token)
		require.NoError(t, err)

		savedToken, err := repo.GetToken(ctx, token.Token)
		require.NoError(t, err)
		assert.Equal(t, token.UserID, savedToken.UserID)
		assert.Equal(t, token.Token, savedToken.Token)
		assert.WithinDuration(t, token.ExpiresAt, savedToken.ExpiresAt, 1*time.Second)
	})

	t.Run("Ошибка при сохранении дубликата токена", func(t *testing.T) {
		token := createValidToken(uuid.New())

		err := repo.SaveToken(ctx, token)
		require.NoError(t, err)

		// Пытаемся сохранить тот же токен еще раз
		err = repo.SaveToken(ctx, token)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "token already exists")
	})

	t.Run("Сохранение нескольких токенов для одного пользователя", func(t *testing.T) {
		userID := uuid.New()
		token1 := createValidToken(userID)
		token2 := createValidToken(userID)

		err := repo.SaveToken(ctx, token1)
		require.NoError(t, err)

		err = repo.SaveToken(ctx, token2)
		require.NoError(t, err)

		// Оба токена должны существовать
		_, err = repo.GetToken(ctx, token1.Token)
		require.NoError(t, err)
		_, err = repo.GetToken(ctx, token2.Token)
		require.NoError(t, err)
	})
}

func TestRefreshTokenRepository_GetToken(t *testing.T) {
	pool, cleanup := setupRefreshTokenTestDB(t)
	defer cleanup()

	repo := repositories.NewRefreshTokenRepository(pool)
	ctx := context.Background()

	t.Run("Успешное получение токена", func(t *testing.T) {
		token := createValidToken(uuid.New())

		err := repo.SaveToken(ctx, token)
		require.NoError(t, err)

		retrievedToken, err := repo.GetToken(ctx, token.Token)
		require.NoError(t, err)
		assert.Equal(t, token.UserID, retrievedToken.UserID)
		assert.Equal(t, token.Token, retrievedToken.Token)
		assert.WithinDuration(t, token.ExpiresAt, retrievedToken.ExpiresAt, 1*time.Second)
	})

	t.Run("Токен не найден", func(t *testing.T) {
		nonExistentToken := "non-existent-token-" + uuid.New().String()

		retrievedToken, err := repo.GetToken(ctx, nonExistentToken)
		require.Error(t, err)
		assert.Nil(t, retrievedToken)
		assert.Contains(t, err.Error(), "token not found")
	})
}

func TestRefreshTokenRepository_DeleteToken(t *testing.T) {
	pool, cleanup := setupRefreshTokenTestDB(t)
	defer cleanup()

	repo := repositories.NewRefreshTokenRepository(pool)
	ctx := context.Background()

	t.Run("Успешное удаление токена", func(t *testing.T) {
		token := createValidToken(uuid.New())

		err := repo.SaveToken(ctx, token)
		require.NoError(t, err)

		// Удаляем токен
		err = repo.DeleteToken(ctx, token.Token)
		require.NoError(t, err)

		// Проверяем, что токен удален
		deletedToken, err := repo.GetToken(ctx, token.Token)
		require.Error(t, err)
		assert.Nil(t, deletedToken)
	})

	t.Run("Ошибка при удалении несуществующего токена", func(t *testing.T) {
		nonExistentToken := "non-existent-delete-" + uuid.New().String()

		err := repo.DeleteToken(ctx, nonExistentToken)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "token not found")
	})
}

func TestRefreshTokenRepository_ComplexScenario(t *testing.T) {
	pool, cleanup := setupRefreshTokenTestDB(t)
	defer cleanup()

	repo := repositories.NewRefreshTokenRepository(pool)
	ctx := context.Background()

	t.Run("Полный жизненный цикл токена", func(t *testing.T) {
		userID := uuid.New()
		token1 := createValidToken(userID)

		// 1. Создаем токен
		err := repo.SaveToken(ctx, token1)
		require.NoError(t, err)

		// 2. Получаем токен
		retrievedToken, err := repo.GetToken(ctx, token1.Token)
		require.NoError(t, err)
		assert.Equal(t, token1.Token, retrievedToken.Token)
		assert.Equal(t, userID, retrievedToken.UserID)

		// 3. Создаем еще один токен для того же пользователя
		token2 := createValidToken(userID)
		err = repo.SaveToken(ctx, token2)
		require.NoError(t, err)

		// 4. Удаляем первый токен
		err = repo.DeleteToken(ctx, token1.Token)
		require.NoError(t, err)

		// 5. Проверяем, что первый токен удален, а второй остался
		_, err = repo.GetToken(ctx, token1.Token)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "token not found")

		foundSecondToken, err := repo.GetToken(ctx, token2.Token)
		require.NoError(t, err)
		assert.Equal(t, token2.Token, foundSecondToken.Token)

		// 6. Удаляем второй токен
		err = repo.DeleteToken(ctx, token2.Token)
		require.NoError(t, err)

		// 7. Проверяем, что все токены удалены
		_, err = repo.GetToken(ctx, token2.Token)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "token not found")
	})
}
