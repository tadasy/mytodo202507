package entity_test

import (
	"testing"
	"time"

	"github.com/tadasy/mytodo202507/server/services/todo/internal/domain/entity"
)

// ========================================
// 外部振る舞いテスト（Black-box Testing）
// エンティティの振る舞いをテスト
// ビジネスルールの検証
// ========================================

func TestTodo_NewTodo(t *testing.T) {
	// Arrange
	id := "test-id"
	userID := "user-123"
	title := "Test Todo"
	description := "Test Description"

	// Act
	todo := entity.NewTodo(id, userID, title, description)

	// Assert
	if todo.ID != id {
		t.Errorf("Expected ID %s, got %s", id, todo.ID)
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
		t.Errorf("Expected Completed to be false")
	}
	if todo.CompletedAt != nil {
		t.Errorf("Expected CompletedAt to be nil")
	}
	if todo.CreatedAt.IsZero() {
		t.Errorf("Expected CreatedAt to be set")
	}
	if todo.UpdatedAt.IsZero() {
		t.Errorf("Expected UpdatedAt to be set")
	}
}

func TestTodo_MarkComplete(t *testing.T) {
	// Arrange
	todo := entity.NewTodo("test-id", "user-123", "Test Todo", "Test Description")
	originalUpdatedAt := todo.UpdatedAt

	// 時間差を作るために少し待つ
	time.Sleep(10 * time.Millisecond)

	// Act - 完了状態にする
	todo.MarkComplete(true)

	// Assert
	if !todo.Completed {
		t.Errorf("Expected todo to be completed")
	}
	if todo.CompletedAt == nil {
		t.Errorf("Expected CompletedAt to be set")
	}
	if !todo.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf("Expected UpdatedAt to be updated")
	}

	// Act - 未完了状態に戻す
	todo.MarkComplete(false)

	// Assert
	if todo.Completed {
		t.Errorf("Expected todo to be incomplete")
	}
	if todo.CompletedAt != nil {
		t.Errorf("Expected CompletedAt to be nil")
	}
}

func TestTodo_Update(t *testing.T) {
	// Arrange
	todo := entity.NewTodo("test-id", "user-123", "Test Todo", "Test Description")
	originalUpdatedAt := todo.UpdatedAt

	// 時間差を作るために少し待つ
	time.Sleep(10 * time.Millisecond)

	newTitle := "Updated Title"
	newDescription := "Updated Description"

	// Act
	todo.Update(newTitle, newDescription)

	// Assert
	if todo.Title != newTitle {
		t.Errorf("Expected title to be updated to %s, got %s", newTitle, todo.Title)
	}
	if todo.Description != newDescription {
		t.Errorf("Expected description to be updated to %s, got %s", newDescription, todo.Description)
	}
	if !todo.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf("Expected UpdatedAt to be updated")
	}
}

func TestTodo_BusinessRules(t *testing.T) {
	// Arrange
	todo := entity.NewTodo("test-id", "user-123", "Test Todo", "Test Description")

	// Act & Assert - 初期状態の確認
	if todo.Completed {
		t.Errorf("New todo should be incomplete by default")
	}

	// Act & Assert - 完了→未完了の遷移
	todo.MarkComplete(true)
	if !todo.Completed || todo.CompletedAt == nil {
		t.Errorf("Todo should be completed with timestamp")
	}

	todo.MarkComplete(false)
	if todo.Completed || todo.CompletedAt != nil {
		t.Errorf("Todo should be incomplete without timestamp")
	}

	// Act & Assert - 更新操作でCompletedは変更されない
	originalCompleted := todo.Completed
	todo.Update("New Title", "New Description")
	if todo.Completed != originalCompleted {
		t.Errorf("Update should not change completion status")
	}
}
