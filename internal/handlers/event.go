package handlers

import (
	"strconv"

	"razum-backend/internal/models"
	"razum-backend/internal/services"
	"razum-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EventHandler struct {
	eventService *services.EventService
}

func NewEventHandler(eventService *services.EventService) *EventHandler {
	return &EventHandler{
		eventService: eventService,
	}
}

// CreateEvent создает мероприятие
func (h *EventHandler) CreateEvent(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c)
		return
	}

	var req services.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	event, err := h.eventService.CreateEvent(userID.(uuid.UUID), &req)
	if err != nil {
		utils.InternalServerError(c, "Failed to create event: "+err.Error())
		return
	}

	utils.CreatedResponse(c, event)
}

// GetEvent получает мероприятие по ID
func (h *EventHandler) GetEvent(c *gin.Context) {
	idParam := c.Param("id")
	eventID, err := uuid.Parse(idParam)
	if err != nil {
		utils.BadRequest(c, "Invalid event ID")
		return
	}

	enriched, err := h.eventService.GetEventWithParticipants(eventID)
	if err != nil {
		if err.Error() == "event not found" {
			utils.NotFound(c)
			return
		}
		utils.InternalServerError(c, "Failed to get event: "+err.Error())
		return
	}

	// Получаем userID из контекста (если есть)
	var userID *uuid.UUID
	if uid, exists := c.Get("user_id"); exists {
		if id, ok := uid.(uuid.UUID); ok {
			userID = &id
		}
	}

	response := h.eventService.ToEventResponse(enriched, userID)
	utils.SuccessResponse(c, response)
}

// ListEvents получает список мероприятий
func (h *EventHandler) ListEvents(c *gin.Context) {
	status := c.Query("status")
	direction := c.Query("direction")
	format := c.Query("format")
	dateFrom := c.Query("dateFrom")
	dateTo := c.Query("dateTo")
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	events, total, err := h.eventService.ListEventsWithFilters(status, direction, format, dateFrom, dateTo, search, page, limit)
	if err != nil {
		utils.InternalServerError(c, "Failed to list events: "+err.Error())
		return
	}

	// Получаем userID из контекста (если есть)
	var userID *uuid.UUID
	if uid, exists := c.Get("user_id"); exists {
		if id, ok := uid.(uuid.UUID); ok {
			userID = &id
		}
	}

	responses := make([]*models.EventResponse, 0, len(events))
	for _, event := range events {
		responses = append(responses, h.eventService.ToEventResponse(&event, userID))
	}

	utils.SuccessResponse(c, gin.H{
		"events": responses,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + limit - 1) / limit,
		},
	})
}

// UpdateEvent обновляет мероприятие
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c)
		return
	}

	idParam := c.Param("id")
	eventID, err := uuid.Parse(idParam)
	if err != nil {
		utils.BadRequest(c, "Invalid event ID")
		return
	}

	var req services.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	event, err := h.eventService.UpdateEvent(eventID, userID.(uuid.UUID), &req)
	if err != nil {
		utils.InternalServerError(c, "Failed to update event: "+err.Error())
		return
	}

	utils.SuccessResponse(c, event)
}

// DeleteEvent удаляет мероприятие
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c)
		return
	}

	idParam := c.Param("id")
	eventID, err := uuid.Parse(idParam)
	if err != nil {
		utils.BadRequest(c, "Invalid event ID")
		return
	}

	if err := h.eventService.DeleteEvent(eventID, userID.(uuid.UUID)); err != nil {
		utils.InternalServerError(c, "Failed to delete event: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Event deleted successfully"})
}
