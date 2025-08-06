package service_test

import (
	"testing"

	"github.com/tadasy/mytodo202507/server/services/user/internal/domain/entity"
	"github.com/tadasy/mytodo202507/server/services/user/internal/domain/repository"
	"github.com/tadasy/mytodo202507/server/services/user/internal/domain/service"
)

// ========================================
// 外部振る舞いテスト（Black-box Testing）
// Service パッケージ外からの視点でのテスト
// 公開インターフェース経由のテスト
// ========================================

// SimpleMockUserRepository は外部振る舞いテスト用のシンプルなモック
type SimpleMockUserRepository struct {
	users  map[string]*entity.User
	emails map[string]*entity.User
}

func NewSimpleMockUserRepository() *SimpleMockUserRepository {
	return &SimpleMockUserRepository{
		users:  make(map[string]*entity.User),
		emails: make(map[string]*entity.User),
	}
}

func (m *SimpleMockUserRepository) Create(user *entity.User) error {
	m.users[user.ID] = user
	m.emails[user.Email] = user
	return nil
}

func (m *SimpleMockUserRepository) GetByID(id string) (*entity.User, error) {
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, repository.ErrUserNotFound
}

func (m *SimpleMockUserRepository) GetByEmail(email string) (*entity.User, error) {
	if user, exists := m.emails[email]; exists {
		return user, nil
	}
	return nil, repository.ErrUserNotFound
}

func (m *SimpleMockUserRepository) Update(user *entity.User) error {
	if _, exists := m.users[user.ID]; !exists {
		return repository.ErrUserNotFound
	}
	
	// Remove old email mapping if email changed
	for email, u := range m.emails {
		if u.ID == user.ID && email != user.Email {
			delete(m.emails, email)
			break
		}
	}
	
	m.users[user.ID] = user
	m.emails[user.Email] = user
	return nil
}

func (m *SimpleMockUserRepository) Delete(id string) error {
	if user, exists := m.users[id]; exists {
		delete(m.users, id)
		delete(m.emails, user.Email)
		return nil
	}
	return repository.ErrUserNotFound
}

func TestUserService_CreateUser(t *testing.T) {
	// Arrange
	repo := NewSimpleMockUserRepository()
	userService := service.NewUserService(repo)
	email := "test@example.com"
	password := "password123"

	// Act
	user, err := userService.CreateUser(email, password)

	// Assert
	if err != nil {
		t.Errorf("CreateUser should not return error: %v", err)
	}
	if user == nil {
		t.Fatalf("User should not be nil")
	}
	if user.Email != email {
		t.Errorf("Expected Email %s, got %s", email, user.Email)
	}
	if user.ID == "" {
		t.Errorf("User ID should be generated")
	}
	if !user.CheckPassword(password) {
		t.Errorf("User should be created with correct password")
	}
}

func TestUserService_CreateUser_DuplicateEmail(t *testing.T) {
	// Arrange
	repo := NewSimpleMockUserRepository()
	userService := service.NewUserService(repo)
	email := "test@example.com"
	password := "password123"

	// Create first user
	userService.CreateUser(email, password)

	// Act - Try to create user with same email
	user, err := userService.CreateUser(email, "differentpassword")

	// Assert
	if err == nil {
		t.Errorf("CreateUser should return error for duplicate email")
	}
	if user != nil {
		t.Errorf("User should be nil when creation fails")
	}
}

func TestUserService_GetUser(t *testing.T) {
	// Arrange
	repo := NewSimpleMockUserRepository()
	userService := service.NewUserService(repo)
	email := "test@example.com"
	password := "password123"
	
	createdUser, _ := userService.CreateUser(email, password)

	// Act
	retrievedUser, err := userService.GetUser(createdUser.ID)

	// Assert
	if err != nil {
		t.Errorf("GetUser should not return error: %v", err)
	}
	if retrievedUser == nil {
		t.Fatalf("User should not be nil")
	}
	if retrievedUser.ID != createdUser.ID {
		t.Errorf("Expected ID %s, got %s", createdUser.ID, retrievedUser.ID)
	}
	if retrievedUser.Email != email {
		t.Errorf("Expected Email %s, got %s", email, retrievedUser.Email)
	}
}

func TestUserService_GetUser_NotFound(t *testing.T) {
	// Arrange
	repo := NewSimpleMockUserRepository()
	userService := service.NewUserService(repo)

	// Act
	user, err := userService.GetUser("nonexistent-id")

	// Assert
	if err == nil {
		t.Errorf("GetUser should return error for nonexistent user")
	}
	if user != nil {
		t.Errorf("User should be nil when not found")
	}
}

func TestUserService_AuthenticateUser(t *testing.T) {
	// Arrange
	repo := NewSimpleMockUserRepository()
	userService := service.NewUserService(repo)
	email := "test@example.com"
	password := "password123"
	
	userService.CreateUser(email, password)

	// Act
	user, err := userService.AuthenticateUser(email, password)

	// Assert
	if err != nil {
		t.Errorf("AuthenticateUser should not return error: %v", err)
	}
	if user == nil {
		t.Fatalf("User should not be nil")
	}
	if user.Email != email {
		t.Errorf("Expected Email %s, got %s", email, user.Email)
	}
}

func TestUserService_AuthenticateUser_InvalidEmail(t *testing.T) {
	// Arrange
	repo := NewSimpleMockUserRepository()
	userService := service.NewUserService(repo)

	// Act
	user, err := userService.AuthenticateUser("nonexistent@example.com", "password123")

	// Assert
	if err == nil {
		t.Errorf("AuthenticateUser should return error for invalid email")
	}
	if user != nil {
		t.Errorf("User should be nil for invalid credentials")
	}
}

func TestUserService_AuthenticateUser_InvalidPassword(t *testing.T) {
	// Arrange
	repo := NewSimpleMockUserRepository()
	userService := service.NewUserService(repo)
	email := "test@example.com"
	password := "password123"
	
	userService.CreateUser(email, password)

	// Act
	user, err := userService.AuthenticateUser(email, "wrongpassword")

	// Assert
	if err == nil {
		t.Errorf("AuthenticateUser should return error for invalid password")
	}
	if user != nil {
		t.Errorf("User should be nil for invalid credentials")
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	// Arrange
	repo := NewSimpleMockUserRepository()
	userService := service.NewUserService(repo)
	email := "test@example.com"
	password := "password123"
	
	createdUser, _ := userService.CreateUser(email, password)
	
	newEmail := "updated@example.com"
	newPassword := "newpassword456"

	// Act
	updatedUser, err := userService.UpdateUser(createdUser.ID, newEmail, newPassword)

	// Assert
	if err != nil {
		t.Errorf("UpdateUser should not return error: %v", err)
	}
	if updatedUser == nil {
		t.Fatalf("UpdatedUser should not be nil")
	}
	if updatedUser.Email != newEmail {
		t.Errorf("Expected Email %s, got %s", newEmail, updatedUser.Email)
	}
	if !updatedUser.CheckPassword(newPassword) {
		t.Errorf("User should have updated password")
	}
	if updatedUser.CheckPassword(password) {
		t.Errorf("User should not have old password")
	}
}

func TestUserService_UpdateUser_EmailOnly(t *testing.T) {
	// Arrange
	repo := NewSimpleMockUserRepository()
	userService := service.NewUserService(repo)
	email := "test@example.com"
	password := "password123"
	
	createdUser, _ := userService.CreateUser(email, password)
	newEmail := "updated@example.com"

	// Act
	updatedUser, err := userService.UpdateUser(createdUser.ID, newEmail, "")

	// Assert
	if err != nil {
		t.Errorf("UpdateUser should not return error: %v", err)
	}
	if updatedUser.Email != newEmail {
		t.Errorf("Expected Email %s, got %s", newEmail, updatedUser.Email)
	}
	if !updatedUser.CheckPassword(password) {
		t.Errorf("User should keep original password")
	}
}

func TestUserService_UpdateUser_PasswordOnly(t *testing.T) {
	// Arrange
	repo := NewSimpleMockUserRepository()
	userService := service.NewUserService(repo)
	email := "test@example.com"
	password := "password123"
	
	createdUser, _ := userService.CreateUser(email, password)
	newPassword := "newpassword456"

	// Act
	updatedUser, err := userService.UpdateUser(createdUser.ID, "", newPassword)

	// Assert
	if err != nil {
		t.Errorf("UpdateUser should not return error: %v", err)
	}
	if updatedUser.Email != email {
		t.Errorf("Expected Email %s, got %s", email, updatedUser.Email)
	}
	if !updatedUser.CheckPassword(newPassword) {
		t.Errorf("User should have updated password")
	}
}

func TestUserService_UpdateUser_NotFound(t *testing.T) {
	// Arrange
	repo := NewSimpleMockUserRepository()
	userService := service.NewUserService(repo)

	// Act
	user, err := userService.UpdateUser("nonexistent-id", "new@example.com", "newpassword")

	// Assert
	if err == nil {
		t.Errorf("UpdateUser should return error for nonexistent user")
	}
	if user != nil {
		t.Errorf("User should be nil when not found")
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	// Arrange
	repo := NewSimpleMockUserRepository()
	userService := service.NewUserService(repo)
	email := "test@example.com"
	password := "password123"
	
	createdUser, _ := userService.CreateUser(email, password)

	// Act
	err := userService.DeleteUser(createdUser.ID)

	// Assert
	if err != nil {
		t.Errorf("DeleteUser should not return error: %v", err)
	}
	
	// Verify user is deleted
	_, getErr := userService.GetUser(createdUser.ID)
	if getErr == nil {
		t.Errorf("User should be deleted")
	}
}

func TestUserService_DeleteUser_NotFound(t *testing.T) {
	// Arrange
	repo := NewSimpleMockUserRepository()
	userService := service.NewUserService(repo)

	// Act
	err := userService.DeleteUser("nonexistent-id")

	// Assert
	if err == nil {
		t.Errorf("DeleteUser should return error for nonexistent user")
	}
}
