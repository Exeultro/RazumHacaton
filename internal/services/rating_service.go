package services

import (
	"fmt"

	"razum-backend/internal/models"
	"razum-backend/internal/repository"

	"github.com/google/uuid"
)

type RatingService struct {
	ratingRepo *repository.RatingRepository
}

func NewRatingService(ratingRepo *repository.RatingRepository) *RatingService {
	return &RatingService{
		ratingRepo: ratingRepo,
	}
}

// GetGlobalRating возвращает глобальную таблицу лидеров
func (s *RatingService) GetGlobalRating(page, limit int) ([]models.RatingCache, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	return s.ratingRepo.GetGlobalRating(limit, offset)
}

// GetMyRating возвращает рейтинг текущего пользователя
func (s *RatingService) GetMyRating(userID uuid.UUID) (*models.RatingCache, error) {
	return s.GetUserRating(userID)
}

// GetRatingByDirection возвращает рейтинг по направлению
func (s *RatingService) GetRatingByDirection(direction string, page, limit int) ([]models.RatingCache, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	return s.ratingRepo.GetRatingByDirection(direction, limit, offset)
}

// GetUserRating получает рейтинг конкретного пользователя
func (s *RatingService) GetUserRating(userID uuid.UUID) (*models.RatingCache, error) {
	rating, err := s.ratingRepo.GetUserRating(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user rating: %w", err)
	}
	if rating == nil {
		// Пользователь есть, но баллов пока нет
		return &models.RatingCache{
			UserID:      userID,
			TotalPoints: 0,
			EventsCount: 0, // 👈 добавили
			GlobalRank:  nil,
		}, nil
	}
	return rating, nil
}

// UpdateRatingCache обновляет кэш рейтинга
func (s *RatingService) UpdateRatingCache() error {
	return s.ratingRepo.UpdateRatingCache()
}

// GetRatingWithUserInfo возвращает рейтинг с информацией о пользователе
func (s *RatingService) GetRatingWithUserInfo(limit, offset int) ([]RatingWithUser, int, error) {
	ratings, total, err := s.ratingRepo.GetGlobalRating(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	result := make([]RatingWithUser, 0, len(ratings))
	for _, r := range ratings {
		// Здесь можно добавить дополнительную информацию о пользователе
		result = append(result, RatingWithUser{
			RatingCache: r,
		})
	}

	return result, total, nil
}

// RatingWithUser расширенная структура рейтинга с данными пользователя
type RatingWithUser struct {
	models.RatingCache
	FullName  string  `json:"full_name"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	City      *string `json:"city,omitempty"`
	Direction *string `json:"direction,omitempty"`
}
