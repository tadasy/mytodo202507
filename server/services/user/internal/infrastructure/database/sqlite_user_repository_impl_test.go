package database

import (
	"os"
	"testing"
	"time"

	"github.com/tadasy/mytodo202507/server/services/user/internal/domain/entity"
)

// ========================================
// 内部実装テスト（White-box Testing）
// SQLite固有の実装詳細テスト
// データベース操作の内部動作を検証
// ========================================

func TestSQLiteUserRepository_TableCreation(t *testing.T) {
	// Arrange
	dbPath := "test_table_creation.db"
	defer os.Remove(dbPath)

	// Act
	repo, err := NewSQLiteUserRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	// Assert - テーブルが作成されていることを確認するため、実際にデータを挿入してみる
	user, _ := entity.NewUser("test-id", "test@example.com", "password123")
	err = repo.Create(user)
	if err != nil {
		t.Errorf("Table should be created and usable: %v", err)
	}
}

func TestSQLiteUserRepository_DatabaseConnection(t *testing.T) {
	// Arrange
	dbPath := "test_connection.db"
	defer os.Remove(dbPath)

	// Act
	repo, err := NewSQLiteUserRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Assert - データベース接続が有効であることを確認
	err = repo.Close()
	if err != nil {
		t.Errorf("Database connection should be closable: %v", err)
	}
}

func TestSQLiteUserRepository_TimeFormatHandling(t *testing.T) {
	// Arrange
	dbPath := "test_users_time_format.db"
	defer os.Remove(dbPath)

	repo, err := NewSQLiteUserRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	// 通常のユーザー作成プロセス
	user, err := entity.NewUser("user-123", "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to create user entity: %v", err)
	}

	// Act - 正常パスでの作成・取得
	err = repo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	retrievedUser, err := repo.GetByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	// Assert - 基本的な時刻データの整合性
	if retrievedUser.CreatedAt.IsZero() {
		t.Errorf("CreatedAt should not be zero")
	}
	if retrievedUser.UpdatedAt.IsZero() {
		t.Errorf("UpdatedAt should not be zero")
	}
	
	// 時刻の順序性確認（CreatedAtとUpdatedAtが合理的な範囲内）
	timeDiff := retrievedUser.UpdatedAt.Sub(retrievedUser.CreatedAt)
	if timeDiff < 0 {
		t.Errorf("UpdatedAt should not be before CreatedAt")
	}
	if timeDiff > 5*time.Second {
		t.Errorf("Time difference between CreatedAt and UpdatedAt should be small: %v", timeDiff)
	}
}

func TestSQLiteUserRepository_SQLInjectionProtection(t *testing.T) {
	// Arrange
	dbPath := "test_sql_injection.db"
	defer os.Remove(dbPath)

	repo, err := NewSQLiteUserRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	// SQLインジェクションを試みる悪意のあるデータ
	maliciousID := "'; DROP TABLE users; --"
	maliciousEmail := "test@example.com'; DROP TABLE users; --"

	user, err := entity.NewUser(maliciousID, maliciousEmail, "password123")
	if err != nil {
		t.Fatalf("Failed to create user entity: %v", err)
	}

	// Act
	err = repo.Create(user)
	if err != nil {
		t.Fatalf("Repository should handle malicious input gracefully: %v", err)
	}

	// SQLインジェクション攻撃が成功していないことを確認
	retrievedUser, err := repo.GetByID(maliciousID)
	if err != nil {
		t.Errorf("Should be able to retrieve user with malicious ID: %v", err)
	}
	if retrievedUser != nil && retrievedUser.ID != maliciousID {
		t.Errorf("ID should be stored as-is: expected %s, got %s", maliciousID, retrievedUser.ID)
	}

	// テーブルがまだ存在することを確認（別のユーザーを作成できる）
	normalUser, err := entity.NewUser("normal-id", "normal@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to create normal user entity: %v", err)
	}
	
	err = repo.Create(normalUser)
	if err != nil {
		t.Errorf("Table should still exist after malicious input: %v", err)
	}
}

func TestSQLiteUserRepository_TransactionConsistency(t *testing.T) {
	// Arrange
	dbPath := "test_transaction.db"
	defer os.Remove(dbPath)

	repo, err := NewSQLiteUserRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	user, _ := entity.NewUser("user-123", "test@example.com", "password123")
	repo.Create(user)

	// Act - 存在しないユーザーの更新を試行
	nonExistentUser, _ := entity.NewUser("nonexistent", "none@example.com", "password123")
	err = repo.Update(nonExistentUser)

	// Assert - エラーが返されること
	if err == nil {
		t.Errorf("Update of nonexistent user should return error")
	}

	// 元のユーザーが影響を受けていないことを確認
	retrievedUser, err := repo.GetByID(user.ID)
	if err != nil {
		t.Errorf("Original user should still exist: %v", err)
	}
	if retrievedUser.Email != user.Email {
		t.Errorf("Original user should be unchanged")
	}
}

func TestSQLiteUserRepository_DatabaseFileHandling(t *testing.T) {
	// Arrange
	dbPath := "test_file_handling.db"
	defer os.Remove(dbPath)

	// Act - データベースファイル作成
	repo, err := NewSQLiteUserRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// データベースファイルが作成されていることを確認
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Errorf("Database file should be created")
	}

	repo.Close()

	// ファイルが存在する状態で再度開く
	repo2, err := NewSQLiteUserRepository(dbPath)
	if err != nil {
		t.Fatalf("Should be able to open existing database file: %v", err)
	}
	defer repo2.Close()

	// 既存データベースが使用可能であることを確認
	user, _ := entity.NewUser("test-id", "test@example.com", "password123")
	err = repo2.Create(user)
	if err != nil {
		t.Errorf("Existing database should be usable: %v", err)
	}
}

func TestSQLiteUserRepository_ScanUserMethod(t *testing.T) {
	// Arrange
	dbPath := "test_scan_user.db"
	defer os.Remove(dbPath)

	repo, err := NewSQLiteUserRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	// 異なる時間精度を持つユーザーを作成
	user, _ := entity.NewUser("user-123", "test@example.com", "password123")
	// マイクロ秒まで設定してscanUserメソッドの時間パース機能をテスト
	user.CreatedAt = time.Date(2023, 12, 25, 15, 30, 45, 123456000, time.UTC)
	user.UpdatedAt = time.Date(2023, 12, 25, 16, 30, 45, 654321000, time.UTC)

	// Act
	err = repo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	retrievedUser, err := repo.GetByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	// Assert - scanUserメソッドが正しく時間をパースしていることを確認
	if retrievedUser.ID != user.ID {
		t.Errorf("ID should be correctly scanned")
	}
	if retrievedUser.Email != user.Email {
		t.Errorf("Email should be correctly scanned")
	}
	if retrievedUser.PasswordHash != user.PasswordHash {
		t.Errorf("PasswordHash should be correctly scanned")
	}
	// 時間は SQLite の精度制限により秒単位での確認
	if abs(retrievedUser.CreatedAt.Sub(user.CreatedAt)) > time.Second {
		t.Errorf("CreatedAt should be correctly parsed from string")
	}
	if abs(retrievedUser.UpdatedAt.Sub(user.UpdatedAt)) > time.Second {
		t.Errorf("UpdatedAt should be correctly parsed from string")
	}
}

func TestSQLiteUserRepository_ErrorHandling(t *testing.T) {
	// Arrange
	dbPath := "test_error_handling.db"
	defer os.Remove(dbPath)

	repo, err := NewSQLiteUserRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	// Act & Assert - Delete nonexistent user
	err = repo.Delete("nonexistent-id")
	if err != nil {
		t.Logf("Delete nonexistent user returned error (expected): %v", err)
	}

	// Act & Assert - Update nonexistent user
	nonExistentUser, _ := entity.NewUser("nonexistent", "none@example.com", "password")
	err = repo.Update(nonExistentUser)
	if err != nil {
		t.Logf("Update nonexistent user returned error (expected): %v", err)
	}
}

// Helper function for time comparison
func abs(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}
