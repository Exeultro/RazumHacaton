package handlers

import (
	"razum-backend/internal/models"
	"razum-backend/internal/services"
	"razum-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProfileHandler struct {
	profileService *services.ProfileService
}

func NewProfileHandler(profileService *services.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

// GetMyProfile получает профиль текущего пользователя
func (h *ProfileHandler) GetMyProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c)
		return
	}

	user, extra, err := h.profileService.GetProfile(userID.(uuid.UUID))
	if err != nil {
		utils.InternalServerError(c, "Failed to get profile: "+err.Error())
		return
	}

	response := gin.H{
		"user": user,
	}
	if extra != nil {
		response["organizer_stats"] = extra
	}

	utils.SuccessResponse(c, response)
}

// UpdateMyProfile обновляет профиль текущего пользователя
func (h *ProfileHandler) UpdateMyProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c)
		return
	}

	var req struct {
		FullName  *string           `json:"full_name"`
		City      *string           `json:"city"`
		Age       *int              `json:"age"`
		Direction *models.Direction `json:"direction"`
		AvatarURL *string           `json:"avatar_url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	user, err := h.profileService.UpdateProfile(
		userID.(uuid.UUID),
		req.FullName,
		req.City,
		req.Age,
		req.Direction,
		req.AvatarURL,
	)
	if err != nil {
		utils.InternalServerError(c, "Failed to update profile: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"user": user})
}

// GetPublicProfile получает публичный профиль пользователя
func (h *ProfileHandler) GetPublicProfile(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		utils.BadRequest(c, "Неверный ID пользователя")
		return
	}

	user, err := h.profileService.GetPublicProfile(userID)
	if err != nil {
		if err.Error() == "user not found" {
			utils.NotFound(c)
			return
		}
		utils.InternalServerError(c, "Ошибка при получении профиля: "+err.Error())
		return
	}

	// Если это организатор, добавляем статистику
	if user.Role == "organizer" {
		stats, _ := h.profileService.GetOrganizerStats(userID)
		utils.SuccessResponse(c, gin.H{
			"user":    user,
			"stats":   stats,
			"reviews": nil, // можно добавить последние отзывы
		})
		return
	}

	utils.SuccessResponse(c, user)
}
