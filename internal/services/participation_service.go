package services

import (
	"errors"
	"fmt"
	"time"

	"razum-backend/internal/models"
	"razum-backend/internal/repository"

	"github.com/google/uuid"
)

type ParticipationService struct {
	participationRepo *repository.ParticipationRepository
	eventRepo         *repository.EventRepository
	userRepo          *repository.UserRepository
	auditRepo         *repository.AuditRepository
	ratingRepo        *repository.RatingRepository
}

func NewParticipationService(
	participationRepo *repository.ParticipationRepository,
	eventRepo *repository.EventRepository,
	userRepo *repository.UserRepository,
	auditRepo *repository.AuditRepository,
	ratingRepo *repository.RatingRepository,
) *ParticipationService {
	return &ParticipationService{
		participationRepo: participationRepo,
		eventRepo:         eventRepo,
		userRepo:          userRepo,
		auditRepo:         auditRepo,
		ratingRepo:        ratingRepo,
	}
}

// RegisterForEvent регистрирует участника на мероприятие
func (s *ParticipationService) RegisterForEvent(userID, eventID uuid.UUID) (*models.EventParticipation, error) {
	// Проверяем существование мероприятия
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to find event: %w", err)
	}
	if event == nil {
		return nil, errors.New("event not found")
	}

	// Проверяем, что мероприятие не отменено
	if event.Status == models.StatusCancelled {
		return nil, errors.New("event is cancelled")
	}

	// Проверяем, что дедлайн регистрации не истек
	if time.Now().After(event.RegistrationDeadline) {
		return nil, errors.New("registration deadline has passed")
	}

	// Проверяем, что мероприятие еще не началось
	if time.Now().After(event.EventDate) {
		return nil, errors.New("event has already started")
	}

	// Проверяем, не зарегистрирован ли пользователь уже
	existing, err := s.participationRepo.FindByUserAndEvent(userID, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing registration: %w", err)
	}
	if existing != nil {
		return nil, errors.New("already registered for this event")
	}

	// Создаем запись об участии
	participation := &models.EventParticipation{
		ID:           uuid.New(),
		EventID:      eventID,
		UserID:       userID,
		Status:       models.ParticipationRegistered,
		QRCodeToken:  uuid.New(),
		PointsEarned: 0,
	}

	if err := s.participationRepo.Create(participation); err != nil {
		return nil, fmt.Errorf("failed to create participation: %w", err)
	}

	return participation, nil
}

// ConfirmParticipation подтверждает участие организатором через QR-код
func (s *ParticipationService) ConfirmParticipation(qrToken uuid.UUID, organizerID uuid.UUID) (*models.EventParticipation, int, error) {
	// Находим запись об участии по QR-токену
	participation, err := s.participationRepo.FindByQRToken(qrToken)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find participation: %w", err)
	}
	if participation == nil {
		return nil, 0, errors.New("invalid QR code")
	}

	// Защита от повторного использования QR-кода
	if participation.Status == models.ParticipationAttended {
		return nil, 0, errors.New("QR code already used, participation already confirmed")
	}

	// Проверяем, что регистрация не была отменена
	if participation.Status == models.ParticipationCancelled {
		return nil, 0, errors.New("participation was cancelled")
	}

	// Получаем информацию о мероприятии
	event, err := s.eventRepo.FindByID(participation.EventID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find event: %w", err)
	}
	if event == nil {
		return nil, 0, errors.New("event not found")
	}

	// Проверяем, что подтверждающий является организатором этого мероприятия
	if event.OrganizerID != organizerID {
		return nil, 0, errors.New("only event organizer can confirm participation")
	}

	// Проверяем, что мероприятие не отменено
	if event.Status == models.StatusCancelled {
		return nil, 0, errors.New("event was cancelled")
	}

	// Проверяем, что мероприятие еще не завершено
	if event.Status == models.StatusCompleted {
		return nil, 0, errors.New("event already completed")
	}

	// Рассчитываем баллы: базовые баллы * коэффициент сложности
	pointsEarned := int(float64(event.PointsForParticipation) * event.DifficultyCoefficient)

	// Обновляем статус участия и начисляем баллы
	if err := s.participationRepo.UpdateStatus(
		participation.ID,
		models.ParticipationAttended,
		pointsEarned,
		organizerID,
	); err != nil {
		return nil, 0, fmt.Errorf("failed to update participation status: %w", err)
	}

	// Записываем начисление баллов в аудит
	audit := &models.PointsAudit{
		ID:           uuid.New(),
		UserID:       participation.UserID,
		EventID:      &participation.EventID,
		PointsChange: pointsEarned,
		Reason:       "participation_confirmed",
		CreatedAt:    time.Now(),
	}
	if err := s.auditRepo.AddPointsAudit(audit); err != nil {
		_ = err
	}

	participation.Status = models.ParticipationAttended
	participation.PointsEarned = pointsEarned
	if s.ratingRepo != nil {
		// Д обновим весь кэш
		if err := s.ratingRepo.UpdateRatingCache(); err != nil {
			// Логируем ошибку, но не блокируем основную операцию
			_ = err
		}
	}
	return participation, pointsEarned, nil
}

// CancelRegistration отменяет регистрацию на мероприятие
func (s *ParticipationService) CancelRegistration(userID, eventID uuid.UUID) error {
	// Находим запись об участии
	participation, err := s.participationRepo.FindByUserAndEvent(userID, eventID)
	if err != nil {
		return fmt.Errorf("failed to find participation: %w", err)
	}
	if participation == nil {
		return errors.New("not registered for this event")
	}

	// Нельзя отменить уже подтвержденное участие
	if participation.Status == models.ParticipationAttended {
		return errors.New("cannot cancel attendance after confirmation")
	}

	// Нельзя отменить уже отмененное
	if participation.Status == models.ParticipationCancelled {
		return errors.New("registration already cancelled")
	}

	// Обновляем статус на отмененный
	if err := s.participationRepo.UpdateStatus(
		participation.ID,
		models.ParticipationCancelled,
		0,
		uuid.Nil,
	); err != nil {
		return fmt.Errorf("failed to cancel registration: %w", err)
	}

	return nil
}

// GetUserParticipations возвращает список участий пользователя
func (s *ParticipationService) GetUserParticipations(userID uuid.UUID) ([]models.EventParticipation, error) {
	participations, err := s.participationRepo.GetUserParticipations(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user participations: %w", err)
	}
	return participations, nil
}

// GetEventParticipants возвращает список участников мероприятия
func (s *ParticipationService) GetEventParticipants(eventID uuid.UUID) ([]models.EventParticipation, error) {
	participants, err := s.participationRepo.GetEventParticipants(eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event participants: %w", err)
	}
	return participants, nil
}

// GetConfirmedParticipantsCount возвращает количество подтвержденных участников мероприятия
func (s *ParticipationService) GetConfirmedParticipantsCount(eventID uuid.UUID) (int, error) {
	count, err := s.participationRepo.GetConfirmedParticipantsCount(eventID)
	if err != nil {
		return 0, fmt.Errorf("failed to get confirmed participants count: %w", err)
	}
	return count, nil
}

func (s *ParticipationService) GetParticipationByUserAndEvent(userID, eventID uuid.UUID) (*models.EventParticipation, error) {
	participation, err := s.participationRepo.FindByUserAndEvent(userID, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to find participation: %w", err)
	}
	return participation, nil
}

// GetQRTokenForEvent возвращает QR-токен для участника (для отображения)
func (s *ParticipationService) GetQRTokenForEvent(userID, eventID uuid.UUID) (*uuid.UUID, error) {
	// Находим участие
	participation, err := s.participationRepo.FindByUserAndEvent(userID, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to find participation: %w", err)
	}
	if participation == nil {
		return nil, errors.New("not registered for this event")
	}

	// Проверяем, что участие еще не подтверждено
	if participation.Status == models.ParticipationAttended {
		return nil, errors.New("participation already confirmed, QR code is no longer valid")
	}

	// Проверяем, что регистрация не отменена
	if participation.Status == models.ParticipationCancelled {
		return nil, errors.New("registration was cancelled")
	}

	// Проверяем, что дедлайн регистрации не прошел (опционально)
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to find event: %w", err)
	}
	if event == nil {
		return nil, errors.New("event not found")
	}

	// Проверяем, что дедлайн регистрации не истек
	if time.Now().After(event.RegistrationDeadline) {
		return nil, errors.New("registration deadline has passed")
	}

	return &participation.QRCodeToken, nil
}
