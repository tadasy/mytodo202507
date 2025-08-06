package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tadasy/mytodo202507/server/services/user/internal/domain/entity"
	"github.com/tadasy/mytodo202507/server/services/user/internal/domain/repository"
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
		user.CreatedAt.Format("2006-01-02 15:04:05"), 
		user.UpdatedAt.Format("2006-01-02 15:04:05"))
	return err
}

func (r *SQLiteUserRepository) GetByID(id string) (*entity.User, error) {
	query := `
	SELECT id, email, password_hash, created_at, updated_at
	FROM users WHERE id = ?`

	row := r.db.QueryRow(query, id)
	user, err := r.scanUser(row)
	if err == sql.ErrNoRows {
		return nil, repository.ErrUserNotFound
	}
	return user, err
}

func (r *SQLiteUserRepository) GetByEmail(email string) (*entity.User, error) {
	query := `
	SELECT id, email, password_hash, created_at, updated_at
	FROM users WHERE email = ?`

	row := r.db.QueryRow(query, email)
	user, err := r.scanUser(row)
	if err == sql.ErrNoRows {
		return nil, repository.ErrUserNotFound
	}
	return user, err
}

func (r *SQLiteUserRepository) Update(user *entity.User) error {
	query := `
	UPDATE users SET email = ?, password_hash = ?, updated_at = ?
	WHERE id = ?`

	result, err := r.db.Exec(query, user.Email, user.PasswordHash,
		user.UpdatedAt.Format("2006-01-02 15:04:05"), user.ID)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return repository.ErrUserNotFound
	}
	
	return nil
}

func (r *SQLiteUserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return repository.ErrUserNotFound
	}
	
	return nil
}

func (r *SQLiteUserRepository) scanUser(row *sql.Row) (*entity.User, error) {
	var user entity.User
	var createdAt, updatedAt string

	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash,
		&createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	// 複数の時刻フォーマットを試行（SQLiteドライバーがフォーマットを変換する場合に対応）
	timeFormats := []string{
		"2006-01-02 15:04:05",           // 保存時の形式
		"2006-01-02T15:04:05Z",          // SQLiteドライバーがISO形式に変換
		time.RFC3339,                    // 標準ISO形式
	}
	
	for _, format := range timeFormats {
		if user.CreatedAt, err = time.Parse(format, createdAt); err == nil {
			// SQLiteは常にUTCで保存されるため、UTCに統一
			user.CreatedAt = user.CreatedAt.UTC()
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse CreatedAt '%s': %w", createdAt, err)
	}
	
	for _, format := range timeFormats {
		if user.UpdatedAt, err = time.Parse(format, updatedAt); err == nil {
			// SQLiteは常にUTCで保存されるため、UTCに統一
			user.UpdatedAt = user.UpdatedAt.UTC()
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse UpdatedAt '%s': %w", updatedAt, err)
	}

	return &user, nil
}

func (r *SQLiteUserRepository) Close() error {
	return r.db.Close()
}
