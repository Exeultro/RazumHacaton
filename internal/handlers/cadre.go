package handlers

import (
	"strconv"

	"razum-backend/internal/repository"
	"razum-backend/internal/services"
	"razum-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CadreHandler struct {
	cadreService  *services.CadreService
	pdfService    *services.PDFService
	filterService *services.FilterService
}

func NewCadreHandler(
	cadreService *services.CadreService,
	pdfService *services.PDFService,
	filterService *services.FilterService,
) *CadreHandler {
	return &CadreHandler{
		cadreService:  cadreService,
		pdfService:    pdfService,
		filterService: filterService,
	}
}

// GetCandidates возвращает список кандидатов с фильтрацией
func (h *CadreHandler) GetCandidates(c *gin.Context) {
	var filter repository.CandidateFilter

	// Если передан filter_id, загружаем сохраненный фильтр
	if filterID := c.Query("filter_id"); filterID != "" {
		userID, exists := c.Get("user_id")
		if exists {
			fid, err := uuid.Parse(filterID)
			if err == nil {
				savedFilter, err := h.filterService.GetFilterByID(fid, userID.(uuid.UUID))
				if err == nil && savedFilter != nil {
					// Загружаем фильтр из сохраненного
					filter = repository.CandidateFilter{
						AgeMin:       savedFilter.Filters.AgeMin,
						AgeMax:       savedFilter.Filters.AgeMax,
						City:         savedFilter.Filters.City,
						Direction:    savedFilter.Filters.Direction,
						MinPoints:    savedFilter.Filters.MinPoints,
						MinEvents:    savedFilter.Filters.MinEvents,
						MinAvgPoints: savedFilter.Filters.MinAvgPoints,
						SortBy:       savedFilter.Filters.SortBy,
						SortOrder:    savedFilter.Filters.SortOrder,
					}
				}
			}
		}
	}

	// Парсим параметры из query (переопределяют сохраненный фильтр)
	if ageMin := c.Query("age_min"); ageMin != "" {
		val, err := strconv.Atoi(ageMin)
		if err == nil {
			filter.AgeMin = &val
		}
	}

	if ageMax := c.Query("age_max"); ageMax != "" {
		val, err := strconv.Atoi(ageMax)
		if err == nil {
			filter.AgeMax = &val
		}
	}

	if city := c.Query("city"); city != "" {
		filter.City = &city
	}

	if direction := c.Query("direction"); direction != "" {
		filter.Direction = &direction
	}

	if minPoints := c.Query("min_points"); minPoints != "" {
		val, err := strconv.Atoi(minPoints)
		if err == nil {
			filter.MinPoints = &val
		}
	}

	if minEvents := c.Query("min_events"); minEvents != "" {
		val, err := strconv.Atoi(minEvents)
		if err == nil {
			filter.MinEvents = &val
		}
	}

	if minAvgPoints := c.Query("min_avg_points"); minAvgPoints != "" {
		val, err := strconv.ParseFloat(minAvgPoints, 64)
		if err == nil {
			filter.MinAvgPoints = &val
		}
	}

	// Сортировка
	if sortBy := c.Query("sort_by"); sortBy != "" {
		filter.SortBy = sortBy
	} else if filter.SortBy == "" {
		filter.SortBy = "points"
	}

	if sortOrder := c.Query("sort_order"); sortOrder != "" {
		filter.SortOrder = sortOrder
	} else if filter.SortOrder == "" {
		filter.SortOrder = "desc"
	}

	// Пагинация
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	candidates, total, err := h.cadreService.GetCandidates(filter, page, limit)
	if err != nil {
		utils.InternalServerError(c, "Ошибка при получении кандидатов: "+err.Error())
		return
	}

	pages := (total + limit - 1) / limit
	if pages < 1 {
		pages = 1
	}

	utils.SuccessResponse(c, gin.H{
		"candidates": candidates,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": pages,
		},
	})
}

// ExportCandidatePDF экспортирует кандидата в PDF
func (h *CadreHandler) ExportCandidatePDF(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		utils.BadRequest(c, "Неверный ID пользователя")
		return
	}

	pdfData, err := h.pdfService.GenerateCandidatePDF(userID)
	if err != nil {
		utils.InternalServerError(c, "Ошибка при генерации PDF: "+err.Error())
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=candidate_report.pdf")
	c.Data(200, "application/pdf", pdfData)
}
