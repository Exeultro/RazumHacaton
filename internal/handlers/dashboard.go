package handlers

import (
	"strconv"

	"razum-backend/internal/services"
	"razum-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DashboardHandler struct {
	dashboardService *services.DashboardService
}

func NewDashboardHandler(dashboardService *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// GetDashboard возвращает полные данные дашборда
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c)
		return
	}

	data, err := h.dashboardService.GetDashboardData(userID.(uuid.UUID))
	if err != nil {
		utils.InternalServerError(c, "Failed to get dashboard data: "+err.Error())
		return
	}

	utils.SuccessResponse(c, data)
}

// GetRecentEvents возвращает ленту мероприятий
func (h *DashboardHandler) GetRecentEvents(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	events, err := h.dashboardService.GetRecentEvents(limit)
	if err != nil {
		utils.InternalServerError(c, "Failed to get recent events: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"events": events,
		"count":  len(events),
	})
}

// GetRatingHistory возвращает историю рейтинга текущего пользователя
func (h *DashboardHandler) GetRatingHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c)
		return
	}

	history, err := h.dashboardService.GetRatingHistory(userID.(uuid.UUID))
	if err != nil {
		utils.InternalServerError(c, "Failed to get rating history: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"history": history,
	})
}

// GetTrendingTags возвращает облако тегов
func (h *DashboardHandler) GetTrendingTags(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	tags, err := h.dashboardService.GetTrendingTags(limit)
	if err != nil {
		utils.InternalServerError(c, "Failed to get trending tags: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"tags": tags,
	})
}

// GetActivityStats возвращает общую статистику
func (h *DashboardHandler) GetActivityStats(c *gin.Context) {
	stats, err := h.dashboardService.GetActivityStats()
	if err != nil {
		utils.InternalServerError(c, "Failed to get activity stats: "+err.Error())
		return
	}

	utils.SuccessResponse(c, stats)
}
