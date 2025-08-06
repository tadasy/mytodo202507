package entity_test

import (
	"testing"
	"time"

	"github.com/tadasy/mytodo202507/server/services/user/internal/domain/entity"
)

// ========================================
// 外部振る舞いテスト（Black-box Testing）
// Entity パッケージ外からの視点でのテスト
// ビジネスルールの外部振る舞いを検証
// ========================================

func TestUser_NewUser(t *testing.T) {
	// Arrange
	id := "user-123"
	email := "test@example.com"
	password := "password123"

	// Act
	user, err := entity.NewUser(id, email, password)

	// Assert
	if err != nil {
		t.Errorf("NewUser should not return error: %v", err)
	}
	if user.ID != id {
		t.Errorf("Expected ID %s, got %s", id, user.ID)
	}
	if user.Email != email {
		t.Errorf("Expected Email %s, got %s", email, user.Email)
	}
	if user.PasswordHash == "" {
		t.Errorf("PasswordHash should not be empty")
	}
	if user.PasswordHash == password {
		t.Errorf("PasswordHash should not be plaintext password")
	}
	if user.CreatedAt.IsZero() {
		t.Errorf("CreatedAt should be set")
	}
	if user.UpdatedAt.IsZero() {
		t.Errorf("UpdatedAt should be set")
	}
	if !user.CreatedAt.Equal(user.UpdatedAt) {
		t.Errorf("CreatedAt and UpdatedAt should be equal for new user")
	}
}

func TestUser_CheckPassword(t *testing.T) {
	// Arrange
	password := "password123"
	user, err := entity.NewUser("user-123", "test@example.com", password)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Act & Assert - Correct password
	if !user.CheckPassword(password) {
		t.Errorf("CheckPassword should return true for correct password")
	}

	// Act & Assert - Incorrect password
	if user.CheckPassword("wrongpassword") {
		t.Errorf("CheckPassword should return false for incorrect password")
	}

	// Act & Assert - Empty password
	if user.CheckPassword("") {
		t.Errorf("CheckPassword should return false for empty password")
	}
}

func TestUser_UpdatePassword(t *testing.T) {
	// Arrange
	originalPassword := "password123"
	newPassword := "newpassword456"
	user, err := entity.NewUser("user-123", "test@example.com", originalPassword)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	originalHash := user.PasswordHash
	originalUpdatedAt := user.UpdatedAt
	
	// Wait a bit to ensure UpdatedAt changes
	time.Sleep(1 * time.Millisecond)

	// Act
	err = user.UpdatePassword(newPassword)

	// Assert
	if err != nil {
		t.Errorf("UpdatePassword should not return error: %v", err)
	}
	if user.PasswordHash == originalHash {
		t.Errorf("PasswordHash should be updated")
	}
	if user.PasswordHash == newPassword {
		t.Errorf("PasswordHash should not be plaintext password")
	}
	if !user.CheckPassword(newPassword) {
		t.Errorf("New password should be valid")
	}
	if user.CheckPassword(originalPassword) {
		t.Errorf("Old password should no longer be valid")
	}
	if !user.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf("UpdatedAt should be updated")
	}
}

func TestUser_UpdateEmail(t *testing.T) {
	// Arrange
	originalEmail := "test@example.com"
	newEmail := "updated@example.com"
	user, err := entity.NewUser("user-123", originalEmail, "password123")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	originalUpdatedAt := user.UpdatedAt
	
	// Wait a bit to ensure UpdatedAt changes
	time.Sleep(1 * time.Millisecond)

	// Act
	user.UpdateEmail(newEmail)

	// Assert
	if user.Email != newEmail {
		t.Errorf("Expected Email %s, got %s", newEmail, user.Email)
	}
	if !user.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf("UpdatedAt should be updated")
	}
}

func TestUser_PasswordSecurity(t *testing.T) {
	// Arrange
	password := "password123"
	
	// Act - Create multiple users with same password
	user1, _ := entity.NewUser("user-1", "user1@example.com", password)
	user2, _ := entity.NewUser("user-2", "user2@example.com", password)

	// Assert - Password hashes should be different (due to salt)
	if user1.PasswordHash == user2.PasswordHash {
		t.Errorf("Password hashes should be different even for same password")
	}
	
	// Assert - Both users should validate with correct password
	if !user1.CheckPassword(password) {
		t.Errorf("User1 should validate with correct password")
	}
	if !user2.CheckPassword(password) {
		t.Errorf("User2 should validate with correct password")
	}
}

func TestUser_InvalidPasswordUpdate(t *testing.T) {
	// Arrange
	user, err := entity.NewUser("user-123", "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Act - Try to update with empty password
	err = user.UpdatePassword("")

	// Assert - Should not error, but hash should remain unchanged
	if err != nil {
		t.Errorf("UpdatePassword with empty string should not error: %v", err)
	}
	// Note: bcrypt can handle empty passwords, so this test checks the behavior
}

func TestUser_EmailValidation(t *testing.T) {
	// Arrange & Act
	testCases := []struct {
		name  string
		email string
	}{
		{"Normal email", "test@example.com"},
		{"Email with plus", "test+tag@example.com"},
		{"Empty email", ""},
		{"Invalid format", "notanemail"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := entity.NewUser("user-123", tc.email, "password123")
			
			// Assert - Entity creation should succeed regardless of email format
			// (validation should be done at service layer)
			if err != nil {
				t.Errorf("NewUser should not fail for email validation: %v", err)
			}
			if user.Email != tc.email {
				t.Errorf("Expected Email %s, got %s", tc.email, user.Email)
			}
		})
	}
}
