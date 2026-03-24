package handlers

import (
	"strconv"

	"razum-backend/internal/services"
	"razum-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReviewHandler struct {
	reviewService *services.ReviewService
}

func NewReviewHandler(reviewService *services.ReviewService) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
	}
}

// CreateReview создает отзыв об организаторе
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c)
		return
	}

	organizerIDParam := c.Param("id")
	organizerID, err := uuid.Parse(organizerIDParam)
	if err != nil {
		utils.BadRequest(c, "Неверный ID организатора")
		return
	}

	var req struct {
		EventID *string `json:"event_id,omitempty"`
		Rating  int     `json:"rating" binding:"required,min=1,max=5"`
		Comment string  `json:"comment" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Неверный формат запроса: "+err.Error())
		return
	}

	var eventID *uuid.UUID
	if req.EventID != nil && *req.EventID != "" {
		eid, err := uuid.Parse(*req.EventID)
		if err == nil {
			eventID = &eid
		}
	}

	review, err := h.reviewService.CreateReview(userID.(uuid.UUID), organizerID, eventID, req.Rating, req.Comment)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.CreatedResponse(c, review)
}

// GetOrganizerReviews получает отзывы об организаторе
func (h *ReviewHandler) GetOrganizerReviews(c *gin.Context) {
	organizerIDParam := c.Param("id")
	organizerID, err := uuid.Parse(organizerIDParam)
	if err != nil {
		utils.BadRequest(c, "Неверный ID организатора")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	reviews, total, err := h.reviewService.GetOrganizerReviews(organizerID, page, limit)
	if err != nil {
		utils.InternalServerError(c, "Ошибка при получении отзывов: "+err.Error())
		return
	}

	stats, _ := h.reviewService.GetOrganizerRatingStats(organizerID)

	utils.SuccessResponse(c, gin.H{
		"reviews": reviews,
		"stats":   stats,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + limit - 1) / limit,
		},
	})
}
