package models

import (
	"time"

	"github.com/google/uuid"
)

type ParticipationStatus string

const (
	ParticipationRegistered ParticipationStatus = "registered"
	ParticipationConfirmed  ParticipationStatus = "confirmed"
	ParticipationAttended   ParticipationStatus = "attended"
	ParticipationCancelled  ParticipationStatus = "cancelled"
)

type EventParticipation struct {
	ID           uuid.UUID           `json:"id" db:"id"`
	EventID      uuid.UUID           `json:"event_id" db:"event_id"`
	UserID       uuid.UUID           `json:"user_id" db:"user_id"`
	Status       ParticipationStatus `json:"status" db:"status"`
	QRCodeToken  uuid.UUID           `json:"qr_code_token" db:"qr_code_token"`
	PointsEarned int                 `json:"points_earned" db:"points_earned"`
	AttendedAt   *time.Time          `json:"attended_at,omitempty" db:"attended_at"`
	ConfirmedBy  *uuid.UUID          `json:"confirmed_by,omitempty" db:"confirmed_by"`
	CreatedAt    time.Time           `json:"created_at" db:"created_at"`
}
