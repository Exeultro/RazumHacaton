package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"razum-backend/internal/models"
	"razum-backend/internal/repository"

	"github.com/google/uuid"
)

type AdminService struct {
	adminRepo *repository.AdminRepository
	userRepo  *repository.UserRepository
	db        *sql.DB
}

func NewAdminService(
	adminRepo *repository.AdminRepository,
	userRepo *repository.UserRepository,
	db *sql.DB,
) *AdminService {
	return &AdminService{
		adminRepo: adminRepo,
		userRepo:  userRepo,
		db:        db,
	}
}

// GetPendingOrganizers возвращает список организаторов на модерации
func (s *AdminService) GetPendingOrganizers() ([]repository.PendingOrganizer, error) {
	return s.adminRepo.GetPendingOrganizers()
}

// ApproveOrganizer одобряет организатора
func (s *AdminService) ApproveOrganizer(userID uuid.UUID) error {
	return s.adminRepo.ApproveOrganizer(userID)
}

// RejectOrganizer отклоняет организатора
func (s *AdminService) RejectOrganizer(userID uuid.UUID) error {
	return s.adminRepo.RejectOrganizer(userID)
}

// GetDifficultySettings получает настройки весов
func (s *AdminService) GetDifficultySettings() (*repository.DifficultySettings, error) {
	return s.adminRepo.GetDifficultySettings()
}

// UpdateDifficultySettings обновляет настройки весов
func (s *AdminService) UpdateDifficultySettings(settings *repository.DifficultySettings) error {
	return s.adminRepo.UpdateDifficultySettings(settings)
}

// GetStats возвращает статистику по платформе
func (s *AdminService) GetStats() (map[string]interface{}, error) {
	var totalUsers, totalOrganizers, totalEvents, totalParticipations int
	eventsByDirection := make(map[string]int)

	// Получаем общее количество пользователей
	err := s.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalUsers)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Получаем количество организаторов
	err = s.db.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'organizer'").Scan(&totalOrganizers)
	if err != nil {
		return nil, fmt.Errorf("failed to count organizers: %w", err)
	}

	// Получаем количество мероприятий
	err = s.db.QueryRow("SELECT COUNT(*) FROM events").Scan(&totalEvents)
	if err != nil {
		return nil, fmt.Errorf("failed to count events: %w", err)
	}

	// Получаем количество участий
	err = s.db.QueryRow("SELECT COUNT(*) FROM event_participations WHERE status = 'attended'").Scan(&totalParticipations)
	if err != nil {
		return nil, fmt.Errorf("failed to count participations: %w", err)
	}

	// Получаем мероприятия по направлениям
	rows, err := s.db.Query("SELECT direction, COUNT(*) FROM events GROUP BY direction")
	if err != nil {
		return nil, fmt.Errorf("failed to get events by direction: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var direction string
		var count int
		if err := rows.Scan(&direction, &count); err != nil {
			continue
		}
		eventsByDirection[direction] = count
	}

	return map[string]interface{}{
		"total_users":          totalUsers,
		"total_organizers":     totalOrganizers,
		"total_events":         totalEvents,
		"total_participations": totalParticipations,
		"events_by_direction":  eventsByDirection,
	}, nil
}

func (s *AdminService) GetAllUsers(page, limit int) ([]models.User, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Считаем общее количество
	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Получаем пользователей с пагинацией
	query := `
        SELECT id, email, full_name, role, city, age, direction, avatar_url, created_at, updated_at
        FROM users
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `

	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
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
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, u)
	}

	return users, total, nil
}

// DeleteUser удаляет пользователя
func (s *AdminService) DeleteUser(userID uuid.UUID) error {
	return s.userRepo.Delete(userID)
}

// ChangeUserRole меняет роль пользователя
func (s *AdminService) ChangeUserRole(userID uuid.UUID, role string) error {
	return s.userRepo.UpdateRole(userID, models.UserRole(role))
}

// EventWithOrganizer структура для мероприятия с информацией об организаторе
type EventWithOrganizer struct {
	ID                     uuid.UUID      `json:"id"`
	Title                  string         `json:"title"`
	Description            string         `json:"description"`
	EventDate              time.Time      `json:"event_date"`
	RegistrationDeadline   time.Time      `json:"registration_deadline"`
	Location               *string        `json:"location,omitempty"`
	Format                 string         `json:"format"`
	Direction              string         `json:"direction"`
	DifficultyCoefficient  float64        `json:"difficulty_coefficient"`
	PointsForParticipation int            `json:"points_for_participation"`
	Prizes                 []models.Prize `json:"prizes"`
	Status                 string         `json:"status"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
	OrganizerID            uuid.UUID      `json:"organizer_id"`
	OrganizerName          string         `json:"organizer_name"`
	OrganizerEmail         string         `json:"organizer_email"`
}

// GetAllEvents возвращает все мероприятия с пагинацией
func (s *AdminService) GetAllEvents(page, limit int) ([]EventWithOrganizer, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM events").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count events: %w", err)
	}

	query := `
        SELECT 
            e.id,
            e.title,
            e.description,
            e.event_date,
            e.registration_deadline,
            e.location,
            e.format,
            e.direction,
            e.difficulty_coefficient,
            e.points_for_participation,
            e.prizes,
            e.status,
            e.created_at,
            e.updated_at,
            e.organizer_id,
            u.full_name as organizer_name,
            u.email as organizer_email
        FROM events e
        JOIN users u ON e.organizer_id = u.id
        ORDER BY e.created_at DESC
        LIMIT $1 OFFSET $2
    `

	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get events: %w", err)
	}
	defer rows.Close()

	events := []EventWithOrganizer{}
	for rows.Next() {
		var e EventWithOrganizer
		var prizesJSON []byte

		err := rows.Scan(
			&e.ID,
			&e.Title,
			&e.Description,
			&e.EventDate,
			&e.RegistrationDeadline,
			&e.Location,
			&e.Format,
			&e.Direction,
			&e.DifficultyCoefficient,
			&e.PointsForParticipation,
			&prizesJSON,
			&e.Status,
			&e.CreatedAt,
			&e.UpdatedAt,
			&e.OrganizerID,
			&e.OrganizerName,
			&e.OrganizerEmail,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan event: %w", err)
		}

		if len(prizesJSON) > 0 {
			json.Unmarshal(prizesJSON, &e.Prizes)
		}

		events = append(events, e)
	}

	return events, total, nil
}

// DeleteEventByAdmin удаляет мероприятие по ID
func (s *AdminService) DeleteEventByAdmin(eventID uuid.UUID) error {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM events WHERE id = $1)", eventID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check event existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("event not found")
	}

	_, err = s.db.Exec("DELETE FROM events WHERE id = $1", eventID)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return nil
}
