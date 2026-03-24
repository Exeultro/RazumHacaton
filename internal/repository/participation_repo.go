package repository

import (
	"database/sql"
	"fmt"

	"razum-backend/internal/models"

	"github.com/google/uuid"
)

type ParticipationRepository struct {
	db *sql.DB
}

func NewParticipationRepository(db *sql.DB) *ParticipationRepository {
	return &ParticipationRepository{db: db}
}

// Create создает запись об участии
func (r *ParticipationRepository) Create(participation *models.EventParticipation) error {
	query := `
        INSERT INTO event_participations (
            id, event_id, user_id, status, qr_code_token, points_earned, created_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING created_at
    `

	err := r.db.QueryRow(
		query,
		participation.ID,
		participation.EventID,
		participation.UserID,
		participation.Status,
		participation.QRCodeToken,
		participation.PointsEarned,
		participation.CreatedAt,
	).Scan(&participation.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create participation: %w", err)
	}

	return nil
}

// FindByUserAndEvent находит участие по пользователю и событию
func (r *ParticipationRepository) FindByUserAndEvent(userID, eventID uuid.UUID) (*models.EventParticipation, error) {
	query := `
        SELECT id, event_id, user_id, status, qr_code_token, points_earned, attended_at, confirmed_by, created_at
        FROM event_participations
        WHERE user_id = $1 AND event_id = $2
    `

	participation := &models.EventParticipation{}
	err := r.db.QueryRow(query, userID, eventID).Scan(
		&participation.ID,
		&participation.EventID,
		&participation.UserID,
		&participation.Status,
		&participation.QRCodeToken,
		&participation.PointsEarned,
		&participation.AttendedAt,
		&participation.ConfirmedBy,
		&participation.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find participation: %w", err)
	}

	return participation, nil
}

// FindByQRToken находит участие по QR токену
func (r *ParticipationRepository) FindByQRToken(token uuid.UUID) (*models.EventParticipation, error) {
	query := `
        SELECT id, event_id, user_id, status, qr_code_token, points_earned, attended_at, confirmed_by, created_at
        FROM event_participations
        WHERE qr_code_token = $1
    `

	participation := &models.EventParticipation{}
	err := r.db.QueryRow(query, token).Scan(
		&participation.ID,
		&participation.EventID,
		&participation.UserID,
		&participation.Status,
		&participation.QRCodeToken,
		&participation.PointsEarned,
		&participation.AttendedAt,
		&participation.ConfirmedBy,
		&participation.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find participation by QR token: %w", err)
	}

	return participation, nil
}

// UpdateStatus обновляет статус участия
func (r *ParticipationRepository) UpdateStatus(participationID uuid.UUID, status models.ParticipationStatus, pointsEarned int, confirmedBy uuid.UUID) error {
	query := `
        UPDATE event_participations
        SET status = $1, points_earned = $2, confirmed_by = $3, attended_at = CURRENT_TIMESTAMP
        WHERE id = $4
    `

	_, err := r.db.Exec(query, status, pointsEarned, confirmedBy, participationID)
	if err != nil {
		return fmt.Errorf("failed to update participation status: %w", err)
	}

	return nil
}

// GetEventParticipants получает список участников мероприятия
func (r *ParticipationRepository) GetEventParticipants(eventID uuid.UUID) ([]models.EventParticipation, error) {
	query := `
        SELECT id, event_id, user_id, status, qr_code_token, points_earned, attended_at, confirmed_by, created_at
        FROM event_participations
        WHERE event_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.Query(query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event participants: %w", err)
	}
	defer rows.Close()

	participants := []models.EventParticipation{}
	for rows.Next() {
		var p models.EventParticipation
		err := rows.Scan(
			&p.ID,
			&p.EventID,
			&p.UserID,
			&p.Status,
			&p.QRCodeToken,
			&p.PointsEarned,
			&p.AttendedAt,
			&p.ConfirmedBy,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan participant: %w", err)
		}
		participants = append(participants, p)
	}

	return participants, nil
}

// GetUserParticipations возвращает список участий пользователя
func (r *ParticipationRepository) GetUserParticipations(userID uuid.UUID) ([]models.EventParticipation, error) {
	query := `
        SELECT id, event_id, user_id, status, qr_code_token, points_earned, attended_at, confirmed_by, created_at
        FROM event_participations
        WHERE user_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user participations: %w", err)
	}
	defer rows.Close()

	participations := []models.EventParticipation{}
	for rows.Next() {
		var p models.EventParticipation
		err := rows.Scan(
			&p.ID,
			&p.EventID,
			&p.UserID,
			&p.Status,
			&p.QRCodeToken,
			&p.PointsEarned,
			&p.AttendedAt,
			&p.ConfirmedBy,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan participation: %w", err)
		}
		participations = append(participations, p)
	}

	return participations, nil
}

// GetConfirmedParticipantsCount возвращает количество подтвержденных участников мероприятия
func (r *ParticipationRepository) GetConfirmedParticipantsCount(eventID uuid.UUID) (int, error) {
	query := `
        SELECT COUNT(*)
        FROM event_participations
        WHERE event_id = $1 AND status = $2
    `

	var count int
	err := r.db.QueryRow(query, eventID, models.ParticipationAttended).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get confirmed participants count: %w", err)
	}

	return count, nil
}
