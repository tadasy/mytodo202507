package entity

import (
	"time"
)

type Todo struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// NewTodo creates a new todo item
func NewTodo(id, userID, title, description string) *Todo {
	now := time.Now()
	return &Todo{
		ID:          id,
		UserID:      userID,
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Update updates the todo's title and description
func (t *Todo) Update(title, description string) {
	if title != "" {
		t.Title = title
	}
	if description != "" {
		t.Description = description
	}
	t.UpdatedAt = time.Now()
}

// MarkComplete marks the todo as completed or uncompleted
func (t *Todo) MarkComplete(completed bool) {
	t.Completed = completed
	t.UpdatedAt = time.Now()

	if completed {
		now := time.Now()
		t.CompletedAt = &now
	} else {
		t.CompletedAt = nil
	}
}
