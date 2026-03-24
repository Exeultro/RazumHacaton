package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// CandidateFilter фильтр для поиска кандидатов
type CandidateFilter struct {
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

// CandidateInfo информация о кандидате для кадрового резерва
type CandidateInfo struct {
	UserID        uuid.UUID `json:"user_id"`
	FullName      string    `json:"full_name"`
	Email         string    `json:"email"`
	Age           *int      `json:"age,omitempty"`
	City          *string   `json:"city,omitempty"`
	Direction     *string   `json:"direction,omitempty"`
	TotalPoints   int       `json:"total_points"`
	EventsCount   int       `json:"events_count"`
	AvgPoints     float64   `json:"avg_points"`
	GlobalRank    *int      `json:"global_rank,omitempty"`
	DirectionRank *int      `json:"direction_rank,omitempty"`
}

type CadreRepository struct {
	db *sql.DB
}

func NewCadreRepository(db *sql.DB) *CadreRepository {
	return &CadreRepository{db: db}
}

// GetCandidates возвращает список кандидатов с фильтрацией
func (r *CadreRepository) GetCandidates(filter CandidateFilter, page, limit int) ([]CandidateInfo, int, error) {
	offset := (page - 1) * limit

	// Базовый запрос с вычислением среднего балла
	baseQuery := `
        FROM users u
        LEFT JOIN rating_cache rc ON u.id = rc.user_id
        WHERE u.role = 'participant'
    `

	args := []interface{}{}
	argIdx := 1
	whereClauses := []string{}

	// Добавляем фильтры
	if filter.AgeMin != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("u.age >= $%d", argIdx))
		args = append(args, *filter.AgeMin)
		argIdx++
	}

	if filter.AgeMax != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("u.age <= $%d", argIdx))
		args = append(args, *filter.AgeMax)
		argIdx++
	}

	if filter.City != nil && *filter.City != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("u.city = $%d", argIdx))
		args = append(args, *filter.City)
		argIdx++
	}

	if filter.Direction != nil && *filter.Direction != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("u.direction = $%d", argIdx))
		args = append(args, *filter.Direction)
		argIdx++
	}

	if filter.MinPoints != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("rc.total_points >= $%d", argIdx))
		args = append(args, *filter.MinPoints)
		argIdx++
	}

	if filter.MinEvents != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("rc.events_count >= $%d", argIdx))
		args = append(args, *filter.MinEvents)
		argIdx++
	}

	if filter.MinAvgPoints != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("(rc.total_points::float / NULLIF(rc.events_count, 0)) >= $%d", argIdx))
		args = append(args, *filter.MinAvgPoints)
		argIdx++
	}

	// Добавляем where условия
	if len(whereClauses) > 0 {
		baseQuery += " AND " + strings.Join(whereClauses, " AND ")
	}

	// Считаем общее количество
	countQuery := "SELECT COUNT(*)" + baseQuery
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count candidates: %w", err)
	}

	if total == 0 {
		return []CandidateInfo{}, 0, nil
	}

	// Определяем сортировку
	sortBy := "rc.total_points"
	switch filter.SortBy {
	case "points":
		sortBy = "rc.total_points"
	case "events":
		sortBy = "rc.events_count"
	case "age":
		sortBy = "u.age"
	case "avg_points":
		sortBy = "avg_points"
	default:
		sortBy = "rc.total_points"
	}

	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	// Основной запрос с данными
	query := `
        SELECT 
            u.id,
            u.full_name,
            u.email,
            u.age,
            u.city,
            u.direction,
            COALESCE(rc.total_points, 0) as total_points,
            COALESCE(rc.events_count, 0) as events_count,
            COALESCE(rc.global_rank, 0) as global_rank,
            CASE 
                WHEN rc.events_count > 0 THEN rc.total_points::float / rc.events_count
                ELSE 0
            END as avg_points
    ` + baseQuery + fmt.Sprintf(`
        ORDER BY %s %s
        LIMIT $%d OFFSET $%d
    `, sortBy, sortOrder, argIdx, argIdx+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get candidates: %w", err)
	}
	defer rows.Close()

	candidates := []CandidateInfo{}
	for rows.Next() {
		var c CandidateInfo
		var globalRank int
		var direction sql.NullString

		err := rows.Scan(
			&c.UserID,
			&c.FullName,
			&c.Email,
			&c.Age,
			&c.City,
			&direction,
			&c.TotalPoints,
			&c.EventsCount,
			&globalRank,
			&c.AvgPoints,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan candidate: %w", err)
		}

		if direction.Valid {
			c.Direction = &direction.String
		}
		if globalRank > 0 {
			c.GlobalRank = &globalRank
		}

		candidates = append(candidates, c)
	}

	// Обновляем ранги по направлению (если нужно)
	if filter.Direction != nil && *filter.Direction != "" {
		for i := range candidates {
			rank := i + 1 + offset
			candidates[i].DirectionRank = &rank
		}
	}

	return candidates, total, nil
}

// GetCandidateByID получает кандидата по ID
func (r *CadreRepository) GetCandidateByID(userID uuid.UUID) (*CandidateInfo, error) {
	query := `
        SELECT 
            u.id,
            u.full_name,
            u.email,
            u.age,
            u.city,
            u.direction,
            COALESCE(rc.total_points, 0) as total_points,
            COALESCE(rc.events_count, 0) as events_count,
            COALESCE(rc.global_rank, 0) as global_rank,
            CASE 
                WHEN rc.events_count > 0 THEN rc.total_points::float / rc.events_count
                ELSE 0
            END as avg_points
        FROM users u
        LEFT JOIN rating_cache rc ON u.id = rc.user_id
        WHERE u.id = $1 AND u.role = 'participant'
    `

	var c CandidateInfo
	var globalRank int
	var direction sql.NullString

	err := r.db.QueryRow(query, userID).Scan(
		&c.UserID,
		&c.FullName,
		&c.Email,
		&c.Age,
		&c.City,
		&direction,
		&c.TotalPoints,
		&c.EventsCount,
		&globalRank,
		&c.AvgPoints,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get candidate: %w", err)
	}

	if direction.Valid {
		c.Direction = &direction.String
	}
	if globalRank > 0 {
		c.GlobalRank = &globalRank
	}

	return &c, nil
}
