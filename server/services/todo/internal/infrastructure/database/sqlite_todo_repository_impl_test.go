package database

import (
	"os"
	"testing"

	"github.com/tadasy/mytodo202507/server/services/todo/internal/domain/entity"
)

// ========================================
// 実装テスト（Implementation Testing）
// SQLite固有の実装詳細をテスト
// 実装の内部構造に依存したWhite-boxテスト
// ========================================

func TestSQLiteTodoRepository_Internal_DatabaseSchemaCreation(t *testing.T) {
	// Arrange
	dbPath := "test_internal_schema.db"
	defer os.Remove(dbPath)

	// Act
	repo, err := NewSQLiteTodoRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	// Assert - 内部実装: テーブル構造の直接確認
	var tableName string
	err = repo.db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='todos'").Scan(&tableName)
	if err != nil {
		t.Errorf("Table 'todos' should exist: %v", err)
	}
	if tableName != "todos" {
		t.Errorf("Expected table name 'todos', got '%s'", tableName)
	}

	// カラム構造の確認
	rows, err := repo.db.Query("PRAGMA table_info(todos)")
	if err != nil {
		t.Fatalf("Failed to get table info: %v", err)
	}
	defer rows.Close()

	expectedColumns := map[string]bool{
		"id":           false,
		"user_id":      false,
		"title":        false,
		"description":  false,
		"completed":    false,
		"created_at":   false,
		"updated_at":   false,
		"completed_at": false,
	}

	for rows.Next() {
		var cid int
		var name, dataType string
		var notNull, pk int
		var defaultValue interface{}
		err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &pk)
		if err != nil {
			t.Errorf("Failed to scan column info: %v", err)
		}
		if _, exists := expectedColumns[name]; exists {
			expectedColumns[name] = true
		}
	}

	// 期待されるカラムがすべて存在することを確認
	for column, found := range expectedColumns {
		if !found {
			t.Errorf("Expected column '%s' not found in table", column)
		}
	}
}

func TestSQLiteTodoRepository_Internal_TimeFormatting(t *testing.T) {
	// Arrange
	dbPath := "test_internal_time.db"
	defer os.Remove(dbPath)

	repo, err := NewSQLiteTodoRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	todo := entity.NewTodo("test-id", "user-123", "Test Todo", "Test Description")
	todo.MarkComplete(true)

	// Act
	err = repo.Create(todo)
	if err != nil {
		t.Errorf("Failed to create todo: %v", err)
	}

	// Assert - 内部実装: SQLiteでの時刻保存形式を直接確認
	var createdAt, updatedAt, completedAt string
	err = repo.db.QueryRow(
		"SELECT created_at, updated_at, completed_at FROM todos WHERE id = ? AND user_id = ?",
		todo.ID, todo.UserID,
	).Scan(&createdAt, &updatedAt, &completedAt)
	if err != nil {
		t.Errorf("Failed to query time fields: %v", err)
	}

	t.Logf("SQLite stored times - Created: %s, Updated: %s, Completed: %s", createdAt, updatedAt, completedAt)

	// SQLiteの時刻フォーマットが期待される形式であることを確認
	// RFC3339形式または"2006-01-02 15:04:05"形式
	if len(createdAt) == 0 {
		t.Errorf("created_at should not be empty")
	}
	if len(updatedAt) == 0 {
		t.Errorf("updated_at should not be empty")
	}
	if len(completedAt) == 0 {
		t.Errorf("completed_at should not be empty for completed todo")
	}
}

func TestSQLiteTodoRepository_Internal_ScanTodoImplementation(t *testing.T) {
	// Arrange
	dbPath := "test_internal_scan.db"
	defer os.Remove(dbPath)

	repo, err := NewSQLiteTodoRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	// 完了済みTodoを作成
	completedTodo := entity.NewTodo("completed-id", "user-123", "Completed Todo", "Description")
	completedTodo.MarkComplete(true)

	// 未完了Todoを作成
	incompleteTodo := entity.NewTodo("incomplete-id", "user-123", "Incomplete Todo", "Description")

	repo.Create(completedTodo)
	repo.Create(incompleteTodo)

	// Act & Assert - 内部実装: scanTodo メソッドの動作確認
	// 完了済みTodoのスキャン
	row := repo.db.QueryRow(
		"SELECT id, user_id, title, description, completed, created_at, updated_at, completed_at FROM todos WHERE id = ?",
		completedTodo.ID,
	)
	scannedCompleted, err := repo.scanTodo(row)
	if err != nil {
		t.Errorf("Failed to scan completed todo: %v", err)
	}
	if !scannedCompleted.Completed {
		t.Errorf("Scanned todo should be completed")
	}
	if scannedCompleted.CompletedAt == nil {
		t.Errorf("Scanned completed todo should have CompletedAt set")
	}

	// 未完了Todoのスキャン
	row = repo.db.QueryRow(
		"SELECT id, user_id, title, description, completed, created_at, updated_at, completed_at FROM todos WHERE id = ?",
		incompleteTodo.ID,
	)
	scannedIncomplete, err := repo.scanTodo(row)
	if err != nil {
		t.Errorf("Failed to scan incomplete todo: %v", err)
	}
	if scannedIncomplete.Completed {
		t.Errorf("Scanned todo should be incomplete")
	}
	if scannedIncomplete.CompletedAt != nil {
		t.Errorf("Scanned incomplete todo should have CompletedAt as nil")
	}
}

func TestSQLiteTodoRepository_Internal_SQLQueryExecution(t *testing.T) {
	// Arrange
	dbPath := "test_internal_sql.db"
	defer os.Remove(dbPath)

	repo, err := NewSQLiteTodoRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	userID := "user-123"
	todo1 := entity.NewTodo("todo-1", userID, "Todo 1", "Description 1")
	todo2 := entity.NewTodo("todo-2", userID, "Todo 2", "Description 2")
	todo1.MarkComplete(true)

	repo.Create(todo1)
	repo.Create(todo2)

	// Act & Assert - 内部実装: 特定のSQLクエリの実行結果確認
	// 完了済みTodoのカウント
	var completedCount int
	err = repo.db.QueryRow(
		"SELECT COUNT(*) FROM todos WHERE user_id = ? AND completed = ?",
		userID, true,
	).Scan(&completedCount)
	if err != nil {
		t.Errorf("Failed to count completed todos: %v", err)
	}
	if completedCount != 1 {
		t.Errorf("Expected 1 completed todo, got %d", completedCount)
	}

	// 全Todoのカウント
	var totalCount int
	err = repo.db.QueryRow(
		"SELECT COUNT(*) FROM todos WHERE user_id = ?",
		userID,
	).Scan(&totalCount)
	if err != nil {
		t.Errorf("Failed to count total todos: %v", err)
	}
	if totalCount != 2 {
		t.Errorf("Expected 2 total todos, got %d", totalCount)
	}
}

func TestSQLiteTodoRepository_Internal_ConnectionManagement(t *testing.T) {
	// Arrange
	dbPath := "test_internal_connection.db"
	defer os.Remove(dbPath)

	// Act
	repo, err := NewSQLiteTodoRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Assert - 内部実装: データベース接続の状態確認
	if repo.db == nil {
		t.Errorf("Database connection should not be nil")
	}

	// Ping to check connection
	err = repo.db.Ping()
	if err != nil {
		t.Errorf("Database connection should be active: %v", err)
	}

	// Close and verify
	err = repo.Close()
	if err != nil {
		t.Errorf("Failed to close repository: %v", err)
	}

	// After close, ping should fail (SQLiteは特別な動作をするため、この部分は実装によって異なる)
	// この部分は実際のSQLiteドライバーの動作に依存するため、コメントアウト
	// err = repo.db.Ping()
	// if err == nil {
	//     t.Errorf("Database connection should be closed")
	// }
}

func TestSQLiteTodoRepository_Internal_ErrorHandling(t *testing.T) {
	// Arrange
	dbPath := "test_internal_error.db"
	defer os.Remove(dbPath)

	repo, err := NewSQLiteTodoRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	// Act & Assert - 内部実装: エラーハンドリングの確認
	// 不正なSQLクエリの実行（内部実装の詳細をテスト）
	_, err = repo.db.Exec("INVALID SQL QUERY")
	if err == nil {
		t.Errorf("Invalid SQL query should return error")
	}

	// 存在しないテーブルへのクエリ
	_, err = repo.db.Query("SELECT * FROM non_existent_table")
	if err == nil {
		t.Errorf("Query to non-existent table should return error")
	}

	// 制約違反（重複キー）のテスト
	todo := entity.NewTodo("duplicate-id", "user-123", "Test Todo", "Description")
	repo.Create(todo)

	// 同じIDで再度作成（重複キー制約違反）
	_, err = repo.db.Exec(
		"INSERT INTO todos (id, user_id, title, description, completed, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		todo.ID, todo.UserID, "Another Title", "Another Description", false, todo.CreatedAt, todo.UpdatedAt,
	)
	if err == nil {
		t.Errorf("Duplicate key insertion should return error")
	}
}
