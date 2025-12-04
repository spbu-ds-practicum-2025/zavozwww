package services_test

import (
	"authServ/internal/domain/entities"
	"authServ/internal/services"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockUserProfilesRepository struct {
	mock.Mock
}

func (m *MockUserProfilesRepository) SaveProfile(ctx context.Context, profile *entities.UserProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockUserProfilesRepository) UpdateProfile(ctx context.Context, profile *entities.UserProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func TestProfileService_SaveProfile(t *testing.T) {
	userID := uuid.New()
	validReq := services.ProfileReq{
		FirstName: "Ivan",
		LastName:  "Ivanov",
		Age:       25,
		Info:      "Developer",
		City:      "Moscow",
	}

	t.Run("Успешное обновление существующего профиля", func(t *testing.T) {
		mockRepo := new(MockUserProfilesRepository)
		service := services.NewProfileService(mockRepo)
		ctx := context.Background()

		mockRepo.On("UpdateProfile", ctx, mock.MatchedBy(func(p *entities.UserProfile) bool {
			return p.UserId == userID && p.FirstName == validReq.FirstName
		})).Return(nil)

		err := service.SaveProfile(ctx, validReq, userID)

		require.NoError(t, err, "Метод не должен возвращать ошибку при успешном обновлении")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Успешное создание нового профиля (Upsert)", func(t *testing.T) {
		mockRepo := new(MockUserProfilesRepository)
		service := services.NewProfileService(mockRepo)
		ctx := context.Background()

		mockRepo.On("UpdateProfile", ctx, mock.Anything).Return(errors.New("user profile not found"))

		mockRepo.On("SaveProfile", ctx, mock.MatchedBy(func(p *entities.UserProfile) bool {
			return p.UserId == userID && p.FirstName == validReq.FirstName
		})).Return(nil)

		err := service.SaveProfile(ctx, validReq, userID)

		require.NoError(t, err, "Метод должен успешно создать профиль, если он не найден для обновления")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Ошибка базы данных при обновлении (не 'not found')", func(t *testing.T) {
		mockRepo := new(MockUserProfilesRepository)
		service := services.NewProfileService(mockRepo)
		ctx := context.Background()

		dbErr := errors.New("connection lost")

		mockRepo.On("UpdateProfile", ctx, mock.Anything).Return(dbErr)

		err := service.SaveProfile(ctx, validReq, userID)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update profile", "Ошибка должна быть обернута с контекстом")
		assert.Contains(t, err.Error(), "connection lost", "Оригинальная ошибка должна присутствовать")

		mockRepo.AssertNotCalled(t, "SaveProfile", ctx, mock.Anything)
	})

	t.Run("Ошибка базы данных при создании (после неудачного обновления)", func(t *testing.T) {
		mockRepo := new(MockUserProfilesRepository)
		service := services.NewProfileService(mockRepo)
		ctx := context.Background()

		mockRepo.On("UpdateProfile", ctx, mock.Anything).Return(errors.New("user profile not found"))

		saveErr := errors.New("duplicate key")
		mockRepo.On("SaveProfile", ctx, mock.Anything).Return(saveErr)

		err := service.SaveProfile(ctx, validReq, userID)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create profile", "Ошибка должна указывать на сбой создания")
		assert.Contains(t, err.Error(), "duplicate key")
	})
}
