package models

import (
	"time"

	"github.com/google/uuid"
)

type OrganizerReview struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	OrganizerID   uuid.UUID  `json:"organizer_id" db:"organizer_id"`
	ParticipantID uuid.UUID  `json:"participant_id" db:"participant_id"`
	EventID       *uuid.UUID `json:"event_id,omitempty" db:"event_id"`
	Rating        int        `json:"rating" db:"rating"`
	Comment       string     `json:"comment" db:"comment"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}

// ReviewWithDetails отзыв с деталями
type ReviewWithDetails struct {
	ID              uuid.UUID  `json:"id"`
	OrganizerID     uuid.UUID  `json:"organizer_id"`
	OrganizerName   string     `json:"organizer_name"`
	ParticipantID   uuid.UUID  `json:"participant_id"`
	ParticipantName string     `json:"participant_name"`
	EventID         *uuid.UUID `json:"event_id,omitempty"`
	EventTitle      *string    `json:"event_title,omitempty"`
	Rating          int        `json:"rating"`
	Comment         string     `json:"comment"`
	CreatedAt       time.Time  `json:"created_at"`
}

// OrganizerRatingStats статистика рейтинга организатора
type OrganizerRatingStats struct {
	OrganizerID   uuid.UUID `json:"organizer_id"`
	AverageRating float64   `json:"average_rating"`
	TotalReviews  int       `json:"total_reviews"`
	Rating5Count  int       `json:"rating_5_count"`
	Rating4Count  int       `json:"rating_4_count"`
	Rating3Count  int       `json:"rating_3_count"`
	Rating2Count  int       `json:"rating_2_count"`
	Rating1Count  int       `json:"rating_1_count"`
}
