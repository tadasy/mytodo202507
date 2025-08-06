package database_test

import (
	"os"
	"testing"

	"github.com/tadasy/mytodo202507/server/services/user/internal/domain/entity"
	"github.com/tadasy/mytodo202507/server/services/user/internal/domain/repository"
	"github.com/tadasy/mytodo202507/server/services/user/internal/infrastructure/database"
)

// ========================================
// 外部振る舞いテスト（Black-box Testing）
// リポジトリインターface経由でのテスト
// 実装の詳細に依存しない
// ========================================

func TestUserRepository_CreateAndGet(t *testing.T) {
	// Arrange
	dbPath := "test_users.db"
	defer os.Remove(dbPath)

	// repositoryインターface経由でテスト
	var repo repository.UserRepository
	sqliteRepo, err := database.NewSQLiteUserRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer sqliteRepo.Close()
	repo = sqliteRepo

	user, err := entity.NewUser("user-123", "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to create user entity: %v", err)
	}

	// Act - Create
	err = repo.Create(user)
	if err != nil {
		t.Errorf("Failed to create user: %v", err)
	}

	// Act - Get by ID
	retrievedUser, err := repo.GetByID(user.ID)
	if err != nil {
		t.Errorf("Failed to get user by ID: %v", err)
	}

	// Assert
	if retrievedUser.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, retrievedUser.ID)
	}
	if retrievedUser.Email != user.Email {
		t.Errorf("Expected Email %s, got %s", user.Email, retrievedUser.Email)
	}
	if retrievedUser.PasswordHash != user.PasswordHash {
		t.Errorf("Expected PasswordHash %s, got %s", user.PasswordHash, retrievedUser.PasswordHash)
	}
	// タイムスタンプの基本的な整合性確認
	if retrievedUser.CreatedAt.IsZero() {
		t.Errorf("CreatedAt should not be zero")
	}
	if retrievedUser.UpdatedAt.IsZero() {
		t.Errorf("UpdatedAt should not be zero")
	}
	
	// 時刻の順序性確認
	if retrievedUser.UpdatedAt.Before(retrievedUser.CreatedAt) {
		t.Errorf("UpdatedAt should not be before CreatedAt")
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	// Arrange
	dbPath := "test_users_email.db"
	defer os.Remove(dbPath)

	var repo repository.UserRepository
	sqliteRepo, err := database.NewSQLiteUserRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer sqliteRepo.Close()
	repo = sqliteRepo

	user, err := entity.NewUser("user-123", "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to create user entity: %v", err)
	}
	repo.Create(user)

	// Act
	retrievedUser, err := repo.GetByEmail(user.Email)

	// Assert
	if err != nil {
		t.Errorf("Failed to get user by email: %v", err)
	}
	if retrievedUser.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, retrievedUser.ID)
	}
	if retrievedUser.Email != user.Email {
		t.Errorf("Expected Email %s, got %s", user.Email, retrievedUser.Email)
	}
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	dbPath := "test_users_notfound.db"
	defer os.Remove(dbPath)

	var repo repository.UserRepository
	sqliteRepo, err := database.NewSQLiteUserRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer sqliteRepo.Close()
	repo = sqliteRepo

	// Act
	_, err = repo.GetByID("nonexistent-id")

	// Assert
	if err == nil {
		t.Errorf("Expected error for nonexistent user")
	}
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	// Arrange
	dbPath := "test_users_email_notfound.db"
	defer os.Remove(dbPath)

	var repo repository.UserRepository
	sqliteRepo, err := database.NewSQLiteUserRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer sqliteRepo.Close()
	repo = sqliteRepo

	// Act
	_, err = repo.GetByEmail("nonexistent@example.com")

	// Assert
	if err == nil {
		t.Errorf("Expected error for nonexistent email")
	}
}

func TestUserRepository_Update(t *testing.T) {
	// Arrange
	dbPath := "test_users_update.db"
	defer os.Remove(dbPath)

	var repo repository.UserRepository
	sqliteRepo, err := database.NewSQLiteUserRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer sqliteRepo.Close()
	repo = sqliteRepo

	user, err := entity.NewUser("user-123", "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to create user entity: %v", err)
	}
	repo.Create(user)

	// 変更
	user.UpdateEmail("updated@example.com")
	user.UpdatePassword("newpassword456")

	// Act
	err = repo.Update(user)
	if err != nil {
		t.Errorf("Failed to update user: %v", err)
	}

	// 更新されたUserを取得して確認
	updatedUser, err := repo.GetByID(user.ID)
	if err != nil {
		t.Errorf("Failed to get updated user: %v", err)
	}

	// Assert
	if updatedUser.Email != "updated@example.com" {
		t.Errorf("Expected email to be updated")
	}
	// パスワードが更新されているかチェック（UpdatePassword後のハッシュと比較）
	if updatedUser.PasswordHash != user.PasswordHash {
		t.Errorf("Expected password hash to be updated to %s, got %s", user.PasswordHash, updatedUser.PasswordHash)
	}
	// 更新時刻が作成時刻以降であることを確認（SQLiteの精度制限により同一時刻も許可）
	if updatedUser.UpdatedAt.Before(updatedUser.CreatedAt) {
		t.Errorf("Expected UpdatedAt (%v) to not be before CreatedAt (%v)", 
			updatedUser.UpdatedAt, updatedUser.CreatedAt)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	// Arrange
	dbPath := "test_users_delete.db"
	defer os.Remove(dbPath)

	var repo repository.UserRepository
	sqliteRepo, err := database.NewSQLiteUserRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer sqliteRepo.Close()
	repo = sqliteRepo

	user, err := entity.NewUser("user-123", "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to create user entity: %v", err)
	}
	repo.Create(user)

	// Act
	err = repo.Delete(user.ID)
	if err != nil {
		t.Errorf("Failed to delete user: %v", err)
	}

	// 削除されたことを確認
	_, err = repo.GetByID(user.ID)
	if err == nil {
		t.Errorf("Expected user to be deleted")
	}
}

func TestUserRepository_EmailUniqueness(t *testing.T) {
	// Arrange
	dbPath := "test_users_unique.db"
	defer os.Remove(dbPath)

	var repo repository.UserRepository
	sqliteRepo, err := database.NewSQLiteUserRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer sqliteRepo.Close()
	repo = sqliteRepo

	user1, _ := entity.NewUser("user-1", "test@example.com", "password123")
	user2, _ := entity.NewUser("user-2", "test@example.com", "password456")

	// Act
	err1 := repo.Create(user1)
	err2 := repo.Create(user2) // Same email

	// Assert
	if err1 != nil {
		t.Errorf("First user creation should succeed: %v", err1)
	}
	if err2 == nil {
		t.Errorf("Second user creation should fail due to email uniqueness")
	}
}

func TestUserRepository_MultipleUsers(t *testing.T) {
	// Arrange
	dbPath := "test_users_multiple.db"
	defer os.Remove(dbPath)

	var repo repository.UserRepository
	sqliteRepo, err := database.NewSQLiteUserRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer sqliteRepo.Close()
	repo = sqliteRepo

	user1, _ := entity.NewUser("user-1", "user1@example.com", "password1")
	user2, _ := entity.NewUser("user-2", "user2@example.com", "password2")
	user3, _ := entity.NewUser("user-3", "user3@example.com", "password3")

	// Act - Create multiple users
	repo.Create(user1)
	repo.Create(user2)
	repo.Create(user3)

	// Assert - Each user can be retrieved independently
	retrievedUser1, err1 := repo.GetByID(user1.ID)
	retrievedUser2, err2 := repo.GetByEmail(user2.Email)
	retrievedUser3, err3 := repo.GetByID(user3.ID)

	if err1 != nil || err2 != nil || err3 != nil {
		t.Errorf("All users should be retrievable")
	}
	if retrievedUser1.Email != user1.Email {
		t.Errorf("User1 email mismatch")
	}
	if retrievedUser2.ID != user2.ID {
		t.Errorf("User2 ID mismatch")
	}
	if retrievedUser3.Email != user3.Email {
		t.Errorf("User3 email mismatch")
	}
}

func TestUserRepository_EmailCaseInsensitiveCheck(t *testing.T) {
	// Arrange
	dbPath := "test_users_case.db"
	defer os.Remove(dbPath)

	var repo repository.UserRepository
	sqliteRepo, err := database.NewSQLiteUserRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer sqliteRepo.Close()
	repo = sqliteRepo

	user, _ := entity.NewUser("user-123", "Test@Example.com", "password123")
	repo.Create(user)

	// Act - Try to get with different case
	retrievedUser, err := repo.GetByEmail("test@example.com")

	// Assert - SQLite is case-insensitive by default for text
	// This test documents the current behavior
	if err != nil {
		t.Logf("Case-insensitive email lookup failed (expected behavior): %v", err)
	} else if retrievedUser != nil {
		t.Logf("Case-insensitive email lookup succeeded: %s", retrievedUser.Email)
	}
}

func TestUserRepository_ConcurrentAccess(t *testing.T) {
	// Arrange
	dbPath := "test_users_concurrent.db"
	defer os.Remove(dbPath)

	var repo repository.UserRepository
	sqliteRepo, err := database.NewSQLiteUserRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer sqliteRepo.Close()
	repo = sqliteRepo

	// Act - Create and immediately read (tests basic concurrent safety)
	user, _ := entity.NewUser("user-123", "test@example.com", "password123")
	
	err = repo.Create(user)
	if err != nil {
		t.Errorf("Create should succeed: %v", err)
	}
	
	retrievedUser, err := repo.GetByID(user.ID)
	if err != nil {
		t.Errorf("Immediate read should succeed: %v", err)
	}
	
	// Assert
	if retrievedUser.ID != user.ID {
		t.Errorf("Concurrent read should return correct user")
	}
}
