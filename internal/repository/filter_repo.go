package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"razum-backend/internal/models"

	"github.com/google/uuid"
)

type FilterRepository struct {
	db *sql.DB
}

func NewFilterRepository(db *sql.DB) *FilterRepository {
	return &FilterRepository{db: db}
}

// CreateFilter создает фильтр
func (r *FilterRepository) CreateFilter(filter *models.ObserverFilter) error {
	filtersJSON, err := json.Marshal(filter.Filters)
	if err != nil {
		return fmt.Errorf("failed to marshal filters: %w", err)
	}

	query := `
        INSERT INTO observer_filters (id, user_id, filter_name, filters, created_at)
        VALUES ($1, $2, $3, $4, $5)
    `

	_, err = r.db.Exec(query, filter.ID, filter.UserID, filter.FilterName, filtersJSON, filter.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create filter: %w", err)
	}

	return nil
}

// GetUserFilters возвращает все фильтры пользователя
func (r *FilterRepository) GetUserFilters(userID uuid.UUID) ([]models.ObserverFilter, error) {
	query := `
        SELECT id, user_id, filter_name, filters, created_at
        FROM observer_filters
        WHERE user_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user filters: %w", err)
	}
	defer rows.Close()

	filters := []models.ObserverFilter{}
	for rows.Next() {
		var f models.ObserverFilter
		var filtersJSON []byte

		err := rows.Scan(&f.ID, &f.UserID, &f.FilterName, &filtersJSON, &f.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan filter: %w", err)
		}

		if err := json.Unmarshal(filtersJSON, &f.Filters); err != nil {
			return nil, fmt.Errorf("failed to unmarshal filters: %w", err)
		}

		filters = append(filters, f)
	}

	return filters, nil
}

// GetFilterByID получает фильтр по ID
func (r *FilterRepository) GetFilterByID(filterID uuid.UUID) (*models.ObserverFilter, error) {
	query := `
        SELECT id, user_id, filter_name, filters, created_at
        FROM observer_filters
        WHERE id = $1
    `

	var f models.ObserverFilter
	var filtersJSON []byte

	err := r.db.QueryRow(query, filterID).Scan(&f.ID, &f.UserID, &f.FilterName, &filtersJSON, &f.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get filter: %w", err)
	}

	if err := json.Unmarshal(filtersJSON, &f.Filters); err != nil {
		return nil, fmt.Errorf("failed to unmarshal filters: %w", err)
	}

	return &f, nil
}

// UpdateFilter обновляет фильтр
func (r *FilterRepository) UpdateFilter(filter *models.ObserverFilter) error {
	filtersJSON, err := json.Marshal(filter.Filters)
	if err != nil {
		return fmt.Errorf("failed to marshal filters: %w", err)
	}

	query := `
        UPDATE observer_filters
        SET filter_name = $1, filters = $2
        WHERE id = $3 AND user_id = $4
    `

	result, err := r.db.Exec(query, filter.FilterName, filtersJSON, filter.ID, filter.UserID)
	if err != nil {
		return fmt.Errorf("failed to update filter: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("filter not found or not owned by user")
	}

	return nil
}

// DeleteFilter удаляет фильтр
func (r *FilterRepository) DeleteFilter(filterID, userID uuid.UUID) error {
	query := `DELETE FROM observer_filters WHERE id = $1 AND user_id = $2`

	result, err := r.db.Exec(query, filterID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete filter: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("filter not found or not owned by user")
	}

	return nil
}
