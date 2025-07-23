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

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	if err != nil {
		return fmt.Errorf("БД не найдена %v/n", err)
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("невозможно открыть БД: %v", err)
	}

	_, err = db.Exec(schema)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы: %v", err)
	}

	return nil
}

func GetDB() *sql.DB {
	return db
}
