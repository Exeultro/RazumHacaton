package models

import (
	"time"

	"github.com/google/uuid"
)

type EventFormat string

const (
	FormatOffline EventFormat = "offline"
	FormatOnline  EventFormat = "online"
	FormatHybrid  EventFormat = "hybrid"
)

// EventResponse унифицированный ответ для мероприятий
type EventResponse struct {
	ID                     uuid.UUID     `json:"id"`
	Title                  string        `json:"title"`
	Description            string        `json:"description"`
	EventDate              time.Time     `json:"event_date"`
	RegistrationDeadline   time.Time     `json:"registration_deadline"`
	Location               *string       `json:"location,omitempty"`
	Format                 EventFormat   `json:"format"`
	Direction              Direction     `json:"direction"`
	DifficultyCoefficient  float64       `json:"difficulty_coefficient"`
	PointsForParticipation int           `json:"points_for_participation"`
	Prizes                 []Prize       `json:"prizes"`
	Status                 EventStatus   `json:"status"`
	Organizer              OrganizerInfo `json:"organizer"`
	ParticipantsCount      int           `json:"participants_count"`
	RegisteredCount        int           `json:"registered_count"`
	AttendedCount          int           `json:"attended_count"`
	ConfirmedCount         int           `json:"confirmed_count"`
	IsRegistered           bool          `json:"is_registered"`
	UserStatus             *string       `json:"user_status,omitempty"`
	CreatedAt              time.Time     `json:"created_at"`
	UpdatedAt              time.Time     `json:"updated_at"`
}

// OrganizerInfo информация об организаторе для ответа
type OrganizerInfo struct {
	ID          uuid.UUID `json:"id"`
	FullName    string    `json:"full_name"`
	TrustRating float64   `json:"trust_rating"`
	EventsCount int       `json:"events_count"`
}
type EventStatus string

const (
	StatusDraft     EventStatus = "draft"
	StatusPublished EventStatus = "published"
	StatusCompleted EventStatus = "completed"
	StatusCancelled EventStatus = "cancelled"
)

type Prize struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type EventWithParticipants struct {
	*Event
	ParticipantsCount int `json:"participants_count"`
	RegisteredCount   int `json:"registered_count"`
	AttendedCount     int `json:"attended_count"`
	ConfirmedCount    int `json:"confirmed_count"`
}

type Event struct {
	ID                     uuid.UUID   `json:"id" db:"id"`
	OrganizerID            uuid.UUID   `json:"organizer_id" db:"organizer_id"`
	Title                  string      `json:"title" db:"title"`
	Description            string      `json:"description" db:"description"`
	EventDate              time.Time   `json:"event_date" db:"event_date"`
	RegistrationDeadline   time.Time   `json:"registration_deadline" db:"registration_deadline"`
	Location               *string     `json:"location,omitempty" db:"location"`
	Format                 EventFormat `json:"format" db:"format"`
	Direction              Direction   `json:"direction" db:"direction"`
	DifficultyCoefficient  float64     `json:"difficulty_coefficient" db:"difficulty_coefficient"`
	PointsForParticipation int         `json:"points_for_participation" db:"points_for_participation"`
	Prizes                 []Prize     `json:"prizes" db:"prizes"`
	Status                 EventStatus `json:"status" db:"status"`
	CreatedAt              time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time   `json:"updated_at" db:"updated_at"`
	OrganizerName          string      `json:"organizer_name,omitempty" db:"-"`
	OrganizerRating        float64     `json:"organizer_rating,omitempty" db:"-"`
}
