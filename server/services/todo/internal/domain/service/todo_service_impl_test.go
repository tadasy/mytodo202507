package service

import (
	"errors"
	"testing"

	"github.com/tadasy/mytodo202507/server/services/todo/internal/domain/entity"
	"github.com/tadasy/mytodo202507/server/services/todo/internal/domain/repository"
)

// ========================================
// 実装テスト（Implementation Testing）
// サービス層の実装詳細をテスト
// モックを使った依存関係の詳細な検証
// ========================================

// DetailedMockTodoRepository - 詳細なエラーケースをテストするためのモック
type DetailedMockTodoRepository struct {
	todos       map[string]*entity.Todo
	userTodos   map[string][]*entity.Todo
	createError error
	getError    error
	listError   error
	updateError error
	deleteError error
	callLog     []string // 呼び出しログ（実装テストの特徴）
}

func NewDetailedMockTodoRepository() *DetailedMockTodoRepository {
	return &DetailedMockTodoRepository{
		todos:     make(map[string]*entity.Todo),
		userTodos: make(map[string][]*entity.Todo),
		callLog:   make([]string, 0),
	}
}

func (m *DetailedMockTodoRepository) Create(todo *entity.Todo) error {
	m.callLog = append(m.callLog, "Create")
	if m.createError != nil {
		return m.createError
	}
	m.todos[todo.ID] = todo
	m.userTodos[todo.UserID] = append(m.userTodos[todo.UserID], todo)
	return nil
}

func (m *DetailedMockTodoRepository) GetByID(id, userID string) (*entity.Todo, error) {
	m.callLog = append(m.callLog, "GetByID")
	if m.getError != nil {
		return nil, m.getError
	}
	todo, exists := m.todos[id]
	if !exists || todo.UserID != userID {
		return nil, errors.New("todo not found")
	}
	return todo, nil
}

func (m *DetailedMockTodoRepository) ListByUserID(userID string) ([]*entity.Todo, error) {
	m.callLog = append(m.callLog, "ListByUserID")
	if m.listError != nil {
		return nil, m.listError
	}
	return m.userTodos[userID], nil
}

func (m *DetailedMockTodoRepository) ListCompletedByUserID(userID string) ([]*entity.Todo, error) {
	m.callLog = append(m.callLog, "ListCompletedByUserID")
	if m.listError != nil {
		return nil, m.listError
	}
	var completed []*entity.Todo
	for _, todo := range m.userTodos[userID] {
		if todo.Completed {
			completed = append(completed, todo)
		}
	}
	return completed, nil
}

func (m *DetailedMockTodoRepository) Update(todo *entity.Todo) error {
	m.callLog = append(m.callLog, "Update")
	if m.updateError != nil {
		return m.updateError
	}
	m.todos[todo.ID] = todo
	return nil
}

func (m *DetailedMockTodoRepository) Delete(id, userID string) error {
	m.callLog = append(m.callLog, "Delete")
	if m.deleteError != nil {
		return m.deleteError
	}
	todo, exists := m.todos[id]
	if !exists || todo.UserID != userID {
		return errors.New("todo not found")
	}
	delete(m.todos, id)

	// userTodosからも削除
	userTodos := m.userTodos[userID]
	for i, t := range userTodos {
		if t.ID == id {
			m.userTodos[userID] = append(userTodos[:i], userTodos[i+1:]...)
			break
		}
	}
	return nil
}

// Interface compliance check
var _ repository.TodoRepository = (*DetailedMockTodoRepository)(nil)

func TestTodoService_Implementation_CreateTodo_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := NewDetailedMockTodoRepository()
	mockRepo.createError = errors.New("database connection failed")
	todoService := NewTodoService(mockRepo)

	// Act
	_, err := todoService.CreateTodo("user-123", "Test Todo", "Description")

	// Assert
	if err == nil {
		t.Errorf("Expected error from repository")
	}
	if err.Error() != "database connection failed" {
		t.Errorf("Expected 'database connection failed', got '%s'", err.Error())
	}
	// 実装テストの特徴：呼び出しログの確認
	if len(mockRepo.callLog) != 1 || mockRepo.callLog[0] != "Create" {
		t.Errorf("Expected Create to be called, got %v", mockRepo.callLog)
	}
}

func TestTodoService_Implementation_GetTodo_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := NewDetailedMockTodoRepository()
	mockRepo.getError = errors.New("query timeout")
	todoService := NewTodoService(mockRepo)

	// Act
	_, err := todoService.GetTodo("todo-id", "user-123")

	// Assert
	if err == nil {
		t.Errorf("Expected error from repository")
	}
	if err.Error() != "query timeout" {
		t.Errorf("Expected 'query timeout', got '%s'", err.Error())
	}
	// 実装詳細：正確に1回GetByIDが呼ばれることを確認
	if len(mockRepo.callLog) != 1 || mockRepo.callLog[0] != "GetByID" {
		t.Errorf("Expected GetByID to be called once, got %v", mockRepo.callLog)
	}
}

func TestTodoService_Implementation_UpdateTodo_CallSequence(t *testing.T) {
	// Arrange
	mockRepo := NewDetailedMockTodoRepository()
	todoService := NewTodoService(mockRepo)

	// 最初にTodoを作成
	todo, _ := todoService.CreateTodo("user-123", "Original", "Description")
	mockRepo.callLog = nil // ログをリセット

	// Act
	_, err := todoService.UpdateTodo(todo.ID, "user-123", "Updated", "Updated Description")

	// Assert
	if err != nil {
		t.Errorf("UpdateTodo should succeed: %v", err)
	}
	// 実装詳細：正確にGetByID -> Updateの順で呼ばれることを確認
	expectedCalls := []string{"GetByID", "Update"}
	if len(mockRepo.callLog) != 2 {
		t.Errorf("Expected 2 calls, got %d", len(mockRepo.callLog))
	}
	for i, expectedCall := range expectedCalls {
		if mockRepo.callLog[i] != expectedCall {
			t.Errorf("Expected call %d to be %s, got %s", i, expectedCall, mockRepo.callLog[i])
		}
	}
}

func TestTodoService_Implementation_MarkTodoComplete_CallSequence(t *testing.T) {
	// Arrange
	mockRepo := NewDetailedMockTodoRepository()
	todoService := NewTodoService(mockRepo)

	// 最初にTodoを作成
	todo, _ := todoService.CreateTodo("user-123", "Test", "Description")
	mockRepo.callLog = nil // ログをリセット

	// Act
	_, err := todoService.MarkTodoComplete(todo.ID, "user-123", true)

	// Assert
	if err != nil {
		t.Errorf("MarkTodoComplete should succeed: %v", err)
	}
	// 実装詳細：GetByID -> Updateの順で呼ばれることを確認
	expectedCalls := []string{"GetByID", "Update"}
	if len(mockRepo.callLog) != 2 {
		t.Errorf("Expected 2 calls, got %d", len(mockRepo.callLog))
	}
	for i, expectedCall := range expectedCalls {
		if mockRepo.callLog[i] != expectedCall {
			t.Errorf("Expected call %d to be %s, got %s", i, expectedCall, mockRepo.callLog[i])
		}
	}
}

func TestTodoService_Implementation_DeleteTodo_ErrorHandling(t *testing.T) {
	// Arrange
	mockRepo := NewDetailedMockTodoRepository()
	mockRepo.deleteError = errors.New("foreign key constraint")
	todoService := NewTodoService(mockRepo)

	// Act
	err := todoService.DeleteTodo("todo-id", "user-123")

	// Assert
	if err == nil {
		t.Errorf("Expected error from repository")
	}
	if err.Error() != "foreign key constraint" {
		t.Errorf("Expected 'foreign key constraint', got '%s'", err.Error())
	}
	// 実装詳細：Deleteが呼ばれることを確認
	if len(mockRepo.callLog) != 1 || mockRepo.callLog[0] != "Delete" {
		t.Errorf("Expected Delete to be called, got %v", mockRepo.callLog)
	}
}

func TestTodoService_Implementation_ListTodos_ErrorHandling(t *testing.T) {
	// Arrange
	mockRepo := NewDetailedMockTodoRepository()
	mockRepo.listError = errors.New("table scan timeout")
	todoService := NewTodoService(mockRepo)

	// Act
	_, err := todoService.ListTodos("user-123")

	// Assert
	if err == nil {
		t.Errorf("Expected error from repository")
	}
	if err.Error() != "table scan timeout" {
		t.Errorf("Expected 'table scan timeout', got '%s'", err.Error())
	}
	// 実装詳細：ListByUserIDが呼ばれることを確認
	if len(mockRepo.callLog) != 1 || mockRepo.callLog[0] != "ListByUserID" {
		t.Errorf("Expected ListByUserID to be called, got %v", mockRepo.callLog)
	}
}

func TestTodoService_Implementation_MockInternalState(t *testing.T) {
	// Arrange
	mockRepo := NewDetailedMockTodoRepository()
	todoService := NewTodoService(mockRepo)

	// Act
	todo, err := todoService.CreateTodo("user-123", "Test", "Description")
	if err != nil {
		t.Errorf("Create should succeed: %v", err)
	}

	// Assert - 実装テストの特徴：モックの内部状態を直接確認
	if len(mockRepo.todos) != 1 {
		t.Errorf("Expected 1 todo in mock storage, got %d", len(mockRepo.todos))
	}
	if len(mockRepo.userTodos["user-123"]) != 1 {
		t.Errorf("Expected 1 todo for user in mock storage, got %d", len(mockRepo.userTodos["user-123"]))
	}

	// モックの内部データ構造を直接確認
	storedTodo, exists := mockRepo.todos[todo.ID]
	if !exists {
		t.Errorf("Todo should be stored in mock todos map")
	}
	if storedTodo.Title != "Test" {
		t.Errorf("Stored todo title should be 'Test', got '%s'", storedTodo.Title)
	}

	// ユーザーごとのインデックスも確認
	userTodo := mockRepo.userTodos["user-123"][0]
	if userTodo.ID != todo.ID {
		t.Errorf("User todo should match created todo ID")
	}
}
