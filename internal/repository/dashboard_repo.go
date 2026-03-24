package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type DashboardRepository struct {
	db *sql.DB
}

func NewDashboardRepository(db *sql.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

// RecentEvent представляет мероприятие для ленты
type RecentEvent struct {
	ID                string    `json:"id"`
	Title             string    `json:"title"`
	EventDate         time.Time `json:"event_date"`
	OrganizerName     string    `json:"organizer_name"`
	OrganizerID       string    `json:"organizer_id"`
	ParticipantsCount int       `json:"participants_count"`
	Direction         string    `json:"direction"`
	Format            string    `json:"format"`
}

// GetRecentEvents возвращает последние 10 мероприятий
func (r *DashboardRepository) GetRecentEvents(limit int) ([]RecentEvent, error) {
	query := `
        SELECT 
            e.id, 
            e.title, 
            e.event_date, 
            u.full_name as organizer_name,
            u.id as organizer_id,
            e.direction,
            e.format,
            COUNT(ep.id) as participants_count
        FROM events e
        JOIN users u ON e.organizer_id = u.id
        LEFT JOIN event_participations ep ON e.id = ep.event_id AND ep.status = 'attended'
        WHERE e.status = 'published' AND e.event_date >= CURRENT_DATE
        GROUP BY e.id, u.full_name, u.id
        ORDER BY e.event_date ASC
        LIMIT $1
    `

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent events: %w", err)
	}
	defer rows.Close()

	events := []RecentEvent{}
	for rows.Next() {
		var event RecentEvent
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.EventDate,
			&event.OrganizerName,
			&event.OrganizerID,
			&event.Direction,
			&event.Format,
			&event.ParticipantsCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}

// RatingHistoryPoint представляет точку графика рейтинга
type RatingHistoryPoint struct {
	Date   string `json:"date"`
	Points int    `json:"points"`
}

// GetRatingHistory возвращает историю роста рейтинга участника
func (r *DashboardRepository) GetRatingHistory(userID uuid.UUID) ([]RatingHistoryPoint, error) {
	query := `
        SELECT 
            DATE(pa.created_at) as date,
            SUM(pa.points_change) OVER (ORDER BY DATE(pa.created_at)) as cumulative_points
        FROM points_audit pa
        WHERE pa.user_id = $1
        ORDER BY date ASC
    `

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get rating history: %w", err)
	}
	defer rows.Close()

	history := []RatingHistoryPoint{}
	for rows.Next() {
		var point RatingHistoryPoint
		var date time.Time
		var points int

		err := rows.Scan(&date, &points)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history point: %w", err)
		}

		point.Date = date.Format("2006-01-02")
		point.Points = points
		history = append(history, point)
	}

	return history, nil
}

// TagStat представляет статистику по тегу
type TagStat struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}

// GetTrendingTags возвращает облако тегов (популярные направления)
func (r *DashboardRepository) GetTrendingTags(limit int) ([]TagStat, error) {
	query := `
        SELECT 
            direction as tag,
            COUNT(*) as count
        FROM events
        WHERE status = 'published'
        GROUP BY direction
        ORDER BY count DESC
        LIMIT $1
    `

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get trending tags: %w", err)
	}
	defer rows.Close()

	tags := []TagStat{}
	for rows.Next() {
		var tag TagStat
		err := rows.Scan(&tag.Tag, &tag.Count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// GetActivityStats возвращает общую статистику активности
type ActivityStats struct {
	TotalEvents       int `json:"total_events"`
	TotalParticipants int `json:"total_participants"`
	TotalPoints       int `json:"total_points"`
	ActiveOrganizers  int `json:"active_organizers"`
}

func (r *DashboardRepository) GetActivityStats() (*ActivityStats, error) {
	query := `
        SELECT 
            (SELECT COUNT(*) FROM events WHERE status = 'published') as total_events,
            (SELECT COUNT(DISTINCT user_id) FROM event_participations WHERE status = 'attended') as total_participants,
            (SELECT COALESCE(SUM(points_earned), 0) FROM event_participations WHERE status = 'attended') as total_points,
            (SELECT COUNT(DISTINCT organizer_id) FROM events WHERE status = 'published') as active_organizers
    `

	var stats ActivityStats
	err := r.db.QueryRow(query).Scan(
		&stats.TotalEvents,
		&stats.TotalParticipants,
		&stats.TotalPoints,
		&stats.ActiveOrganizers,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity stats: %w", err)
	}

	return &stats, nil
}
