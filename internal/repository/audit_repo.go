package repository

import (
	"database/sql"
	"fmt"

	"razum-backend/internal/models"

	"github.com/google/uuid"
)

type AuditRepository struct {
	db *sql.DB
}

func NewAuditRepository(db *sql.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

// AddPointsAudit добавляет запись в журнал аудита начисления баллов
func (r *AuditRepository) AddPointsAudit(audit *models.PointsAudit) error {
	query := `
        INSERT INTO points_audit (id, user_id, event_id, points_change, reason, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

	_, err := r.db.Exec(
		query,
		audit.ID,
		audit.UserID,
		audit.EventID,
		audit.PointsChange,
		audit.Reason,
		audit.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to add points audit: %w", err)
	}

	return nil
}

// GetUserPointsHistory возвращает историю начисления баллов пользователя
func (r *AuditRepository) GetUserPointsHistory(userID uuid.UUID) ([]models.PointsAudit, error) {
	query := `
        SELECT id, user_id, event_id, points_change, reason, created_at
        FROM points_audit
        WHERE user_id = $1
        ORDER BY created_at DESC
        LIMIT 100
    `

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user points history: %w", err)
	}
	defer rows.Close()

	audits := []models.PointsAudit{}
	for rows.Next() {
		var a models.PointsAudit
		err := rows.Scan(
			&a.ID,
			&a.UserID,
			&a.EventID,
			&a.PointsChange,
			&a.Reason,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit: %w", err)
		}
		audits = append(audits, a)
	}

	return audits, nil
}
