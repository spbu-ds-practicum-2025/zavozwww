package repositories_test

import (
	"authServ/internal/domain/entities"
	repositories "authServ/internal/repositories/postgres"
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDB подключается к существующей БД PostgreSQL
func setupTestDB(t *testing.T) (*pgxpool.Pool, func()) {
	ctx := context.Background()

	// 1. Сначала подключаемся к стандартной базе "postgres", чтобы создать нашу базу
	defaultConnStr := "user=postgres password=6324eB6324 host=localhost port=5432 dbname=postgres sslmode=disable"

	tempPool, err := pgxpool.New(ctx, defaultConnStr)
	require.NoError(t, err, "Не удалось подключиться к системной БД postgres")

	// Создаем базу auth_user, если её нет
	// Игнорируем ошибку, если база уже существует
	_, _ = tempPool.Exec(ctx, "CREATE DATABASE auth_user")
	tempPool.Close()

	// 2. Теперь подключаемся к целевой базе "auth_user"
	targetConnStr := os.Getenv("TEST_DATABASE_URL")
	if targetConnStr == "" {
		targetConnStr = "user=postgres password=6324eB6324 host=localhost port=5432 dbname=auth_user sslmode=disable"
	}

	t.Logf("Попытка подключения к целевой БД: %s", targetConnStr)

	pool, err := pgxpool.New(ctx, targetConnStr)
	require.NoError(t, err, "Не удалось подключиться к БД auth_user")

	_, err = pool.Exec(ctx, `
        CREATE TABLE IF NOT EXISTS user_profiles (
            user_id UUID PRIMARY KEY,
            first_name TEXT NOT NULL,
            last_name TEXT NOT NULL,
            age INTEGER NOT NULL,
            info TEXT,
            city TEXT
        )
    `)
	require.NoError(t, err, "Не удалось создать таблицу user_profiles")

	cleanup := func() {
		_, _ = pool.Exec(ctx, "TRUNCATE TABLE user_profiles CASCADE")
		pool.Close()
	}

	return pool, cleanup
}

// createValidProfile создает валидный профиль для тестов
func createValidProfile(userID uuid.UUID) *entities.UserProfile {
	return &entities.UserProfile{
		UserId:    userID,
		FirstName: "John",
		LastName:  "Doe",
		Age:       30,
		Info:      "Test user profile",
		City:      "Moscow",
	}
}

func TestUserProfileRepository_SaveProfile(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewUserProfileRepository(pool)
	ctx := context.Background()

	t.Run("Успешное сохранение профиля", func(t *testing.T) {
		userID := uuid.New()
		profile := createValidProfile(userID)

		err := repo.SaveProfile(ctx, profile)
		require.NoError(t, err)

		savedProfile, err := repo.GetProfileByUserID(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, profile.UserId, savedProfile.UserId)
		assert.Equal(t, profile.FirstName, savedProfile.FirstName)
		assert.Equal(t, profile.LastName, savedProfile.LastName)
		assert.Equal(t, profile.Age, savedProfile.Age)
		assert.Equal(t, profile.Info, savedProfile.Info)
		assert.Equal(t, profile.City, savedProfile.City)
	})

	t.Run("Ошибка при сохранении дубликата", func(t *testing.T) {
		userID := uuid.New()
		profile := createValidProfile(userID)

		err := repo.SaveProfile(ctx, profile)
		require.NoError(t, err)

		err = repo.SaveProfile(ctx, profile)
		require.Error(t, err)
		assert.Contains(t, err.Error(), repositories.ErrProfileAlreadyExists.Error())
	})

	t.Run("Сохранение профиля с минимальными данными", func(t *testing.T) {
		userID := uuid.New()
		profile := createValidProfile(userID)
		profile.Info = ""
		profile.City = ""

		err := repo.SaveProfile(ctx, profile)
		require.NoError(t, err)

		savedProfile, err := repo.GetProfileByUserID(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, "", savedProfile.Info)
		assert.Equal(t, "", savedProfile.City)
	})
}

func TestUserProfileRepository_GetProfileByUserID(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewUserProfileRepository(pool)
	ctx := context.Background()

	t.Run("Успешное получение профиля", func(t *testing.T) {
		userID := uuid.New()
		profile := createValidProfile(userID)

		err := repo.SaveProfile(ctx, profile)
		require.NoError(t, err)

		retrievedProfile, err := repo.GetProfileByUserID(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, profile.UserId, retrievedProfile.UserId)
		assert.Equal(t, profile.FirstName, retrievedProfile.FirstName)
	})

	t.Run("Профиль не найден", func(t *testing.T) {
		nonExistentID := uuid.New()

		profile, err := repo.GetProfileByUserID(ctx, nonExistentID)
		require.Error(t, err)
		assert.Nil(t, profile)
		assert.Contains(t, err.Error(), repositories.ErrUserProfileNotFound.Error())
	})
}

func TestUserProfileRepository_UpdateProfile(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewUserProfileRepository(pool)
	ctx := context.Background()

	t.Run("Успешное обновление профиля", func(t *testing.T) {
		userID := uuid.New()
		profile := createValidProfile(userID)

		err := repo.SaveProfile(ctx, profile)
		require.NoError(t, err)

		profile.FirstName = "Jane"
		profile.LastName = "Smith"
		profile.Age = 25
		profile.City = "Saint Petersburg"
		profile.Info = "Updated info"

		err = repo.UpdateProfile(ctx, profile)
		require.NoError(t, err)

		updatedProfile, err := repo.GetProfileByUserID(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, "Jane", updatedProfile.FirstName)
		assert.Equal(t, "Smith", updatedProfile.LastName)
		assert.Equal(t, 25, updatedProfile.Age)
		assert.Equal(t, "Saint Petersburg", updatedProfile.City)
		assert.Equal(t, "Updated info", updatedProfile.Info)
	})

	t.Run("Ошибка при обновлении несуществующего профиля", func(t *testing.T) {
		nonExistentID := uuid.New()
		profile := createValidProfile(nonExistentID)

		err := repo.UpdateProfile(ctx, profile)
		require.Error(t, err)
		assert.Contains(t, err.Error(), repositories.ErrUserProfileNotFound.Error())
	})
}

func TestUserProfileRepository_DeleteProfile(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewUserProfileRepository(pool)
	ctx := context.Background()

	t.Run("Успешное удаление профиля", func(t *testing.T) {
		userID := uuid.New()
		profile := createValidProfile(userID)

		err := repo.SaveProfile(ctx, profile)
		require.NoError(t, err)

		err = repo.DeleteProfile(ctx, userID)
		require.NoError(t, err)

		deletedProfile, err := repo.GetProfileByUserID(ctx, userID)
		require.Error(t, err)
		assert.Nil(t, deletedProfile)
	})

	t.Run("Ошибка при удалении несуществующего профиля", func(t *testing.T) {
		nonExistentID := uuid.New()

		err := repo.DeleteProfile(ctx, nonExistentID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), repositories.ErrUserProfileNotFound.Error())
	})
}

func TestUserProfileRepository_ComplexScenario(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewUserProfileRepository(pool)
	ctx := context.Background()

	t.Run("Полный жизненный цикл профиля", func(t *testing.T) {
		userID := uuid.New()

		profile := createValidProfile(userID)
		err := repo.SaveProfile(ctx, profile)
		require.NoError(t, err)

		retrievedProfile, err := repo.GetProfileByUserID(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, "John", retrievedProfile.FirstName)

		retrievedProfile.FirstName = "UpdatedName"
		retrievedProfile.Age = 35
		err = repo.UpdateProfile(ctx, retrievedProfile)
		require.NoError(t, err)

		updatedProfile, err := repo.GetProfileByUserID(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, "UpdatedName", updatedProfile.FirstName)
		assert.Equal(t, 35, updatedProfile.Age)

		err = repo.DeleteProfile(ctx, userID)
		require.NoError(t, err)

		_, err = repo.GetProfileByUserID(ctx, userID)
		require.Error(t, err)
	})
}
