package services

import (
	"fmt"

	"razum-backend/internal/repository"

	"github.com/google/uuid"
)

type DashboardService struct {
	dashboardRepo *repository.DashboardRepository
}

func NewDashboardService(dashboardRepo *repository.DashboardRepository) *DashboardService {
	return &DashboardService{
		dashboardRepo: dashboardRepo,
	}
}

// GetDashboardData возвращает все данные для дашборда
func (s *DashboardService) GetDashboardData(userID uuid.UUID) (map[string]interface{}, error) {
	// Получаем ленту мероприятий
	recentEvents, err := s.dashboardRepo.GetRecentEvents(10)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent events: %w", err)
	}

	// Получаем историю рейтинга для пользователя
	ratingHistory, err := s.dashboardRepo.GetRatingHistory(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get rating history: %w", err)
	}

	// Получаем облако тегов
	trendingTags, err := s.dashboardRepo.GetTrendingTags(10)
	if err != nil {
		return nil, fmt.Errorf("failed to get trending tags: %w", err)
	}

	// Получаем общую статистику
	stats, err := s.dashboardRepo.GetActivityStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get activity stats: %w", err)
	}

	return map[string]interface{}{
		"recent_events":  recentEvents,
		"rating_history": ratingHistory,
		"trending_tags":  trendingTags,
		"stats":          stats,
	}, nil
}

// GetRecentEvents возвращает только ленту мероприятий
func (s *DashboardService) GetRecentEvents(limit int) ([]repository.RecentEvent, error) {
	if limit < 1 || limit > 50 {
		limit = 10
	}
	return s.dashboardRepo.GetRecentEvents(limit)
}

// GetRatingHistory возвращает историю рейтинга
func (s *DashboardService) GetRatingHistory(userID uuid.UUID) ([]repository.RatingHistoryPoint, error) {
	return s.dashboardRepo.GetRatingHistory(userID)
}

// GetTrendingTags возвращает облако тегов
func (s *DashboardService) GetTrendingTags(limit int) ([]repository.TagStat, error) {
	if limit < 1 || limit > 20 {
		limit = 10
	}
	return s.dashboardRepo.GetTrendingTags(limit)
}

// GetActivityStats возвращает общую статистику
func (s *DashboardService) GetActivityStats() (*repository.ActivityStats, error) {
	return s.dashboardRepo.GetActivityStats()
}
