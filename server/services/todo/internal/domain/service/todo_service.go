package service

import (
	"github.com/google/uuid"
	"github.com/tadasy/mytodo202507/server/services/todo/internal/domain/entity"
	"github.com/tadasy/mytodo202507/server/services/todo/internal/domain/repository"
)

type TodoService struct {
	todoRepo repository.TodoRepository
}

func NewTodoService(todoRepo repository.TodoRepository) *TodoService {
	return &TodoService{
		todoRepo: todoRepo,
	}
}

func (s *TodoService) CreateTodo(userID, title, description string) (*entity.Todo, error) {
	todoID := uuid.New().String()
	todo := entity.NewTodo(todoID, userID, title, description)

	if err := s.todoRepo.Create(todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *TodoService) GetTodo(id, userID string) (*entity.Todo, error) {
	return s.todoRepo.GetByID(id, userID)
}

func (s *TodoService) ListTodos(userID string) ([]*entity.Todo, error) {
	return s.todoRepo.ListByUserID(userID)
}

func (s *TodoService) ListCompletedTodos(userID string) ([]*entity.Todo, error) {
	return s.todoRepo.ListCompletedByUserID(userID)
}

func (s *TodoService) UpdateTodo(id, userID, title, description string) (*entity.Todo, error) {
	todo, err := s.todoRepo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}

	todo.Update(title, description)

	if err := s.todoRepo.Update(todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *TodoService) MarkTodoComplete(id, userID string, completed bool) (*entity.Todo, error) {
	todo, err := s.todoRepo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}

	todo.MarkComplete(completed)

	if err := s.todoRepo.Update(todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *TodoService) DeleteTodo(id, userID string) error {
	return s.todoRepo.Delete(id, userID)
}
