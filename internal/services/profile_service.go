package services

import (
	"errors"
	"fmt"

	"razum-backend/internal/models"
	"razum-backend/internal/repository"

	"github.com/google/uuid"
)

type ProfileService struct {
	userRepo *repository.UserRepository
}

func NewProfileService(userRepo *repository.UserRepository) *ProfileService {
	return &ProfileService{
		userRepo: userRepo,
	}
}

func (s *ProfileService) GetMyProfile(userID uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}

// GetProfile получает профиль пользователя с дополнительной информацией
func (s *ProfileService) GetProfile(userID uuid.UUID) (*models.User, interface{}, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, nil, fmt.Errorf("user not found")
	}

	var extra interface{}

	// Если организатор, добавляем статистику
	if user.Role == models.RoleOrganizer {
		stats, err := s.userRepo.GetOrganizerStats(userID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get organizer stats: %w", err)
		}
		extra = stats
	}

	return user, extra, nil
}

// UpdateProfile обновляет профиль пользователя
func (s *ProfileService) UpdateProfile(userID uuid.UUID, fullName *string, city *string, age *int, direction *models.Direction, avatarURL *string) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	if fullName != nil {
		user.FullName = *fullName
	}
	if city != nil {
		user.City = city
	}
	if age != nil {
		user.Age = age
	}
	if direction != nil {
		user.Direction = direction
	}
	if avatarURL != nil {
		user.AvatarURL = avatarURL
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// GetPublicProfile получает публичный профиль другого пользователя
func (s *ProfileService) GetPublicProfile(userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *ProfileService) GetOrganizerStats(userID uuid.UUID) (*models.Organizer, error) {
	return s.userRepo.GetOrganizerStats(userID)
}
