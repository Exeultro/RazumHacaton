package handlers

import (
	"razum-backend/internal/models"
	"razum-backend/internal/services"
	"razum-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ParticipationHandler struct {
	participationService *services.ParticipationService
	eventService         *services.EventService
}

func NewParticipationHandler(
	participationService *services.ParticipationService,
	eventService *services.EventService,
) *ParticipationHandler {
	return &ParticipationHandler{
		participationService: participationService,
		eventService:         eventService,
	}
}

// RegisterForEvent регистрирует текущего пользователя на мероприятие
// POST /api/events/:id/register
func (h *ParticipationHandler) RegisterForEvent(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c)
		return
	}

	eventIDParam := c.Param("id")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		utils.BadRequest(c, "Неверный ID мероприятия")
		return
	}

	participation, err := h.participationService.RegisterForEvent(userID.(uuid.UUID), eventID)
	if err != nil {
		switch err.Error() {
		case "event not found":
			utils.EventNotFound(c)
		case "event is cancelled":
			utils.EventCancelled(c)
		case "registration deadline has passed":
			utils.RegistrationDeadlinePassed(c)
		case "event has already started":
			utils.BadRequest(c, "Мероприятие уже началось")
		case "already registered for this event":
			utils.AlreadyRegistered(c)
		default:
			utils.BadRequest(c, err.Error())
		}
		return
	}

	utils.SuccessResponse(c, gin.H{
		"participation_id": participation.ID,
		"qr_code_token":    participation.QRCodeToken,
		"message":          "Вы успешно зарегистрированы на мероприятие",
	})
}

// GetMyQRCode возвращает QR-токен для отображения участнику
// GET /api/events/:id/my-qr
func (h *ParticipationHandler) GetMyQRCode(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c)
		return
	}

	eventIDParam := c.Param("id")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		utils.BadRequest(c, "Неверный ID мероприятия")
		return
	}

	token, err := h.participationService.GetQRTokenForEvent(userID.(uuid.UUID), eventID)
	if err != nil {
		switch err.Error() {
		case "not registered for this event":
			utils.BadRequest(c, "Вы не зарегистрированы на это мероприятие")
		case "participation already confirmed, QR code is no longer valid":
			utils.QRCodeAlreadyUsed(c)
		case "registration was cancelled":
			utils.BadRequest(c, "Регистрация была отменена")
		case "registration deadline has passed":
			utils.RegistrationDeadlinePassed(c)
		default:
			utils.BadRequest(c, err.Error())
		}
		return
	}

	utils.SuccessResponse(c, gin.H{
		"qr_code_token": token.String(),
		"message":       "Отсканируйте этот QR-код на мероприятии",
	})
}

// ConfirmParticipation подтверждает участие (сканирование QR-кода)
// POST /api/events/:id/confirm
func (h *ParticipationHandler) ConfirmParticipation(c *gin.Context) {
	organizerID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c)
		return
	}

	var req struct {
		QRCodeToken string `json:"qr_code_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Не передан QR-код")
		return
	}

	qrToken, err := uuid.Parse(req.QRCodeToken)
	if err != nil {
		utils.InvalidQRCode(c)
		return
	}

	participation, points, err := h.participationService.ConfirmParticipation(qrToken, organizerID.(uuid.UUID))
	if err != nil {
		switch err.Error() {
		case "invalid QR code":
			utils.InvalidQRCode(c)
		case "QR code already used, participation already confirmed":
			utils.QRCodeAlreadyUsed(c)
		case "only event organizer can confirm participation":
			utils.NotOrganizer(c)
		case "event was cancelled":
			utils.EventCancelled(c)
		default:
			utils.BadRequest(c, err.Error())
		}
		return
	}

	utils.SuccessResponse(c, gin.H{
		"participation_id": participation.ID,
		"points_earned":    points,
		"message":          "Участие подтверждено! Начислено баллов: " + string(rune(points)),
	})
}

// CancelRegistration отменяет регистрацию на мероприятие
// DELETE /api/events/:id/cancel
func (h *ParticipationHandler) CancelRegistration(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c)
		return
	}

	eventIDParam := c.Param("id")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		utils.BadRequest(c, "Неверный ID мероприятия")
		return
	}

	if err := h.participationService.CancelRegistration(userID.(uuid.UUID), eventID); err != nil {
		switch err.Error() {
		case "not registered for this event":
			utils.BadRequest(c, "Вы не зарегистрированы на это мероприятие")
		case "cannot cancel attendance after confirmation":
			utils.CannotCancelConfirmed(c)
		case "registration already cancelled":
			utils.BadRequest(c, "Регистрация уже отменена")
		default:
			utils.BadRequest(c, err.Error())
		}
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Регистрация успешно отменена"})
}

// GetEventParticipants возвращает список участников мероприятия
// GET /api/events/:id/participants
func (h *ParticipationHandler) GetEventParticipants(c *gin.Context) {
	eventIDParam := c.Param("id")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		utils.BadRequest(c, "Неверный ID мероприятия")
		return
	}

	participants, err := h.participationService.GetEventParticipants(eventID)
	if err != nil {
		utils.InternalServerError(c, "Ошибка при получении участников: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"participants": participants})
}

// GetMyParticipations возвращает список участий текущего пользователя
// GET /api/my/participations
func (h *ParticipationHandler) GetMyParticipations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c)
		return
	}

	uid := userID.(uuid.UUID)

	participations, err := h.participationService.GetUserParticipations(uid)
	if err != nil {
		utils.InternalServerError(c, "Ошибка при получении участий: "+err.Error())
		return
	}

	// Обогащаем каждое участие деталями мероприятия
	enrichedParticipations := make([]gin.H, 0, len(participations))
	for _, p := range participations {
		// Получаем мероприятие с участниками
		event, _ := h.eventService.GetEventWithParticipants(p.EventID)

		// Получаем EventResponse через ToEventResponse
		var eventResponse *models.EventResponse
		if event != nil {
			eventResponse = h.eventService.ToEventResponse(event, &uid)
		}

		enrichedParticipations = append(enrichedParticipations, gin.H{
			"id":            p.ID,
			"event_id":      p.EventID,
			"user_id":       p.UserID,
			"status":        p.Status,
			"qr_code_token": p.QRCodeToken,
			"points_earned": p.PointsEarned,
			"attended_at":   p.AttendedAt,
			"confirmed_by":  p.ConfirmedBy,
			"created_at":    p.CreatedAt,
			"event":         eventResponse,
		})
	}

	utils.SuccessResponse(c, gin.H{
		"participations": enrichedParticipations,
	})
}
