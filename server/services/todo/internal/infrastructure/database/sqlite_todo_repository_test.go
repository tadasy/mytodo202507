package database_test

import (
	"os"
	"testing"

	"github.com/tadasy/mytodo202507/server/services/todo/internal/domain/entity"
	"github.com/tadasy/mytodo202507/server/services/todo/internal/domain/repository"
	"github.com/tadasy/mytodo202507/server/services/todo/internal/infrastructure/database"
)

// ========================================
// 外部振る舞いテスト（Black-box Testing）
// リポジトリインターface経由でのテスト
// 実装の詳細に依存しない
// ========================================

func TestTodoRepository_CreateAndGet(t *testing.T) {
	// Arrange
	dbPath := "test_todos.db"
	defer os.Remove(dbPath)

	// repositoryインターface経由でテスト
	var repo repository.TodoRepository
	sqliteRepo, err := database.NewSQLiteTodoRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer sqliteRepo.Close()
	repo = sqliteRepo

	todo := entity.NewTodo("test-id", "user-123", "Test Todo", "Test Description")

	// Act - Create
	err = repo.Create(todo)
	if err != nil {
		t.Errorf("Failed to create todo: %v", err)
	}

	// Act - Get
	retrievedTodo, err := repo.GetByID(todo.ID, todo.UserID)
	if err != nil {
		t.Errorf("Failed to get todo: %v", err)
	}

	// Assert
	if retrievedTodo.ID != todo.ID {
		t.Errorf("Expected ID %s, got %s", todo.ID, retrievedTodo.ID)
	}
	if retrievedTodo.UserID != todo.UserID {
		t.Errorf("Expected UserID %s, got %s", todo.UserID, retrievedTodo.UserID)
	}
	if retrievedTodo.Title != todo.Title {
		t.Errorf("Expected Title %s, got %s", todo.Title, retrievedTodo.Title)
	}
	if retrievedTodo.Description != todo.Description {
		t.Errorf("Expected Description %s, got %s", todo.Description, retrievedTodo.Description)
	}
	if retrievedTodo.Completed != todo.Completed {
		t.Errorf("Expected Completed %v, got %v", todo.Completed, retrievedTodo.Completed)
	}
}

func TestTodoRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	dbPath := "test_todos_notfound.db"
	defer os.Remove(dbPath)

	var repo repository.TodoRepository
	sqliteRepo, err := database.NewSQLiteTodoRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer sqliteRepo.Close()
	repo = sqliteRepo

	// Act
	_, err = repo.GetByID("nonexistent-id", "user-123")

	// Assert
	if err == nil {
		t.Errorf("Expected error for nonexistent todo")
	}
}

func TestTodoRepository_ListByUserID(t *testing.T) {
	// Arrange
	dbPath := "test_todos_list.db"
	defer os.Remove(dbPath)

	var repo repository.TodoRepository
	sqliteRepo, err := database.NewSQLiteTodoRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer sqliteRepo.Close()
	repo = sqliteRepo

	userID := "user-123"
	todo1 := entity.NewTodo("todo-1", userID, "Todo 1", "Description 1")
	todo2 := entity.NewTodo("todo-2", userID, "Todo 2", "Description 2")
	todo3 := entity.NewTodo("todo-3", "different-user", "Todo 3", "Description 3")

	// 複数のTodoを作成
	repo.Create(todo1)
	repo.Create(todo2)
	repo.Create(todo3) // 異なるユーザーのTodo

	// Act
	todos, err := repo.ListByUserID(userID)

	// Assert
	if err != nil {
		t.Errorf("Failed to list todos: %v", err)
	}
	if len(todos) != 2 {
		t.Errorf("Expected 2 todos for user, got %d", len(todos))
	}

	// 正しいユーザーのTodoのみ取得されていることを確認
	for _, todo := range todos {
		if todo.UserID != userID {
			t.Errorf("Expected all todos to belong to user %s", userID)
		}
	}
}

func TestTodoRepository_ListCompletedByUserID(t *testing.T) {
	// Arrange
	dbPath := "test_todos_completed.db"
	defer os.Remove(dbPath)

	var repo repository.TodoRepository
	sqliteRepo, err := database.NewSQLiteTodoRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer sqliteRepo.Close()
	repo = sqliteRepo

	userID := "user-123"
	todo1 := entity.NewTodo("todo-1", userID, "Todo 1", "Description 1")
	todo2 := entity.NewTodo("todo-2", userID, "Todo 2", "Description 2")

	// 1つのTodoを完了状態にする
	todo1.MarkComplete(true)

	repo.Create(todo1)
	repo.Create(todo2)

	// Act
	completedTodos, err := repo.ListCompletedByUserID(userID)

	// Assert
	if err != nil {
		t.Errorf("Failed to list completed todos: %v", err)
	}
	if len(completedTodos) != 1 {
		t.Errorf("Expected 1 completed todo, got %d", len(completedTodos))
	}
	if len(completedTodos) > 0 && completedTodos[0].ID != todo1.ID {
		t.Errorf("Expected completed todo to be %s", todo1.ID)
	}
}

func TestTodoRepository_Update(t *testing.T) {
	// Arrange
	dbPath := "test_todos_update.db"
	defer os.Remove(dbPath)

	var repo repository.TodoRepository
	sqliteRepo, err := database.NewSQLiteTodoRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer sqliteRepo.Close()
	repo = sqliteRepo

	todo := entity.NewTodo("test-id", "user-123", "Original Title", "Original Description")
	repo.Create(todo)

	// 変更
	todo.Update("Updated Title", "Updated Description")
	todo.MarkComplete(true)

	// Act
	err = repo.Update(todo)
	if err != nil {
		t.Errorf("Failed to update todo: %v", err)
	}

	// 更新されたTodoを取得して確認
	updatedTodo, err := repo.GetByID(todo.ID, todo.UserID)
	if err != nil {
		t.Errorf("Failed to get updated todo: %v", err)
	}

	// Assert
	if updatedTodo.Title != "Updated Title" {
		t.Errorf("Expected title to be updated")
	}
	if updatedTodo.Description != "Updated Description" {
		t.Errorf("Expected description to be updated")
	}
	if !updatedTodo.Completed {
		t.Errorf("Expected todo to be completed")
	}
	if updatedTodo.CompletedAt == nil {
		t.Errorf("Expected CompletedAt to be set")
	}
}

func TestTodoRepository_Delete(t *testing.T) {
	// Arrange
	dbPath := "test_todos_delete.db"
	defer os.Remove(dbPath)

	var repo repository.TodoRepository
	sqliteRepo, err := database.NewSQLiteTodoRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer sqliteRepo.Close()
	repo = sqliteRepo

	todo := entity.NewTodo("test-id", "user-123", "Test Todo", "Test Description")
	repo.Create(todo)

	// Act
	err = repo.Delete(todo.ID, todo.UserID)
	if err != nil {
		t.Errorf("Failed to delete todo: %v", err)
	}

	// 削除されたことを確認
	_, err = repo.GetByID(todo.ID, todo.UserID)
	if err == nil {
		t.Errorf("Expected todo to be deleted")
	}
}

func TestTodoRepository_UserIsolation(t *testing.T) {
	// Arrange
	dbPath := "test_todos_isolation.db"
	defer os.Remove(dbPath)

	var repo repository.TodoRepository
	sqliteRepo, err := database.NewSQLiteTodoRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer sqliteRepo.Close()
	repo = sqliteRepo

	user1ID := "user-1"
	user2ID := "user-2"
	todo1 := entity.NewTodo("todo-1", user1ID, "User1 Todo", "Description")
	todo2 := entity.NewTodo("todo-2", user2ID, "User2 Todo", "Description")

	// Act
	repo.Create(todo1)
	repo.Create(todo2)

	// Assert - ユーザー分離の確認
	user1Todos, err := repo.ListByUserID(user1ID)
	if err != nil {
		t.Errorf("ListByUserID should succeed: %v", err)
	}
	if len(user1Todos) != 1 {
		t.Errorf("User1 should have 1 todo, got %d", len(user1Todos))
	}

	// user2からuser1のTodoにアクセスできないことを確認
	_, err = repo.GetByID(todo1.ID, user2ID)
	if err == nil {
		t.Errorf("User2 should not access User1's todo")
	}
}

func TestTodoRepository_CompletionFlow(t *testing.T) {
	// Arrange
	dbPath := "test_todos_completion.db"
	defer os.Remove(dbPath)

	var repo repository.TodoRepository
	sqliteRepo, err := database.NewSQLiteTodoRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer sqliteRepo.Close()
	repo = sqliteRepo

	userID := "user-123"
	todo := entity.NewTodo("todo-1", userID, "Test Todo", "Description")

	// Act - 未完了状態で作成
	repo.Create(todo)

	// Assert - 初期状態は未完了
	incompleteTodos, err := repo.ListByUserID(userID)
	if err != nil {
		t.Errorf("ListByUserID should succeed: %v", err)
	}
	if len(incompleteTodos) != 1 || incompleteTodos[0].Completed {
		t.Errorf("Todo should be incomplete initially")
	}

	// Act - 完了状態に変更
	todo.MarkComplete(true)
	repo.Update(todo)

	// Assert - 完了済みリストに表示される
	completedTodos, err := repo.ListCompletedByUserID(userID)
	if err != nil {
		t.Errorf("ListCompletedByUserID should succeed: %v", err)
	}
	if len(completedTodos) != 1 {
		t.Errorf("Should have 1 completed todo, got %d", len(completedTodos))
	}
	if !completedTodos[0].Completed {
		t.Errorf("Todo should be marked as completed")
	}
	if completedTodos[0].CompletedAt == nil {
		t.Errorf("CompletedAt should be set")
	}
}
