package services

import (
	"authServ/internal/domain/entities"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// UserProfilesRepository определяет интерфейс для репозитория профилей пользователей.
type UserProfilesRepository interface {
	SaveProfile(ctx context.Context, profile *entities.UserProfile) error
	UpdateProfile(ctx context.Context, profile *entities.UserProfile) error
}

// ProfileReq определяет структуру данных для входящего запроса на сохранение профиля.
type ProfileReq struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
	Info      string `json:"info"`
	City      string `json:"city"`
}

// ProfileService определяет интерфейс для сервиса профилей.
type ProfileService interface {
	SaveProfile(ctx context.Context, inputProfile ProfileReq, usersId uuid.UUID) error
}

// profileService реализует интерфейс ProfileService.
type profileService struct {
	profileRepo UserProfilesRepository
}

// NewProfileService создает новый экземпляр сервиса профилей.
func NewProfileService(profileRepo UserProfilesRepository) ProfileService {
	return &profileService{
		profileRepo: profileRepo,
	}
}

// SaveProfile сохраняет или обновляет профиль пользователя.
func (s *profileService) SaveProfile(ctx context.Context, inputProfile ProfileReq, usersId uuid.UUID) error {
	const op = "services.profileService.SaveProfile"

	userProfile, err := entities.NewUserProfile(usersId, inputProfile.FirstName, inputProfile.LastName, inputProfile.Age, inputProfile.Info, inputProfile.City)
	if err != nil {
		return fmt.Errorf("%s: failed to create user profile entity: %w", op, err)
	}

	err = s.profileRepo.UpdateProfile(ctx, userProfile)
	if err != nil {
		if errors.Is(err, errors.New("user profile not found")) {
			if createErr := s.profileRepo.SaveProfile(ctx, userProfile); createErr != nil {
				return fmt.Errorf("%s: failed to create profile: %w", op, createErr)
			}
			return nil
		}
		return fmt.Errorf("%s: failed to update profile: %w", op, err)
	}

	return nil
}
