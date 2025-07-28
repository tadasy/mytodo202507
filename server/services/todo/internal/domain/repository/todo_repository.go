package repository

import (
	"github.com/tadasy/todo-app/server/services/todo/internal/domain/entity"
)

type TodoRepository interface {
	Create(todo *entity.Todo) error
	GetByID(id, userID string) (*entity.Todo, error)
	ListByUserID(userID string) ([]*entity.Todo, error)
	ListCompletedByUserID(userID string) ([]*entity.Todo, error)
	Update(todo *entity.Todo) error
	Delete(id, userID string) error
}
