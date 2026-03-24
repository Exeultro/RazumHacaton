package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleOrganizer   UserRole = "organizer"
	RoleParticipant UserRole = "participant"
	RoleObserver    UserRole = "observer"
	RoleAdmin       UserRole = "admin"
)

type Direction string

const (
	DirectionIT     Direction = "IT"
	DirectionSocial Direction = "social"
	DirectionMedia  Direction = "media"
)

type User struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Email     string     `json:"email" db:"email"`
	Password  string     `json:"-" db:"password"`
	FullName  string     `json:"full_name" db:"full_name"`
	Role      UserRole   `json:"role" db:"role"`
	City      *string    `json:"city,omitempty" db:"city"`
	Age       *int       `json:"age,omitempty" db:"age"`
	Direction *Direction `json:"direction,omitempty" db:"direction"`
	AvatarURL *string    `json:"avatar_url,omitempty" db:"avatar_url"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

type Organizer struct {
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	Rating       float64   `json:"rating" db:"rating"`
	EventsCount  int       `json:"events_count" db:"events_count"`
	CommonPrizes []string  `json:"common_prizes" db:"common_prizes"`
}
