package database

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tadasy/mytodo202507/server/services/user/internal/domain/entity"
)

type SQLiteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(dbPath string) (*SQLiteUserRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	repo := &SQLiteUserRepository{db: db}
	if err := repo.createTable(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *SQLiteUserRepository) createTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	)`
	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteUserRepository) Create(user *entity.User) error {
	query := `
	INSERT INTO users (id, email, password_hash, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, user.ID, user.Email, user.PasswordHash,
		user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *SQLiteUserRepository) GetByID(id string) (*entity.User, error) {
	query := `
	SELECT id, email, password_hash, created_at, updated_at
	FROM users WHERE id = ?`

	row := r.db.QueryRow(query, id)
	return r.scanUser(row)
}

func (r *SQLiteUserRepository) GetByEmail(email string) (*entity.User, error) {
	query := `
	SELECT id, email, password_hash, created_at, updated_at
	FROM users WHERE email = ?`

	row := r.db.QueryRow(query, email)
	return r.scanUser(row)
}

func (r *SQLiteUserRepository) Update(user *entity.User) error {
	query := `
	UPDATE users SET email = ?, password_hash = ?, updated_at = ?
	WHERE id = ?`

	_, err := r.db.Exec(query, user.Email, user.PasswordHash,
		user.UpdatedAt, user.ID)
	return err
}

func (r *SQLiteUserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *SQLiteUserRepository) scanUser(row *sql.Row) (*entity.User, error) {
	var user entity.User
	var createdAt, updatedAt string

	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash,
		&createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	user.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	return &user, nil
}

func (r *SQLiteUserRepository) Close() error {
	return r.db.Close()
}
