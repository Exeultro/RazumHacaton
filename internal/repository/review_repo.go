package repository

import (
	"database/sql"
	"fmt"

	"razum-backend/internal/models"

	"github.com/google/uuid"
)

type ReviewRepository struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

// CreateReview создает отзыв
func (r *ReviewRepository) CreateReview(review *models.OrganizerReview) error {
	query := `
        INSERT INTO organizer_reviews (id, organizer_id, participant_id, event_id, rating, comment, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	_, err := r.db.Exec(query,
		review.ID,
		review.OrganizerID,
		review.ParticipantID,
		review.EventID,
		review.Rating,
		review.Comment,
		review.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}

	return nil
}

// GetReviewsByOrganizer получает отзывы об организаторе
func (r *ReviewRepository) GetReviewsByOrganizer(organizerID uuid.UUID, limit, offset int) ([]models.ReviewWithDetails, int, error) {
	// Считаем общее количество
	var total int
	err := r.db.QueryRow(`
        SELECT COUNT(*) FROM organizer_reviews WHERE organizer_id = $1
    `, organizerID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count reviews: %w", err)
	}

	// Получаем отзывы с деталями
	query := `
        SELECT 
            r.id,
            r.organizer_id,
            u_o.full_name as organizer_name,
            r.participant_id,
            u_p.full_name as participant_name,
            r.event_id,
            e.title as event_title,
            r.rating,
            r.comment,
            r.created_at
        FROM organizer_reviews r
        JOIN users u_o ON r.organizer_id = u_o.id
        JOIN users u_p ON r.participant_id = u_p.id
        LEFT JOIN events e ON r.event_id = e.id
        WHERE r.organizer_id = $1
        ORDER BY r.created_at DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := r.db.Query(query, organizerID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get reviews: %w", err)
	}
	defer rows.Close()

	reviews := []models.ReviewWithDetails{}
	for rows.Next() {
		var rev models.ReviewWithDetails
		var eventTitle sql.NullString

		err := rows.Scan(
			&rev.ID,
			&rev.OrganizerID,
			&rev.OrganizerName,
			&rev.ParticipantID,
			&rev.ParticipantName,
			&rev.EventID,
			&eventTitle,
			&rev.Rating,
			&rev.Comment,
			&rev.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan review: %w", err)
		}

		if eventTitle.Valid {
			rev.EventTitle = &eventTitle.String
		}

		reviews = append(reviews, rev)
	}

	return reviews, total, nil
}

// GetOrganizerRatingStats получает статистику рейтинга организатора
func (r *ReviewRepository) GetOrganizerRatingStats(organizerID uuid.UUID) (*models.OrganizerRatingStats, error) {
	query := `
        SELECT 
            COALESCE(AVG(rating), 0) as average_rating,
            COUNT(*) as total_reviews,
            COUNT(CASE WHEN rating = 5 THEN 1 END) as rating_5,
            COUNT(CASE WHEN rating = 4 THEN 1 END) as rating_4,
            COUNT(CASE WHEN rating = 3 THEN 1 END) as rating_3,
            COUNT(CASE WHEN rating = 2 THEN 1 END) as rating_2,
            COUNT(CASE WHEN rating = 1 THEN 1 END) as rating_1
        FROM organizer_reviews
        WHERE organizer_id = $1
    `

	stats := &models.OrganizerRatingStats{
		OrganizerID: organizerID,
	}

	err := r.db.QueryRow(query, organizerID).Scan(
		&stats.AverageRating,
		&stats.TotalReviews,
		&stats.Rating5Count,
		&stats.Rating4Count,
		&stats.Rating3Count,
		&stats.Rating2Count,
		&stats.Rating1Count,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get rating stats: %w", err)
	}

	return stats, nil
}

// CheckUserCanReview проверяет, может ли пользователь оставить отзыв
func (r *ReviewRepository) CheckUserCanReview(participantID, organizerID uuid.UUID) (bool, error) {
	query := `
        SELECT 
            EXISTS(
                SELECT 1 
                FROM event_participations ep
                JOIN events e ON ep.event_id = e.id
                WHERE ep.user_id = $1 
                AND e.organizer_id = $2 
                AND ep.status = 'attended'
            ) AND NOT EXISTS(
                SELECT 1 
                FROM organizer_reviews 
                WHERE participant_id = $1 AND organizer_id = $2
            )
    `

	var canReview bool
	err := r.db.QueryRow(query, participantID, organizerID).Scan(&canReview)
	if err != nil {
		return false, fmt.Errorf("failed to check review permission: %w", err)
	}

	return canReview, nil
}

// UpdateOrganizerRating обновляет рейтинг организатора
func (r *ReviewRepository) UpdateOrganizerRating(organizerID uuid.UUID, rating float64) error {
	query := `
        UPDATE organizers 
        SET rating = $1 
        WHERE user_id = $2
    `

	_, err := r.db.Exec(query, rating, organizerID)
	if err != nil {
		return fmt.Errorf("failed to update organizer rating: %w", err)
	}

	return nil
}
