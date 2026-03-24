package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

// PendingOrganizer информация об организаторе на модерации
type PendingOrganizer struct {
	UserID    string `json:"user_id"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	City      string `json:"city"`
	CreatedAt string `json:"created_at"`
}

// GetPendingOrganizers возвращает список организаторов на модерации
func (r *AdminRepository) GetPendingOrganizers() ([]PendingOrganizer, error) {
	query := `
        SELECT u.id, u.full_name, u.email, u.city, u.created_at
        FROM users u
        WHERE u.role = 'organizer' 
        AND NOT EXISTS (SELECT 1 FROM organizers o WHERE o.user_id = u.id)
    `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending organizers: %w", err)
	}
	defer rows.Close()

	organizers := []PendingOrganizer{}
	for rows.Next() {
		var o PendingOrganizer
		err := rows.Scan(&o.UserID, &o.FullName, &o.Email, &o.City, &o.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan organizer: %w", err)
		}
		organizers = append(organizers, o)
	}

	return organizers, nil
}

// ApproveOrganizer одобряет организатора
func (r *AdminRepository) ApproveOrganizer(userID uuid.UUID) error {
	query := `
        INSERT INTO organizers (user_id, rating, events_count, common_prizes)
        VALUES ($1, 0, 0, '{}')
    `

	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to approve organizer: %w", err)
	}

	return nil
}

// RejectOrganizer отклоняет организатора (удаляет пользователя)
func (r *AdminRepository) RejectOrganizer(userID uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1 AND role = 'organizer'`

	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to reject organizer: %w", err)
	}

	return nil
}

// DifficultySettings настройки весов баллов
type DifficultySettings struct {
	ITCoefficient     float64 `json:"it_coefficient"`
	SocialCoefficient float64 `json:"social_coefficient"`
	MediaCoefficient  float64 `json:"media_coefficient"`
}

// GetDifficultySettings получает текущие настройки
func (r *AdminRepository) GetDifficultySettings() (*DifficultySettings, error) {
	// Временно храним настройки в отдельной таблице или возвращаем значения по умолчанию
	// Для простоты вернем значения по умолчанию
	return &DifficultySettings{
		ITCoefficient:     2.0,
		SocialCoefficient: 1.5,
		MediaCoefficient:  1.2,
	}, nil
}

// UpdateDifficultySettings обновляет настройки весов
func (r *AdminRepository) UpdateDifficultySettings(settings *DifficultySettings) error {
	// Здесь можно сохранять в отдельную таблицу settings
	// Для простоты пока пропустим
	return nil
}
