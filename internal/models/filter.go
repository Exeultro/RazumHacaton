package models

import (
	"time"

	"github.com/google/uuid"
)

// ObserverFilter фильтр наблюдателя
type ObserverFilter struct {
	ID         uuid.UUID      `json:"id" db:"id"`
	UserID     uuid.UUID      `json:"user_id" db:"user_id"`
	FilterName string         `json:"filter_name" db:"filter_name"`
	Filters    FilterCriteria `json:"filters" db:"filters"`
	CreatedAt  time.Time      `json:"created_at" db:"created_at"`
}

// FilterCriteria критерии фильтрации
type FilterCriteria struct {
	AgeMin       *int     `json:"age_min,omitempty"`
	AgeMax       *int     `json:"age_max,omitempty"`
	City         *string  `json:"city,omitempty"`
	Direction    *string  `json:"direction,omitempty"`
	MinPoints    *int     `json:"min_points,omitempty"`
	MinEvents    *int     `json:"min_events,omitempty"`
	MinAvgPoints *float64 `json:"min_avg_points,omitempty"`
	SortBy       string   `json:"sort_by,omitempty"`
	SortOrder    string   `json:"sort_order,omitempty"`
}
