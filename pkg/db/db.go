package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    date    CHAR(8)      NOT NULL,
    title   VARCHAR(255) NOT NULL,
    comment TEXT         NOT NULL,
    repeat  VARCHAR(128) NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

var DB *sql.DB

// Init открывает или создаёт файл БД и при необходимости создаёт схему.
func Init(dbFile string) error {
	install := false
	if _, err := os.Stat(dbFile); err != nil {
		if os.IsNotExist(err) {
			install = true
		} else {
			return fmt.Errorf("stat %s: %w", dbFile, err)
		}
	}

	conn, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	if install {
		if _, err := conn.Exec(schema); err != nil {
			return fmt.Errorf("exec schema: %w", err)
		}
	}

	DB = conn
	return nil
}
