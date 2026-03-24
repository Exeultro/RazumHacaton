package services

import (
	"fmt"

	"razum-backend/internal/repository"

	"github.com/google/uuid"
)

type CadreService struct {
	cadreRepo *repository.CadreRepository
}

func NewCadreService(cadreRepo *repository.CadreRepository) *CadreService {
	return &CadreService{
		cadreRepo: cadreRepo,
	}
}

// GetCandidates возвращает список кандидатов с фильтрацией
func (s *CadreService) GetCandidates(filter repository.CandidateFilter, page, limit int) ([]repository.CandidateInfo, int, error) {
	// Валидация параметров пагинации
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	return s.cadreRepo.GetCandidates(filter, page, limit)
}

// GetCandidateByID возвращает кандидата по ID
func (s *CadreService) GetCandidateByID(userID uuid.UUID) (*repository.CandidateInfo, error) {
	return s.cadreRepo.GetCandidateByID(userID)
}

// ExportCandidatePDF экспортирует кандидата в PDF
func (s *CadreService) ExportCandidatePDF(userID uuid.UUID) ([]byte, error) {
	candidate, err := s.GetCandidateByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get candidate: %w", err)
	}
	if candidate == nil {
		return nil, fmt.Errorf("candidate not found")
	}

	// Здесь будет генерация PDF
	// Пока возвращаем заглушку
	pdf := []byte("PDF content will be generated here")
	return pdf, nil
}
