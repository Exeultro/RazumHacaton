package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"razum-backend/internal/models"
	"razum-backend/internal/repository"

	"github.com/google/uuid"
)

// parseTime парсит время из строки в поддерживаемых форматах
func parseTime(timeStr string) (time.Time, error) {
	// Поддерживаемые форматы
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid time format: %s", timeStr)
}

type EventService struct {
	eventRepo         *repository.EventRepository
	userRepo          *repository.UserRepository
	participationRepo *repository.ParticipationRepository
}

func NewEventService(
	eventRepo *repository.EventRepository,
	userRepo *repository.UserRepository,
	participationRepo *repository.ParticipationRepository,
) *EventService {
	return &EventService{
		eventRepo:         eventRepo,
		userRepo:          userRepo,
		participationRepo: participationRepo,
	}
}

// CreateEventRequest структура для создания мероприятия
type CreateEventRequest struct {
	Title                  string             `json:"title" binding:"required"`
	Description            string             `json:"description"`
	EventDate              string             `json:"event_date" binding:"required"`
	RegistrationDeadline   string             `json:"registration_deadline" binding:"required"`
	Location               *string            `json:"location"`
	Format                 models.EventFormat `json:"format" binding:"required,oneof=offline online hybrid"`
	Direction              models.Direction   `json:"direction" binding:"required,oneof=IT social media"`
	DifficultyCoefficient  float64            `json:"difficulty_coefficient" binding:"required,min=0.5,max=3.0"`
	PointsForParticipation int                `json:"points_for_participation" binding:"required,min=0"`
	Prizes                 []models.Prize     `json:"prizes"`
}

// CreateEvent создает новое мероприятие
func (s *EventService) CreateEvent(organizerID uuid.UUID, req *CreateEventRequest) (*models.Event, error) {
	// Проверяем, что пользователь - организатор ИЛИ админ
	user, err := s.userRepo.FindByID(organizerID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Разрешаем организаторам и админам
	if user.Role != models.RoleOrganizer && user.Role != models.RoleAdmin {
		return nil, errors.New("only organizers and admins can create events")
	}

	// Парсим даты
	eventDate, err := parseTime(req.EventDate)
	if err != nil {
		return nil, fmt.Errorf("invalid event_date: %w", err)
	}

	registrationDeadline, err := parseTime(req.RegistrationDeadline)
	if err != nil {
		return nil, fmt.Errorf("invalid registration_deadline: %w", err)
	}

	// Создаем мероприятие
	event := &models.Event{
		ID:                     uuid.New(),
		OrganizerID:            organizerID,
		Title:                  req.Title,
		Description:            req.Description,
		EventDate:              eventDate,
		RegistrationDeadline:   registrationDeadline,
		Location:               req.Location,
		Format:                 req.Format,
		Direction:              req.Direction,
		DifficultyCoefficient:  req.DifficultyCoefficient,
		PointsForParticipation: req.PointsForParticipation,
		Prizes:                 req.Prizes,
		Status:                 models.StatusPublished,
	}

	if err := s.eventRepo.Create(event); err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	// Если пользователь админ, нужно добавить запись в organizers (если нет)
	if user.Role == models.RoleAdmin {
		// Проверяем, есть ли запись в organizers
		exists, err := s.userRepo.CheckOrganizerExists(organizerID)
		if err != nil {
			// Логируем ошибку, но не блокируем создание мероприятия
			fmt.Printf("Warning: failed to check organizer exists: %v\n", err)
		} else if !exists {
			// Создаем запись в organizers для админа
			if err := s.userRepo.CreateOrganizer(organizerID); err != nil {
				fmt.Printf("Warning: failed to create organizer record: %v\n", err)
			}
		}
	}

	return event, nil
}

// GetEvent получает мероприятие по ID
func (s *EventService) GetEvent(eventID uuid.UUID) (*models.Event, error) {
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to find event: %w", err)
	}
	if event == nil {
		return nil, errors.New("event not found")
	}
	return event, nil
}

// ListEvents получает список мероприятий
func (s *EventService) ListEvents(status, direction, format string, page, limit int) ([]models.Event, int, error) {
	log.Printf("EventService.ListEvents: status=%s, direction=%s, format=%s, page=%d, limit=%d",
		status, direction, format, page, limit)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	return s.eventRepo.List(status, direction, format, page, limit)
}

// UpdateEvent обновляет мероприятие
func (s *EventService) UpdateEvent(eventID uuid.UUID, organizerID uuid.UUID, req *CreateEventRequest) (*models.Event, error) {
	// Получаем существующее мероприятие
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to find event: %w", err)
	}
	if event == nil {
		return nil, errors.New("event not found")
	}

	// Проверяем, что пользователь - организатор этого мероприятия
	if event.OrganizerID != organizerID {
		return nil, errors.New("you can only update your own events")
	}

	// Обновляем поля
	event.Title = req.Title
	event.Description = req.Description

	eventDate, err := parseTime(req.EventDate)
	if err != nil {
		return nil, fmt.Errorf("invalid event_date: %w", err)
	}
	event.EventDate = eventDate

	registrationDeadline, err := parseTime(req.RegistrationDeadline)
	if err != nil {
		return nil, fmt.Errorf("invalid registration_deadline: %w", err)
	}
	event.RegistrationDeadline = registrationDeadline

	event.Location = req.Location
	event.Format = req.Format
	event.Direction = req.Direction
	event.DifficultyCoefficient = req.DifficultyCoefficient
	event.PointsForParticipation = req.PointsForParticipation
	event.Prizes = req.Prizes

	if err := s.eventRepo.Update(event); err != nil {
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	return event, nil
}

// DeleteEvent удаляет мероприятие
func (s *EventService) DeleteEvent(eventID uuid.UUID, organizerID uuid.UUID) error {
	// Получаем существующее мероприятие
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return fmt.Errorf("failed to find event: %w", err)
	}
	if event == nil {
		return errors.New("event not found")
	}

	// Проверяем, что пользователь - организатор этого мероприятия
	if event.OrganizerID != organizerID {
		return errors.New("you can only delete your own events")
	}

	if err := s.eventRepo.Delete(eventID); err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return nil
}

// enrichEventWithParticipants добавляет информацию об участниках к событию
func (s *EventService) enrichEventWithParticipants(event *models.Event) (*models.EventWithParticipants, error) {
	log.Printf("Enriching event %s with participants", event.ID)
	counts, err := s.eventRepo.GetEventParticipantsCounts(event.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants counts: %w", err)
	}
	log.Printf("Event %s counts: total=%d, registered=%d, attended=%d",
		event.ID, counts["total"], counts["registered"], counts["attended"])

	return &models.EventWithParticipants{
		Event:             event,
		ParticipantsCount: counts["total"],
		RegisteredCount:   counts["registered"],
		AttendedCount:     counts["attended"],
		ConfirmedCount:    counts["confirmed"],
	}, nil
}

// ListEventsWithParticipants получает список мероприятий с количеством участников
func (s *EventService) ListEventsWithParticipants(status, direction, format string, page, limit int) ([]models.EventWithParticipants, int, error) {
	events, total, err := s.eventRepo.List(status, direction, format, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list events: %w", err)
	}

	eventsWithParticipants := make([]models.EventWithParticipants, 0, len(events))
	for i := 0; i < len(events); i++ {
		enriched, err := s.enrichEventWithParticipants(&events[i])
		if err != nil {
			continue
		}
		eventsWithParticipants = append(eventsWithParticipants, *enriched)
	}

	return eventsWithParticipants, total, nil
}

// GetEventWithParticipants получает одно мероприятие с участниками
func (s *EventService) GetEventWithParticipants(eventID uuid.UUID) (*models.EventWithParticipants, error) {
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to find event: %w", err)
	}
	if event == nil {
		return nil, errors.New("event not found")
	}

	return s.enrichEventWithParticipants(event)
}

// ToEventResponse преобразует EventWithParticipants в EventResponse
func (s *EventService) ToEventResponse(
	enriched *models.EventWithParticipants,
	userID *uuid.UUID,
) *models.EventResponse {
	log.Printf("ToEventResponse: called for event %s, userID=%v", enriched.ID, userID)

	// Получаем информацию об организаторе
	organizerInfo := models.OrganizerInfo{
		ID:          enriched.OrganizerID,
		FullName:    "",
		TrustRating: 0,
		EventsCount: 0,
	}

	// Получаем пользователя
	user, _ := s.userRepo.FindByID(enriched.OrganizerID)
	if user != nil {
		organizerInfo.FullName = user.FullName
	}

	// Получаем статистику организатора
	orgStats, _ := s.userRepo.GetOrganizerStats(enriched.OrganizerID)
	if orgStats != nil {
		organizerInfo.TrustRating = orgStats.Rating
		organizerInfo.EventsCount = orgStats.EventsCount
	}

	response := &models.EventResponse{
		ID:                     enriched.ID,
		Title:                  enriched.Title,
		Description:            enriched.Description,
		EventDate:              enriched.EventDate,
		RegistrationDeadline:   enriched.RegistrationDeadline,
		Location:               enriched.Location,
		Format:                 enriched.Format,
		Direction:              enriched.Direction,
		DifficultyCoefficient:  enriched.DifficultyCoefficient,
		PointsForParticipation: enriched.PointsForParticipation,
		Prizes:                 enriched.Prizes,
		Status:                 enriched.Status,
		Organizer:              organizerInfo,
		ParticipantsCount:      enriched.ParticipantsCount,
		RegisteredCount:        enriched.RegisteredCount,
		AttendedCount:          enriched.AttendedCount,
		ConfirmedCount:         enriched.ConfirmedCount,
		IsRegistered:           false,
		CreatedAt:              enriched.CreatedAt,
		UpdatedAt:              enriched.UpdatedAt,
	}

	// Если есть userID, проверяем статус регистрации
	if userID != nil {
		log.Printf("ToEventResponse: checking participation for user %s in event %s", userID, enriched.ID)
		participation, err := s.participationRepo.FindByUserAndEvent(*userID, enriched.ID)
		if err != nil {
			log.Printf("ToEventResponse: error finding participation: %v", err)
		}
		if participation != nil {
			log.Printf("ToEventResponse: found participation with status %s", participation.Status)
			response.IsRegistered = true
			statusStr := string(participation.Status)
			response.UserStatus = &statusStr
		} else {
			log.Printf("ToEventResponse: no participation found for user %s", userID)
		}
	} else {
		log.Printf("ToEventResponse: userID is nil")
	}

	return response
}

// ListEventsWithFilters получает список мероприятий с фильтрацией
func (s *EventService) ListEventsWithFilters(status, direction, format, dateFrom, dateTo, search string, page, limit int) ([]models.EventWithParticipants, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	events, total, err := s.eventRepo.ListWithFilters(status, direction, format, dateFrom, dateTo, search, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list events: %w", err)
	}

	eventsWithParticipants := make([]models.EventWithParticipants, 0, len(events))
	for i := range events {
		enriched, err := s.enrichEventWithParticipants(&events[i])
		if err != nil {
			continue
		}
		eventsWithParticipants = append(eventsWithParticipants, *enriched)
	}

	return eventsWithParticipants, total, nil
}
