package handlers

import (
	"razum-backend/internal/services"
	"razum-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RatingHandler struct {
	ratingService *services.RatingService
}

func NewRatingHandler(ratingService *services.RatingService) *RatingHandler {
	return &RatingHandler{
		ratingService: ratingService,
	}
}

// GetGlobalRating возвращает глобальный рейтинг
func (h *RatingHandler) GetGlobalRating(c *gin.Context) {
	page := 1
	limit := 20

	if p := c.Query("page"); p != "" {
		// парсим page
	}
	if l := c.Query("limit"); l != "" {
		// парсим limit
	}

	ratings, total, err := h.ratingService.GetGlobalRating(page, limit)
	if err != nil {
		utils.InternalServerError(c, "Ошибка при получении рейтинга: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"rating": ratings,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + limit - 1) / limit,
		},
	})
}

// GetRatingByDirection возвращает рейтинг по направлению
func (h *RatingHandler) GetRatingByDirection(c *gin.Context) {
	direction := c.Param("direction")
	if direction == "" {
		utils.BadRequest(c, "Направление не указано")
		return
	}

	page := 1
	limit := 20

	ratings, total, err := h.ratingService.GetRatingByDirection(direction, page, limit)
	if err != nil {
		utils.InternalServerError(c, "Ошибка при получении рейтинга: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"rating": ratings,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + limit - 1) / limit,
		},
	})
}

// GetUserRating возвращает рейтинг пользователя
func (h *RatingHandler) GetUserRating(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		utils.BadRequest(c, "Неверный ID пользователя")
		return
	}

	rating, err := h.ratingService.GetUserRating(userID)
	if err != nil {
		utils.InternalServerError(c, "Ошибка при получении рейтинга: "+err.Error())
		return
	}

	utils.SuccessResponse(c, rating)
}

// GetMyRating возвращает рейтинг текущего пользователя
func (h *RatingHandler) GetMyRating(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c)
		return
	}

	rating, err := h.ratingService.GetMyRating(userID.(uuid.UUID))
	if err != nil {
		utils.InternalServerError(c, "Ошибка при получении рейтинга: "+err.Error())
		return
	}

	utils.SuccessResponse(c, rating)
}

// RefreshRatingCache обновляет кэш рейтинга (только для админа)
func (h *RatingHandler) RefreshRatingCache(c *gin.Context) {
	if err := h.ratingService.UpdateRatingCache(); err != nil {
		utils.InternalServerError(c, "Ошибка при обновлении кэша: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message": "Кэш рейтинга успешно обновлен",
	})
}
