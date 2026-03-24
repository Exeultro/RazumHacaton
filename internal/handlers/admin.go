package handlers

import (
	"razum-backend/internal/repository"
	"razum-backend/internal/services"
	"razum-backend/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminHandler struct {
	adminService *services.AdminService
}

func NewAdminHandler(adminService *services.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

// GetPendingOrganizers возвращает список организаторов на модерации
func (h *AdminHandler) GetPendingOrganizers(c *gin.Context) {
	organizers, err := h.adminService.GetPendingOrganizers()
	if err != nil {
		utils.InternalServerError(c, "Failed to get pending organizers: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"pending_organizers": organizers,
		"count":              len(organizers),
	})
}

// ApproveOrganizer одобряет организатора
func (h *AdminHandler) ApproveOrganizer(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID")
		return
	}

	if err := h.adminService.ApproveOrganizer(userID); err != nil {
		utils.InternalServerError(c, "Failed to approve organizer: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Organizer approved successfully"})
}

// RejectOrganizer отклоняет организатора
func (h *AdminHandler) RejectOrganizer(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID")
		return
	}

	if err := h.adminService.RejectOrganizer(userID); err != nil {
		utils.InternalServerError(c, "Failed to reject organizer: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Organizer rejected successfully"})
}

// GetDifficultySettings получает настройки весов
func (h *AdminHandler) GetDifficultySettings(c *gin.Context) {
	settings, err := h.adminService.GetDifficultySettings()
	if err != nil {
		utils.InternalServerError(c, "Failed to get difficulty settings: "+err.Error())
		return
	}

	utils.SuccessResponse(c, settings)
}

// UpdateDifficultySettings обновляет настройки весов
func (h *AdminHandler) UpdateDifficultySettings(c *gin.Context) {
	var req repository.DifficultySettings
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := h.adminService.UpdateDifficultySettings(&req); err != nil {
		utils.InternalServerError(c, "Failed to update difficulty settings: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Settings updated successfully"})
}

// GetStats возвращает статистику по платформе
func (h *AdminHandler) GetStats(c *gin.Context) {
	stats, err := h.adminService.GetStats()
	if err != nil {
		utils.InternalServerError(c, "Failed to get stats: "+err.Error())
		return
	}
	utils.SuccessResponse(c, stats)
}

// GetAllUsers возвращает список всех пользователей
func (h *AdminHandler) GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	users, total, err := h.adminService.GetAllUsers(page, limit)
	if err != nil {
		utils.InternalServerError(c, "Ошибка при получении пользователей: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"users": users,
		"count": len(users),
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + limit - 1) / limit,
		},
	})
}

// DeleteUser удаляет пользователя
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID")
		return
	}

	if err := h.adminService.DeleteUser(userID); err != nil {
		utils.InternalServerError(c, "Failed to delete user: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "User deleted successfully"})
}

// ChangeUserRole меняет роль пользователя
func (h *AdminHandler) ChangeUserRole(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID")
		return
	}

	var req struct {
		Role string `json:"role" binding:"required,oneof=participant organizer observer admin"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := h.adminService.ChangeUserRole(userID, req.Role); err != nil {
		utils.InternalServerError(c, "Failed to change role: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "User role changed successfully"})
}

// GetAllEvents возвращает все мероприятия (для админа)
func (h *AdminHandler) GetAllEvents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	events, total, err := h.adminService.GetAllEvents(page, limit)
	if err != nil {
		utils.InternalServerError(c, "Ошибка при получении мероприятий: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"events": events,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + limit - 1) / limit,
		},
	})
}

// DeleteEventByAdmin удаляет любое мероприятие (для админа)
func (h *AdminHandler) DeleteEventByAdmin(c *gin.Context) {
	eventIDParam := c.Param("id")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		utils.BadRequest(c, "Неверный ID мероприятия")
		return
	}

	if err := h.adminService.DeleteEventByAdmin(eventID); err != nil {
		utils.InternalServerError(c, "Ошибка при удалении мероприятия: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message": "Мероприятие успешно удалено",
	})
}
