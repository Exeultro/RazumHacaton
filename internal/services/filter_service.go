package services

import (
	"fmt"
	"time"

	"razum-backend/internal/models"
	"razum-backend/internal/repository"

	"github.com/google/uuid"
)

type FilterService struct {
	filterRepo *repository.FilterRepository
}

func NewFilterService(filterRepo *repository.FilterRepository) *FilterService {
	return &FilterService{
		filterRepo: filterRepo,
	}
}

// CreateFilter создает новый фильтр
func (s *FilterService) CreateFilter(userID uuid.UUID, name string, criteria models.FilterCriteria) (*models.ObserverFilter, error) {
	filter := &models.ObserverFilter{
		ID:         uuid.New(),
		UserID:     userID,
		FilterName: name,
		Filters:    criteria,
		CreatedAt:  time.Now(),
	}

	if err := s.filterRepo.CreateFilter(filter); err != nil {
		return nil, fmt.Errorf("failed to create filter: %w", err)
	}

	return filter, nil
}

// GetUserFilters возвращает все фильтры пользователя
func (s *FilterService) GetUserFilters(userID uuid.UUID) ([]models.ObserverFilter, error) {
	return s.filterRepo.GetUserFilters(userID)
}

// GetFilterByID получает фильтр по ID
func (s *FilterService) GetFilterByID(filterID, userID uuid.UUID) (*models.ObserverFilter, error) {
	filter, err := s.filterRepo.GetFilterByID(filterID)
	if err != nil {
		return nil, err
	}
	if filter == nil {
		return nil, fmt.Errorf("filter not found")
	}
	if filter.UserID != userID {
		return nil, fmt.Errorf("you don't have permission to view this filter")
	}
	return filter, nil
}

// UpdateFilter обновляет фильтр
func (s *FilterService) UpdateFilter(filterID, userID uuid.UUID, name string, criteria models.FilterCriteria) (*models.ObserverFilter, error) {
	filter := &models.ObserverFilter{
		ID:         filterID,
		UserID:     userID,
		FilterName: name,
		Filters:    criteria,
	}

	if err := s.filterRepo.UpdateFilter(filter); err != nil {
		return nil, fmt.Errorf("failed to update filter: %w", err)
	}

	return filter, nil
}

// DeleteFilter удаляет фильтр
func (s *FilterService) DeleteFilter(filterID, userID uuid.UUID) error {
	return s.filterRepo.DeleteFilter(filterID, userID)
}
