package database

import (
	"fmt"

	"todo-list/lib/model"

	"github.com/jmoiron/sqlx"
)

var dbInstance *DB

type DB struct {
	conn *sqlx.DB
}

func NewDB(dsn string) (*DB, error) {
	conn, err := sqlx.Connect("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	dbInstance = &DB{conn: conn}

	return dbInstance, nil
}

func GetDB() *DB {
	return dbInstance
}

func (db *DB) GetConnection() *sqlx.DB {
	return db.conn
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) Initialize() error {
	_, err := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS scheduler(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT,
			title TEXT,
			comment TEXT,
			repeat VARCHAR(128)
		);
		CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler(date);
	`)

	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	return nil
}

func (db *DB) CreateTask(task *model.Task) (int, error) {
	query := `
        INSERT INTO scheduler (date, title, comment, repeat)
        VALUES (:date, :title, :comment, :repeat)
    `

	result, err := db.conn.NamedExec(query, map[string]any{
		"date":    task.GetTime().Format("20060102"),
		"title":   task.Title,
		"comment": task.Comment,
		"repeat":  task.Repeat,
	})

	if err != nil {
		return 0, fmt.Errorf("failed to create task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return int(id), nil
}
