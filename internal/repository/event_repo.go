package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"razum-backend/internal/models"

	"github.com/google/uuid"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

// Create создает новое мероприятие
func (r *EventRepository) Create(event *models.Event) error {
	prizesJSON, err := json.Marshal(event.Prizes)
	if err != nil {
		return fmt.Errorf("failed to marshal prizes: %w", err)
	}

	query := `
        INSERT INTO events (
            id, organizer_id, title, description, event_date, registration_deadline,
            location, format, direction, difficulty_coefficient, points_for_participation,
            prizes, status, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
        RETURNING created_at, updated_at
    `

	err = r.db.QueryRow(
		query,
		event.ID,
		event.OrganizerID,
		event.Title,
		event.Description,
		event.EventDate,
		event.RegistrationDeadline,
		event.Location,
		event.Format,
		event.Direction,
		event.DifficultyCoefficient,
		event.PointsForParticipation,
		prizesJSON,
		event.Status,
		event.CreatedAt,
		event.UpdatedAt,
	).Scan(&event.CreatedAt, &event.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	// Обновляем счетчик мероприятий у организатора
	_, err = r.db.Exec(`
        UPDATE organizers 
        SET events_count = events_count + 1 
        WHERE user_id = $1
    `, event.OrganizerID)

	if err != nil {
		return fmt.Errorf("failed to update organizer events count: %w", err)
	}

	return nil
}

// FindByID находит мероприятие по ID
func (r *EventRepository) FindByID(id uuid.UUID) (*models.Event, error) {
	query := `
        SELECT id, organizer_id, title, description, event_date, registration_deadline,
               location, format, direction, difficulty_coefficient, points_for_participation,
               prizes, status, created_at, updated_at
        FROM events
        WHERE id = $1
    `

	event := &models.Event{}
	var prizesJSON []byte

	err := r.db.QueryRow(query, id).Scan(
		&event.ID,
		&event.OrganizerID,
		&event.Title,
		&event.Description,
		&event.EventDate,
		&event.RegistrationDeadline,
		&event.Location,
		&event.Format,
		&event.Direction,
		&event.DifficultyCoefficient,
		&event.PointsForParticipation,
		&prizesJSON,
		&event.Status,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find event by id: %w", err)
	}

	if err := json.Unmarshal(prizesJSON, &event.Prizes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal prizes: %w", err)
	}

	return event, nil
}

// List возвращает список мероприятий с фильтрацией
func (r *EventRepository) List(status, direction, format string, page, limit int) ([]models.Event, int, error) {
	offset := (page - 1) * limit

	log.Printf("List events: status=%s, direction=%s, format=%s, page=%d, limit=%d",
		status, direction, format, page, limit)

	// Базовый запрос
	baseQuery := `FROM events e WHERE 1=1`

	whereArgs := []interface{}{}
	argCounter := 1

	if status != "" {
		baseQuery += fmt.Sprintf(" AND e.status = $%d", argCounter)
		whereArgs = append(whereArgs, status)
		argCounter++
	}

	if direction != "" {
		baseQuery += fmt.Sprintf(" AND e.direction = $%d", argCounter)
		whereArgs = append(whereArgs, direction)
		argCounter++
	}

	if format != "" {
		baseQuery += fmt.Sprintf(" AND e.format = $%d", argCounter)
		whereArgs = append(whereArgs, format)
		argCounter++
	}

	// Подсчет количества
	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int
	err := r.db.QueryRow(countQuery, whereArgs...).Scan(&total)
	if err != nil {
		log.Printf("ERROR in count query: %v", err)
		return nil, 0, fmt.Errorf("failed to count events: %w", err)
	}

	log.Printf("Total events: %d", total)

	if total == 0 {
		return []models.Event{}, 0, nil
	}

	// Запрос данных (ТОЛЬКО поля из models.Event!)
	dataQuery := `SELECT 
        e.id, 
        e.organizer_id, 
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
        e.updated_at
        ` + baseQuery +
		" ORDER BY e.event_date DESC LIMIT $" + fmt.Sprint(argCounter) + " OFFSET $" + fmt.Sprint(argCounter+1)

	dataArgs := append(whereArgs, limit, offset)

	rows, err := r.db.Query(dataQuery, dataArgs...)
	if err != nil {
		log.Printf("ERROR in data query: %v", err)
		return nil, 0, fmt.Errorf("failed to list events: %w", err)
	}
	defer rows.Close()

	events := []models.Event{}
	for rows.Next() {
		var event models.Event
		var prizesJSON []byte

		err := rows.Scan(
			&event.ID,
			&event.OrganizerID,
			&event.Title,
			&event.Description,
			&event.EventDate,
			&event.RegistrationDeadline,
			&event.Location,
			&event.Format,
			&event.Direction,
			&event.DifficultyCoefficient,
			&event.PointsForParticipation,
			&prizesJSON,
			&event.Status,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			log.Printf("ERROR in scan: %v", err)
			return nil, 0, fmt.Errorf("failed to scan event: %w", err)
		}
		log.Printf("Loaded event: ID=%s, Title=%s", event.ID, event.Title)
		if len(prizesJSON) > 0 {
			if err := json.Unmarshal(prizesJSON, &event.Prizes); err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal prizes: %w", err)
			}
		}

		events = append(events, event)
	}

	log.Printf("Total events loaded: %d", len(events))
	return events, total, nil
}

// Update обновляет мероприятие
func (r *EventRepository) Update(event *models.Event) error {
	prizesJSON, err := json.Marshal(event.Prizes)
	if err != nil {
		return fmt.Errorf("failed to marshal prizes: %w", err)
	}

	query := `
        UPDATE events
        SET title = $1, description = $2, event_date = $3, registration_deadline = $4,
            location = $5, format = $6, direction = $7, difficulty_coefficient = $8,
            points_for_participation = $9, prizes = $10, status = $11, updated_at = CURRENT_TIMESTAMP
        WHERE id = $12
        RETURNING updated_at
    `

	err = r.db.QueryRow(
		query,
		event.Title,
		event.Description,
		event.EventDate,
		event.RegistrationDeadline,
		event.Location,
		event.Format,
		event.Direction,
		event.DifficultyCoefficient,
		event.PointsForParticipation,
		prizesJSON,
		event.Status,
		event.ID,
	).Scan(&event.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	return nil
}

// Delete удаляет мероприятие
func (r *EventRepository) Delete(id uuid.UUID) error {
	// Сначала получаем organizer_id
	var organizerID uuid.UUID
	err := r.db.QueryRow("SELECT organizer_id FROM events WHERE id = $1", id).Scan(&organizerID)
	if err != nil {
		return fmt.Errorf("failed to get event organizer: %w", err)
	}

	// Удаляем мероприятие
	_, err = r.db.Exec("DELETE FROM events WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	// Обновляем счетчик мероприятий у организатора
	_, err = r.db.Exec(`
        UPDATE organizers 
        SET events_count = events_count - 1 
        WHERE user_id = $1 AND events_count > 0
    `, organizerID)

	if err != nil {
		return fmt.Errorf("failed to update organizer events count: %w", err)
	}

	return nil
}

// GetEventParticipantsCount возвращает количество участников мероприятия по статусу
func (r *EventRepository) GetEventParticipantsCount(eventID uuid.UUID, status string) (int, error) {
	var count int
	query := `
        SELECT COUNT(*)
        FROM event_participations
        WHERE event_id = $1
    `

	if status != "" {
		query += " AND status = $2"
		err := r.db.QueryRow(query, eventID, status).Scan(&count)
		return count, err
	}

	err := r.db.QueryRow(query, eventID).Scan(&count)
	return count, err
}

// GetEventParticipantsCounts возвращает все счетчики участников по статусам
func (r *EventRepository) GetEventParticipantsCounts(eventID uuid.UUID) (map[string]int, error) {
	query := `
        SELECT status, COUNT(*)
        FROM event_participations
        WHERE event_id = $1
        GROUP BY status
    `

	rows, err := r.db.Query(query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants counts: %w", err)
	}
	defer rows.Close()

	counts := map[string]int{
		"registered": 0,
		"confirmed":  0,
		"attended":   0,
		"cancelled":  0,
		"total":      0,
	}

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("failed to scan count: %w", err)
		}
		counts[status] = count
		counts["total"] += count
	}

	// Вычитаем cancelled из total
	counts["total"] -= counts["cancelled"]

	return counts, nil
}

// ListWithFilters получает список мероприятий с фильтрацией по датам
func (r *EventRepository) ListWithFilters(status, direction, format, dateFrom, dateTo, search string, page, limit int) ([]models.Event, int, error) {
	offset := (page - 1) * limit

	query := `
        SELECT e.id, e.organizer_id, e.title, e.description, e.event_date, e.registration_deadline,
               e.location, e.format, e.direction, e.difficulty_coefficient, e.points_for_participation,
               e.prizes, e.status, e.created_at, e.updated_at	
        FROM events e
        WHERE 1=1
    `
	countQuery := "SELECT COUNT(*) FROM events e WHERE 1=1"
	args := []interface{}{}
	argIdx := 1

	if status != "" {
		query += fmt.Sprintf(" AND e.status = $%d", argIdx)
		countQuery += fmt.Sprintf(" AND e.status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	if direction != "" {
		query += fmt.Sprintf(" AND e.direction = $%d", argIdx)
		countQuery += fmt.Sprintf(" AND e.direction = $%d", argIdx)
		args = append(args, direction)
		argIdx++
	}

	if format != "" {
		query += fmt.Sprintf(" AND e.format = $%d", argIdx)
		countQuery += fmt.Sprintf(" AND e.format = $%d", argIdx)
		args = append(args, format)
		argIdx++
	}

	// Фильтр по дате начала (только дата)
	if dateFrom != "" {
		dateOnly := dateFrom[:10] // YYYY-MM-DD
		query += fmt.Sprintf(" AND e.event_date::date >= $%d", argIdx)
		countQuery += fmt.Sprintf(" AND e.event_date::date >= $%d", argIdx)
		args = append(args, dateOnly)
		argIdx++
	}

	// Фильтр по дате окончания (только дата)
	if dateTo != "" {
		dateOnly := dateTo[:10] // YYYY-MM-DD
		query += fmt.Sprintf(" AND e.event_date::date <= $%d", argIdx)
		countQuery += fmt.Sprintf(" AND e.event_date::date <= $%d", argIdx)
		args = append(args, dateOnly)
		argIdx++
	}

	// Поиск по названию
	if search != "" {
		query += fmt.Sprintf(" AND e.title ILIKE $%d", argIdx)
		countQuery += fmt.Sprintf(" AND e.title ILIKE $%d", argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	// Считаем общее количество
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count events: %w", err)
	}

	if total == 0 {
		return []models.Event{}, 0, nil
	}

	query += " ORDER BY e.event_date ASC LIMIT $" + strconv.Itoa(argIdx) + " OFFSET $" + strconv.Itoa(argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	events := []models.Event{}
	for rows.Next() {
		var e models.Event
		var prizesJSON []byte
		err := rows.Scan(
			&e.ID, &e.OrganizerID, &e.Title, &e.Description,
			&e.EventDate, &e.RegistrationDeadline, &e.Location,
			&e.Format, &e.Direction, &e.DifficultyCoefficient,
			&e.PointsForParticipation, &prizesJSON, &e.Status,
			&e.CreatedAt, &e.UpdatedAt,
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
