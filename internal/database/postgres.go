package database

import (
	"database/sql"
	"fmt"
	"log"

	"razum-backend/internal/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init(cfg *config.Config) error {
	// Добавляем параметры для правильной кодировки
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s client_encoding=UTF8",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Проверяем кодировку соединения
	var encoding string
	err = DB.QueryRow("SHOW client_encoding").Scan(&encoding)
	if err != nil {
		return fmt.Errorf("failed to get client encoding: %w", err)
	}
	log.Printf("Client encoding: %s", encoding)

	// Устанавливаем кодировку для всего соединения
	_, err = DB.Exec("SET client_encoding = 'UTF8'")
	if err != nil {
		return fmt.Errorf("failed to set client encoding: %w", err)
	}

	// Проверяем подключение
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connected successfully")
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}

func GetDB() *sql.DB {
	return DB
}
