package service_test

import (
	"errors"
	"testing"

	"github.com/tadasy/mytodo202507/server/services/todo/internal/domain/entity"
	"github.com/tadasy/mytodo202507/server/services/todo/internal/domain/repository"
	"github.com/tadasy/mytodo202507/server/services/todo/internal/domain/service"
)

// ========================================
// 外部振る舞いテスト（Black-box Testing）
// サービス層の振る舞いをテスト
// ビジネスロジックの検証
// ========================================

// SimpleMockRepository - テスト用の簡単なモックリポジトリ
type SimpleMockRepository struct {
	todos map[string]*entity.Todo
}

func NewSimpleMockRepository() *SimpleMockRepository {
	return &SimpleMockRepository{
		todos: make(map[string]*entity.Todo),
	}
}

func (m *SimpleMockRepository) Create(todo *entity.Todo) error {
	m.todos[todo.ID] = todo
	return nil
}

func (m *SimpleMockRepository) GetByID(id, userID string) (*entity.Todo, error) {
	todo, exists := m.todos[id]
	if !exists || todo.UserID != userID {
		return nil, errors.New("todo not found")
	}
	return todo, nil
}

func (m *SimpleMockRepository) ListByUserID(userID string) ([]*entity.Todo, error) {
	var result []*entity.Todo
	for _, todo := range m.todos {
		if todo.UserID == userID {
			result = append(result, todo)
		}
	}
	return result, nil
}

func (m *SimpleMockRepository) ListCompletedByUserID(userID string) ([]*entity.Todo, error) {
	var result []*entity.Todo
	for _, todo := range m.todos {
		if todo.UserID == userID && todo.Completed {
			result = append(result, todo)
		}
	}
	return result, nil
}

func (m *SimpleMockRepository) Update(todo *entity.Todo) error {
	m.todos[todo.ID] = todo
	return nil
}

func (m *SimpleMockRepository) Delete(id, userID string) error {
	todo, exists := m.todos[id]
	if !exists || todo.UserID != userID {
		return errors.New("todo not found")
	}
	delete(m.todos, id)
	return nil
}

// Interface compliance check
var _ repository.TodoRepository = (*SimpleMockRepository)(nil)

func TestTodoService_CreateTodo(t *testing.T) {
	// Arrange
	repo := NewSimpleMockRepository()
	todoService := service.NewTodoService(repo)
	
	userID := "user-123"
	title := "Test Todo"
	description := "Test Description"

	// Act
	todo, err := todoService.CreateTodo(userID, title, description)

	// Assert
	if err != nil {
		t.Errorf("CreateTodo should succeed: %v", err)
	}
	if todo.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, todo.UserID)
	}
	if todo.Title != title {
		t.Errorf("Expected Title %s, got %s", title, todo.Title)
	}
	if todo.Description != description {
		t.Errorf("Expected Description %s, got %s", description, todo.Description)
	}
	if todo.Completed {
		t.Errorf("New todo should be incomplete")
	}
}

func TestTodoService_GetTodo(t *testing.T) {
	// Arrange
	repo := NewSimpleMockRepository()
	todoService := service.NewTodoService(repo)
	
	userID := "user-123"
	todo, _ := todoService.CreateTodo(userID, "Test Todo", "Description")

	// Act
	retrievedTodo, err := todoService.GetTodo(todo.ID, userID)

	// Assert
	if err != nil {
		t.Errorf("GetTodo should succeed: %v", err)
	}
	if retrievedTodo.ID != todo.ID {
		t.Errorf("Expected ID %s, got %s", todo.ID, retrievedTodo.ID)
	}
}

func TestTodoService_ListTodos(t *testing.T) {
	// Arrange
	repo := NewSimpleMockRepository()
	todoService := service.NewTodoService(repo)
	
	userID := "user-123"
	todoService.CreateTodo(userID, "Todo 1", "Description 1")
	todoService.CreateTodo(userID, "Todo 2", "Description 2")
	todoService.CreateTodo("other-user", "Other Todo", "Other Description")

	// Act
	todos, err := todoService.ListTodos(userID)

	// Assert
	if err != nil {
		t.Errorf("ListTodos should succeed: %v", err)
	}
	if len(todos) != 2 {
		t.Errorf("Expected 2 todos, got %d", len(todos))
	}
	for _, todo := range todos {
		if todo.UserID != userID {
			t.Errorf("All todos should belong to user %s", userID)
		}
	}
}

func TestTodoService_UpdateTodo(t *testing.T) {
	// Arrange
	repo := NewSimpleMockRepository()
	todoService := service.NewTodoService(repo)
	
	userID := "user-123"
	todo, _ := todoService.CreateTodo(userID, "Original Title", "Original Description")

	newTitle := "Updated Title"
	newDescription := "Updated Description"

	// Act
	updatedTodo, err := todoService.UpdateTodo(todo.ID, userID, newTitle, newDescription)

	// Assert
	if err != nil {
		t.Errorf("UpdateTodo should succeed: %v", err)
	}
	if updatedTodo.Title != newTitle {
		t.Errorf("Expected title %s, got %s", newTitle, updatedTodo.Title)
	}
	if updatedTodo.Description != newDescription {
		t.Errorf("Expected description %s, got %s", newDescription, updatedTodo.Description)
	}
}

func TestTodoService_MarkTodoComplete(t *testing.T) {
	// Arrange
	repo := NewSimpleMockRepository()
	todoService := service.NewTodoService(repo)
	
	userID := "user-123"
	todo, _ := todoService.CreateTodo(userID, "Test Todo", "Description")

	// Act
	completedTodo, err := todoService.MarkTodoComplete(todo.ID, userID, true)

	// Assert
	if err != nil {
		t.Errorf("MarkTodoComplete should succeed: %v", err)
	}
	if !completedTodo.Completed {
		t.Errorf("Todo should be marked as completed")
	}
	if completedTodo.CompletedAt == nil {
		t.Errorf("CompletedAt should be set")
	}
}

func TestTodoService_DeleteTodo(t *testing.T) {
	// Arrange
	repo := NewSimpleMockRepository()
	todoService := service.NewTodoService(repo)
	
	userID := "user-123"
	todo, _ := todoService.CreateTodo(userID, "Test Todo", "Description")

	// Act
	err := todoService.DeleteTodo(todo.ID, userID)

	// Assert
	if err != nil {
		t.Errorf("DeleteTodo should succeed: %v", err)
	}

	// 削除後の確認
	_, getErr := todoService.GetTodo(todo.ID, userID)
	if getErr == nil {
		t.Errorf("Todo should be deleted")
	}
}

func TestTodoService_ListCompletedTodos(t *testing.T) {
	// Arrange
	repo := NewSimpleMockRepository()
	todoService := service.NewTodoService(repo)
	
	userID := "user-123"
	todo1, _ := todoService.CreateTodo(userID, "Todo 1", "Description 1")
	_, _ = todoService.CreateTodo(userID, "Todo 2", "Description 2")
	
	// 1つを完了状態にする
	todoService.MarkTodoComplete(todo1.ID, userID, true)

	// Act
	completedTodos, err := todoService.ListCompletedTodos(userID)

	// Assert
	if err != nil {
		t.Errorf("ListCompletedTodos should succeed: %v", err)
	}
	if len(completedTodos) != 1 {
		t.Errorf("Expected 1 completed todo, got %d", len(completedTodos))
	}
	if completedTodos[0].ID != todo1.ID {
		t.Errorf("Expected completed todo to be %s", todo1.ID)
	}
}

func TestTodoService_UserIsolation(t *testing.T) {
	// Arrange
	repo := NewSimpleMockRepository()
	todoService := service.NewTodoService(repo)
	
	user1ID := "user-1"
	user2ID := "user-2"
	
	todo1, _ := todoService.CreateTodo(user1ID, "User1 Todo", "Description")
	todoService.CreateTodo(user2ID, "User2 Todo", "Description")

	// Act & Assert - user2はuser1のTodoにアクセスできない
	_, err := todoService.GetTodo(todo1.ID, user2ID)
	if err == nil {
		t.Errorf("User2 should not access User1's todo")
	}

	// Act & Assert - user2はuser1のTodoを更新できない
	_, updateErr := todoService.UpdateTodo(todo1.ID, user2ID, "Hacked", "Hacked")
	if updateErr == nil {
		t.Errorf("User2 should not update User1's todo")
	}
}
