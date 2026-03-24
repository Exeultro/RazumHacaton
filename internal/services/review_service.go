package services

import (
	"fmt"
	"time"

	"razum-backend/internal/models"
	"razum-backend/internal/repository"

	"github.com/google/uuid"
)

type ReviewService struct {
	reviewRepo *repository.ReviewRepository
	eventRepo  *repository.EventRepository
	userRepo   *repository.UserRepository
}

func NewReviewService(
	reviewRepo *repository.ReviewRepository,
	eventRepo *repository.EventRepository,
	userRepo *repository.UserRepository,
) *ReviewService {
	return &ReviewService{
		reviewRepo: reviewRepo,
		eventRepo:  eventRepo,
		userRepo:   userRepo,
	}
}

// CreateReview создает отзыв об организаторе
func (s *ReviewService) CreateReview(participantID, organizerID uuid.UUID, eventID *uuid.UUID, rating int, comment string) (*models.OrganizerReview, error) {
	// Проверяем, может ли пользователь оставить отзыв
	canReview, err := s.reviewRepo.CheckUserCanReview(participantID, organizerID)
	if err != nil {
		return nil, fmt.Errorf("failed to check review permission: %w", err)
	}
	if !canReview {
		return nil, fmt.Errorf("Вы не можете оставить отзыв об этом организаторе")
	}

	// Создаем отзыв
	review := &models.OrganizerReview{
		ID:            uuid.New(),
		OrganizerID:   organizerID,
		ParticipantID: participantID,
		EventID:       eventID,
		Rating:        rating,
		Comment:       comment,
		CreatedAt:     time.Now(),
	}

	if err := s.reviewRepo.CreateReview(review); err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	// Обновляем рейтинг организатора
	if err := s.updateOrganizerRating(organizerID); err != nil {
		// Логируем, но не блокируем
		fmt.Printf("Warning: failed to update organizer rating: %v\n", err)
	}

	return review, nil
}

// GetOrganizerReviews получает отзывы об организаторе
func (s *ReviewService) GetOrganizerReviews(organizerID uuid.UUID, page, limit int) ([]models.ReviewWithDetails, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	return s.reviewRepo.GetReviewsByOrganizer(organizerID, limit, offset)
}

// GetOrganizerRatingStats получает статистику рейтинга
func (s *ReviewService) GetOrganizerRatingStats(organizerID uuid.UUID) (*models.OrganizerRatingStats, error) {
	return s.reviewRepo.GetOrganizerRatingStats(organizerID)
}

// updateOrganizerRating обновляет рейтинг организатора
func (s *ReviewService) updateOrganizerRating(organizerID uuid.UUID) error {
	stats, err := s.reviewRepo.GetOrganizerRatingStats(organizerID)
	if err != nil {
		return err
	}

	// Обновляем рейтинг в таблице organizers
	return s.reviewRepo.UpdateOrganizerRating(organizerID, stats.AverageRating)
}
