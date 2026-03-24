package repository

import (
	"database/sql"
	"fmt"

	"razum-backend/internal/models"

	"github.com/google/uuid"
)

type RatingRepository struct {
	db *sql.DB
}

func NewRatingRepository(db *sql.DB) *RatingRepository {
	return &RatingRepository{db: db}
}

// UpdateRatingCache обновляет кэш рейтинга для всех пользователей
func (r *RatingRepository) UpdateRatingCache() error {
	// Сначала очищаем кэш
	_, err := r.db.Exec("TRUNCATE rating_cache")
	if err != nil {
		return fmt.Errorf("failed to truncate rating cache: %w", err)
	}

	// Заполняем кэш актуальными данными
	query := `
        INSERT INTO rating_cache (user_id, total_points, events_count, top_direction)
        SELECT 
            u.id,
            COALESCE(SUM(ep.points_earned), 0) as total_points,
            COUNT(CASE WHEN ep.status = 'attended' THEN 1 END) as events_count,
            (
                SELECT direction
                FROM (
                    SELECT 
                        e.direction,
                        SUM(ep2.points_earned) as points_by_direction
                    FROM event_participations ep2
                    JOIN events e ON ep2.event_id = e.id
                    WHERE ep2.user_id = u.id AND ep2.status = 'attended'
                    GROUP BY e.direction
                    ORDER BY points_by_direction DESC
                    LIMIT 1
                ) sub
            ) as top_direction
        FROM users u
        LEFT JOIN event_participations ep ON u.id = ep.user_id AND ep.status = 'attended'
        GROUP BY u.id
    `

	_, err = r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to update rating cache: %w", err)
	}

	// Обновляем глобальные ранги
	_, err = r.db.Exec(`
        UPDATE rating_cache rc
        SET global_rank = sub.rank
        FROM (
            SELECT user_id, ROW_NUMBER() OVER (ORDER BY total_points DESC) as rank
            FROM rating_cache
        ) sub
        WHERE rc.user_id = sub.user_id
    `)
	if err != nil {
		return fmt.Errorf("failed to update global ranks: %w", err)
	}

	return nil
}

// GetGlobalRating возвращает глобальную таблицу лидеров
func (r *RatingRepository) GetGlobalRating(limit, offset int) ([]models.RatingCache, int, error) {
	// Получаем общее количество
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM rating_cache WHERE total_points > 0").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count ratings: %w", err)
	}

	// Получаем данные с пагинацией
	query := `
        SELECT rc.user_id, rc.total_points, rc.events_count, rc.top_direction, rc.global_rank, 
               u.full_name, u.avatar_url, u.direction
        FROM rating_cache rc
        JOIN users u ON rc.user_id = u.id
        WHERE rc.total_points > 0
        ORDER BY rc.global_rank
        LIMIT $1 OFFSET $2
    `

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get global rating: %w", err)
	}
	defer rows.Close()

	ratings := []models.RatingCache{}
	for rows.Next() {
		var rating models.RatingCache
		var fullName string
		var avatarURL sql.NullString
		var direction sql.NullString

		err := rows.Scan(
			&rating.UserID,
			&rating.TotalPoints,
			&rating.EventsCount,
			&rating.TopDirection,
			&rating.GlobalRank,
			&fullName,
			&avatarURL,
			&direction,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan rating: %w", err)
		}

		rating.FullName = fullName
		if avatarURL.Valid {
			rating.AvatarURL = &avatarURL.String
		}
		if direction.Valid {
			rating.Direction = &direction.String
		}

		ratings = append(ratings, rating)
	}

	return ratings, total, nil
}

// GetRatingByDirection возвращает рейтинг по направлению
func (r *RatingRepository) GetRatingByDirection(direction string, limit, offset int) ([]models.RatingCache, int, error) {
	// Получаем общее количество
	var total int
	err := r.db.QueryRow(`
        SELECT COUNT(*) 
        FROM rating_cache rc
        JOIN users u ON rc.user_id = u.id
        WHERE rc.total_points > 0 AND u.direction = $1
    `, direction).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count direction ratings: %w", err)
	}

	// Получаем данные с пагинацией
	query := `
        SELECT rc.user_id, rc.total_points, rc.events_count, rc.top_direction, 
               u.full_name, u.avatar_url
        FROM rating_cache rc
        JOIN users u ON rc.user_id = u.id
        WHERE rc.total_points > 0 AND u.direction = $1
        ORDER BY rc.total_points DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := r.db.Query(query, direction, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get direction rating: %w", err)
	}
	defer rows.Close()

	ratings := []models.RatingCache{}
	rank := offset + 1
	for rows.Next() {
		var rating models.RatingCache
		var fullName string
		var avatarURL sql.NullString

		err := rows.Scan(
			&rating.UserID,
			&rating.TotalPoints,
			&rating.EventsCount,
			&rating.TopDirection,
			&fullName,
			&avatarURL,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan rating: %w", err)
		}

		rating.FullName = fullName
		if avatarURL.Valid {
			rating.AvatarURL = &avatarURL.String
		}
		rating.GlobalRank = &rank
		rank++

		ratings = append(ratings, rating)
	}

	return ratings, total, nil
}

// GetUserRating получает рейтинг конкретного пользователя
func (r *RatingRepository) GetUserRating(userID uuid.UUID) (*models.RatingCache, error) {
	query := `
        SELECT rc.user_id, rc.total_points, rc.events_count, rc.top_direction, rc.global_rank, 
               u.full_name, u.avatar_url, u.direction
        FROM rating_cache rc
        JOIN users u ON rc.user_id = u.id
        WHERE rc.user_id = $1
    `

	var rating models.RatingCache
	var fullName string
	var avatarURL sql.NullString
	var direction sql.NullString

	err := r.db.QueryRow(query, userID).Scan(
		&rating.UserID,
		&rating.TotalPoints,
		&rating.EventsCount,
		&rating.TopDirection,
		&rating.GlobalRank,
		&fullName,
		&avatarURL,
		&direction,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user rating: %w", err)
	}

	rating.FullName = fullName
	if avatarURL.Valid {
		rating.AvatarURL = &avatarURL.String
	}
	if direction.Valid {
		rating.Direction = &direction.String
	}

	return &rating, nil
}
