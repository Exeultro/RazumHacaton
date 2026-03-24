package handlers

import (
	"razum-backend/internal/models"
	"razum-backend/internal/services"
	"razum-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FilterHandler struct {
	filterService *services.FilterService
}

func NewFilterHandler(filterService *services.FilterService) *FilterHandler {
	return &FilterHandler{
		filterService: filterService,
	}
}

// CreateFilter создает фильтр
func (h *FilterHandler) CreateFilter(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c)
		return
	}

	var req struct {
		FilterName string                `json:"filter_name" binding:"required"`
		Filters    models.FilterCriteria `json:"filters" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Неверный формат запроса: "+err.Error())
		return
	}

	filter, err := h.filterService.CreateFilter(userID.(uuid.UUID), req.FilterName, req.Filters)
	if err != nil {
		utils.InternalServerError(c, "Ошибка при создании фильтра: "+err.Error())
		return
	}

	utils.CreatedResponse(c, filter)
}

// GetMyFilters возвращает все фильтры текущего пользователя
func (h *FilterHandler) GetMyFilters(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c)
		return
	}

	filters, err := h.filterService.GetUserFilters(userID.(uuid.UUID))
	if err != nil {
		utils.InternalServerError(c, "Ошибка при получении фильтров: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"filters": filters,
		"count":   len(filters),
	})
}

// GetFilter получает фильтр по ID
func (h *FilterHandler) GetFilter(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c)
		return
	}

	filterIDParam := c.Param("id")
	filterID, err := uuid.Parse(filterIDParam)
	if err != nil {
		utils.BadRequest(c, "Неверный ID фильтра")
		return
	}

	filter, err := h.filterService.GetFilterByID(filterID, userID.(uuid.UUID))
	if err != nil {
		utils.NotFound(c)
		return
	}

	utils.SuccessResponse(c, filter)
}

// UpdateFilter обновляет фильтр
func (h *FilterHandler) UpdateFilter(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c)
		return
	}

	filterIDParam := c.Param("id")
	filterID, err := uuid.Parse(filterIDParam)
	if err != nil {
		utils.BadRequest(c, "Неверный ID фильтра")
		return
	}

	var req struct {
		FilterName string                `json:"filter_name" binding:"required"`
		Filters    models.FilterCriteria `json:"filters" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Неверный формат запроса: "+err.Error())
		return
	}

	filter, err := h.filterService.UpdateFilter(filterID, userID.(uuid.UUID), req.FilterName, req.Filters)
	if err != nil {
		utils.InternalServerError(c, "Ошибка при обновлении фильтра: "+err.Error())
		return
	}

	utils.SuccessResponse(c, filter)
}

// DeleteFilter удаляет фильтр
func (h *FilterHandler) DeleteFilter(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c)
		return
	}

	filterIDParam := c.Param("id")
	filterID, err := uuid.Parse(filterIDParam)
	if err != nil {
		utils.BadRequest(c, "Неверный ID фильтра")
		return
	}

	if err := h.filterService.DeleteFilter(filterID, userID.(uuid.UUID)); err != nil {
		utils.InternalServerError(c, "Ошибка при удалении фильтра: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message": "Фильтр успешно удален",
	})
}
