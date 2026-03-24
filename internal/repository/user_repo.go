package repository

import (
	"database/sql"
	"fmt"

	"razum-backend/internal/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create создает нового пользователя
func (r *UserRepository) Create(user *models.User) error {
	query := `
        INSERT INTO users (id, email, password, full_name, role, city, age, direction, avatar_url)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING created_at, updated_at
    `

	err := r.db.QueryRow(
		query,
		user.ID,
		user.Email,
		user.Password,
		user.FullName,
		user.Role,
		user.City,
		user.Age,
		user.Direction,
		user.AvatarURL,
	).Scan(&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Если пользователь - организатор, создаем запись в organizers
	if user.Role == models.RoleOrganizer {
		queryOrganizer := `
            INSERT INTO organizers (user_id, rating, events_count, common_prizes)
            VALUES ($1, $2, $3, $4)
        `
		_, err = r.db.Exec(queryOrganizer, user.ID, 0.0, 0, "{}")
		if err != nil {
			return fmt.Errorf("failed to create organizer record: %w", err)
		}
	}

	return nil
}

// FindByEmail находит пользователя по email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	query := `
        SELECT id, email, password, full_name, role, city, age, direction, avatar_url, created_at, updated_at
        FROM users
        WHERE email = $1
    `

	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FullName,
		&user.Role,
		&user.City,
		&user.Age,
		&user.Direction,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return user, nil
}

// FindByID находит пользователя по ID
func (r *UserRepository) FindByID(userID uuid.UUID) (*models.User, error) {
	query := `
        SELECT id, email, full_name, role, city, age, direction, avatar_url, created_at, updated_at
        FROM users
        WHERE id = $1
    `

	var user models.User
	err := r.db.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.FullName,
		&user.Role,
		&user.City,
		&user.Age,
		&user.Direction,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &user, nil
}

// Update обновляет данные пользователя
func (r *UserRepository) Update(user *models.User) error {
	query := `
        UPDATE users
        SET full_name = $1, city = $2, age = $3, direction = $4, avatar_url = $5, updated_at = CURRENT_TIMESTAMP
        WHERE id = $6
        RETURNING updated_at
    `

	err := r.db.QueryRow(
		query,
		user.FullName,
		user.City,
		user.Age,
		user.Direction,
		user.AvatarURL,
		user.ID,
	).Scan(&user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// GetOrganizerStats получает статистику организатора
func (r *UserRepository) GetOrganizerStats(userID uuid.UUID) (*models.Organizer, error) {
	query := `
        SELECT user_id, rating, events_count, common_prizes
        FROM organizers
        WHERE user_id = $1
    `

	stats := &models.Organizer{}

	err := r.db.QueryRow(query, userID).Scan(
		&stats.UserID,
		&stats.Rating,
		&stats.EventsCount,
		pq.Array(&stats.CommonPrizes),
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get organizer stats: %w", err)
	}

	return stats, nil
}

// GetAll возвращает всех пользователей
func (r *UserRepository) GetAll() ([]models.User, error) {
	query := `
        SELECT id, email, full_name, role, city, age, direction, avatar_url, created_at, updated_at
        FROM users
        ORDER BY created_at DESC
    `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.FullName,
			&u.Role,
			&u.City,
			&u.Age,
			&u.Direction,
			&u.AvatarURL,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, u)
	}

	return users, nil
}

// Delete удаляет пользователя
func (r *UserRepository) Delete(userID uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// UpdateRole обновляет роль пользователя
func (r *UserRepository) UpdateRole(userID uuid.UUID, role models.UserRole) error {
	query := `UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(query, role, userID)
	if err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}
	return nil
}

// CheckOrganizerExists проверяет, есть ли запись в organizers
func (r *UserRepository) CheckOrganizerExists(userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM organizers WHERE user_id = $1)", userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check organizer exists: %w", err)
	}
	return exists, nil
}

// CreateOrganizer создает запись в organizers
func (r *UserRepository) CreateOrganizer(userID uuid.UUID) error {
	query := `
        INSERT INTO organizers (user_id, rating, events_count, common_prizes)
        VALUES ($1, 0, 0, '{}')
    `
	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to create organizer: %w", err)
	}
	return nil
}
