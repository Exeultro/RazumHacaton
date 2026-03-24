package services

import (
	"database/sql"
)

type SearchService struct {
	db *sql.DB
}

func NewSearchService(db *sql.DB) *SearchService {
	return &SearchService{db: db}
}

func (s *SearchService) Search(query string, searchType string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	searchPattern := "%" + query + "%"

	if searchType == "all" || searchType == "events" {
		rows, err := s.db.Query(`
            SELECT id, title, event_date FROM events 
            WHERE title ILIKE $1 OR description ILIKE $1
            LIMIT 10
        `, searchPattern)
		if err == nil {
			events := []map[string]interface{}{}
			for rows.Next() {
				var id, title string
				var eventDate interface{}
				rows.Scan(&id, &title, &eventDate)
				events = append(events, map[string]interface{}{
					"id":         id,
					"title":      title,
					"event_date": eventDate,
				})
			}
			rows.Close()
			result["events"] = events
		}
	}

	if searchType == "all" || searchType == "users" {
		rows, err := s.db.Query(`
            SELECT id, full_name, role FROM users 
            WHERE full_name ILIKE $1 OR email ILIKE $1
            LIMIT 10
        `, searchPattern)
		if err == nil {
			users := []map[string]interface{}{}
			for rows.Next() {
				var id, fullName, role string
				rows.Scan(&id, &fullName, &role)
				users = append(users, map[string]interface{}{
					"id":        id,
					"full_name": fullName,
					"role":      role,
				})
			}
			rows.Close()
			result["users"] = users
		}
	}

	return result, nil
}
