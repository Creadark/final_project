package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL,
    comment TEXT,
    repeat VARCHAR(128)
);
CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
`

// Init инициализирует подключение к базе данных
func Init(dbFile string) error {
	if _, err := os.Stat(dbFile); err != nil {
		return fmt.Errorf("база данных не найдена: %w", err)
	}

	var err error
	db, err = sql.Open("sqlite3", dbFile) // Исправлено с "sqlite" на "sqlite3"
	if err != nil {
		return fmt.Errorf("ошибка открытия базы данных: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("ошибка подключения к базе: %w", err)
	}

	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("ошибка создания схемы: %w", err)
	}

	return nil
}

// Close закрывает соединение с базой данных
func Close() error {
	if db != nil {
		if err := db.Close(); err != nil {
			return fmt.Errorf("ошибка закрытия соединения: %w", err)
		}
		db = nil
	}
	return nil
}

// GetDB возвращает текущее подключение к БД
func GetDB() *sql.DB {
	return db
}
