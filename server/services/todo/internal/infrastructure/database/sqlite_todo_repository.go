package database

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tadasy/mytodo202507/server/services/todo/internal/domain/entity"
)

type SQLiteTodoRepository struct {
	db *sql.DB
}

func NewSQLiteTodoRepository(dbPath string) (*SQLiteTodoRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	repo := &SQLiteTodoRepository{db: db}
	if err := repo.createTable(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *SQLiteTodoRepository) createTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS todos (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		title TEXT NOT NULL,
		description TEXT,
		completed BOOLEAN NOT NULL DEFAULT FALSE,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		completed_at DATETIME
	)`
	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteTodoRepository) Create(todo *entity.Todo) error {
	query := `
	INSERT INTO todos (id, user_id, title, description, completed, created_at, updated_at, completed_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	var completedAt interface{}
	if todo.CompletedAt != nil {
		completedAt = todo.CompletedAt
	}

	_, err := r.db.Exec(query, todo.ID, todo.UserID, todo.Title, todo.Description,
		todo.Completed, todo.CreatedAt, todo.UpdatedAt, completedAt)
	return err
}

func (r *SQLiteTodoRepository) GetByID(id, userID string) (*entity.Todo, error) {
	query := `
	SELECT id, user_id, title, description, completed, created_at, updated_at, completed_at
	FROM todos WHERE id = ? AND user_id = ?`

	row := r.db.QueryRow(query, id, userID)
	return r.scanTodo(row)
}

func (r *SQLiteTodoRepository) ListByUserID(userID string) ([]*entity.Todo, error) {
	query := `
	SELECT id, user_id, title, description, completed, created_at, updated_at, completed_at
	FROM todos WHERE user_id = ? ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*entity.Todo
	for rows.Next() {
		todo, err := r.scanTodoFromRows(rows)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func (r *SQLiteTodoRepository) ListCompletedByUserID(userID string) ([]*entity.Todo, error) {
	query := `
	SELECT id, user_id, title, description, completed, created_at, updated_at, completed_at
	FROM todos WHERE user_id = ? AND completed = TRUE ORDER BY completed_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*entity.Todo
	for rows.Next() {
		todo, err := r.scanTodoFromRows(rows)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func (r *SQLiteTodoRepository) Update(todo *entity.Todo) error {
	query := `
	UPDATE todos SET title = ?, description = ?, completed = ?, updated_at = ?, completed_at = ?
	WHERE id = ? AND user_id = ?`

	var completedAt interface{}
	if todo.CompletedAt != nil {
		completedAt = todo.CompletedAt
	}

	_, err := r.db.Exec(query, todo.Title, todo.Description, todo.Completed,
		todo.UpdatedAt, completedAt, todo.ID, todo.UserID)
	return err
}

func (r *SQLiteTodoRepository) Delete(id, userID string) error {
	query := `DELETE FROM todos WHERE id = ? AND user_id = ?`
	_, err := r.db.Exec(query, id, userID)
	return err
}

func (r *SQLiteTodoRepository) scanTodo(row *sql.Row) (*entity.Todo, error) {
	var todo entity.Todo
	var createdAt, updatedAt string
	var completedAt sql.NullString

	err := row.Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description,
		&todo.Completed, &createdAt, &updatedAt, &completedAt)
	if err != nil {
		return nil, err
	}

	todo.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	todo.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	if completedAt.Valid {
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", completedAt.String)
		todo.CompletedAt = &parsedTime
	}

	return &todo, nil
}

func (r *SQLiteTodoRepository) scanTodoFromRows(rows *sql.Rows) (*entity.Todo, error) {
	var todo entity.Todo
	var createdAt, updatedAt string
	var completedAt sql.NullString

	err := rows.Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description,
		&todo.Completed, &createdAt, &updatedAt, &completedAt)
	if err != nil {
		return nil, err
	}

	todo.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	todo.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	if completedAt.Valid {
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", completedAt.String)
		todo.CompletedAt = &parsedTime
	}

	return &todo, nil
}

func (r *SQLiteTodoRepository) Close() error {
	return r.db.Close()
}
