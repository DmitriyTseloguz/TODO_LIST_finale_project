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

func (db *DB) GetTask(taskID int) (model.Task, error) {
	query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE id = ?
	`

	var task model.Task

	err := db.conn.QueryRowx(query, taskID).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		return model.Task{}, fmt.Errorf("failed to get task: %w", err)
	}

	return task, nil
}

func (db *DB) GetAllTasks() ([]model.Task, error) {
	query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
	`

	var tasks []model.Task = make([]model.Task, 0)

	err := db.conn.Select(&tasks, query)

	if err != nil {
		return tasks, fmt.Errorf("failed to get tasks: %w", err)
	}

	return tasks, nil
}

func (db *DB) CreateTask(task *model.Task) (int, error) {
	query := `
        INSERT INTO scheduler (date, title, comment, repeat)
        VALUES (:date, :title, :comment, :repeat)
    `

	result, err := db.conn.NamedExec(query, map[string]any{
		"date":    task.Date,
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

func (db *DB) UpdateTask(taskID int, task *model.Task) error {
	_, err := db.GetTask(taskID)
	if err != nil {
		return fmt.Errorf("task not found")
	}

	query := `
		UPDATE scheduler
		SET date = :date, title = :title, comment = :comment, repeat = :repeat
		WHERE id = :id
	`

	data := map[string]interface{}{
		"id":      taskID,
		"date":    task.Date,
		"title":   task.Title,
		"comment": task.Comment,
		"repeat":  task.Repeat,
	}

	result, err := db.conn.NamedExec(query, data)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func (db *DB) DeleteTask(taskID int) error {
	query := `
		DELETE FROM scheduler
		WHERE id = ?
	`

	_, err := db.conn.Exec(query, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}
