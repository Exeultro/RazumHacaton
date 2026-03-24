package models

import (
	"time"

	"github.com/google/uuid"
)

type RatingCache struct {
    UserID        uuid.UUID  `json:"user_id" db:"user_id"`
    FullName      string     `json:"full_name" db:"-"`                 
    AvatarURL     *string    `json:"avatar_url,omitempty" db:"-"`        
    Direction     *string    `json:"direction,omitempty" db:"-"`         
    TotalPoints   int        `json:"total_points" db:"total_points"`
    EventsCount   int        `json:"events_count" db:"events_count"`
    TopDirection  *string    `json:"top_direction,omitempty" db:"top_direction"`
    GlobalRank    *int       `json:"global_rank" db:"global_rank"`
    ITRank        *int       `json:"it_rank,omitempty" db:"it_rank"`
    SocialRank    *int       `json:"social_rank,omitempty" db:"social_rank"`
    MediaRank     *int       `json:"media_rank,omitempty" db:"media_rank"`
    UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

type PointsAudit struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	EventID      *uuid.UUID `json:"event_id,omitempty" db:"event_id"`
	PointsChange int        `json:"points_change" db:"points_change"`
	Reason       string     `json:"reason" db:"reason"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
}
