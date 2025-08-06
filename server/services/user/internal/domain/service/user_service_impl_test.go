package service

import (
	"testing"

	"github.com/tadasy/mytodo202507/server/services/user/internal/domain/entity"
	"github.com/tadasy/mytodo202507/server/services/user/internal/domain/repository"
)

// ========================================
// 内部実装テスト（White-box Testing）
// Service パッケージ内部の実装詳細テスト
// プライベートメンバーへのアクセス可能
// ========================================

// DetailedMockUserRepository は内部実装テスト用の詳細なモック
type DetailedMockUserRepository struct {
	users       map[string]*entity.User
	emails      map[string]*entity.User
	createCalls []string
	getCalls    []string
	updateCalls []string
	deleteCalls []string
}

func NewDetailedMockUserRepository() *DetailedMockUserRepository {
	return &DetailedMockUserRepository{
		users:       make(map[string]*entity.User),
		emails:      make(map[string]*entity.User),
		createCalls: make([]string, 0),
		getCalls:    make([]string, 0),
		updateCalls: make([]string, 0),
		deleteCalls: make([]string, 0),
	}
}

func (m *DetailedMockUserRepository) Create(user *entity.User) error {
	m.createCalls = append(m.createCalls, user.ID)
	m.users[user.ID] = user
	m.emails[user.Email] = user
	return nil
}

func (m *DetailedMockUserRepository) GetByID(id string) (*entity.User, error) {
	m.getCalls = append(m.getCalls, "ID:"+id)
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, repository.ErrUserNotFound
}

func (m *DetailedMockUserRepository) GetByEmail(email string) (*entity.User, error) {
	m.getCalls = append(m.getCalls, "EMAIL:"+email)
	if user, exists := m.emails[email]; exists {
		return user, nil
	}
	return nil, repository.ErrUserNotFound
}

func (m *DetailedMockUserRepository) Update(user *entity.User) error {
	m.updateCalls = append(m.updateCalls, user.ID)
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

func (m *DetailedMockUserRepository) Delete(id string) error {
	m.deleteCalls = append(m.deleteCalls, id)
	if user, exists := m.users[id]; exists {
		delete(m.users, id)
		delete(m.emails, user.Email)
		return nil
	}
	return repository.ErrUserNotFound
}

// Helper methods for testing
func (m *DetailedMockUserRepository) GetCreateCalls() []string   { return m.createCalls }
func (m *DetailedMockUserRepository) GetGetCalls() []string      { return m.getCalls }
func (m *DetailedMockUserRepository) GetUpdateCalls() []string   { return m.updateCalls }
func (m *DetailedMockUserRepository) GetDeleteCalls() []string   { return m.deleteCalls }
func (m *DetailedMockUserRepository) ClearCalls() {
	m.createCalls = m.createCalls[:0]
	m.getCalls = m.getCalls[:0]
	m.updateCalls = m.updateCalls[:0]
	m.deleteCalls = m.deleteCalls[:0]
}

func TestUserService_CreateUser_ImplementationDetails(t *testing.T) {
	// Arrange
	repo := NewDetailedMockUserRepository()
	userService := NewUserService(repo)
	email := "test@example.com"
	password := "password123"

	// Act
	user, err := userService.CreateUser(email, password)

	// Assert - 外部振る舞い
	if err != nil {
		t.Errorf("CreateUser should not return error: %v", err)
	}

	// Assert - 内部実装詳細
	createCalls := repo.GetCreateCalls()
	if len(createCalls) != 1 {
		t.Errorf("Expected 1 Create call, got %d", len(createCalls))
	}
	if createCalls[0] != user.ID {
		t.Errorf("Expected Create call with ID %s, got %s", user.ID, createCalls[0])
	}

	getCalls := repo.GetGetCalls()
	if len(getCalls) != 1 {
		t.Errorf("Expected 1 GetByEmail call, got %d", len(getCalls))
	}
	if getCalls[0] != "EMAIL:"+email {
		t.Errorf("Expected GetByEmail call with %s, got %s", email, getCalls[0])
	}
}

func TestUserService_CreateUser_DuplicateEmailCheck(t *testing.T) {
	// Arrange
	repo := NewDetailedMockUserRepository()
	userService := NewUserService(repo)
	email := "test@example.com"

	// Create first user
	userService.CreateUser(email, "password123")
	repo.ClearCalls()

	// Act - Try to create duplicate
	user, err := userService.CreateUser(email, "password456")

	// Assert - 外部振る舞い
	if err == nil {
		t.Errorf("CreateUser should return error for duplicate email")
	}
	if user != nil {
		t.Errorf("User should be nil for duplicate email")
	}

	// Assert - 内部実装詳細
	getCalls := repo.GetGetCalls()
	if len(getCalls) != 1 {
		t.Errorf("Expected 1 GetByEmail call, got %d", len(getCalls))
	}
	
	createCalls := repo.GetCreateCalls()
	if len(createCalls) != 0 {
		t.Errorf("Expected 0 Create calls for duplicate, got %d", len(createCalls))
	}
}

func TestUserService_AuthenticateUser_CallSequence(t *testing.T) {
	// Arrange
	repo := NewDetailedMockUserRepository()
	userService := NewUserService(repo)
	email := "test@example.com"
	password := "password123"
	
	userService.CreateUser(email, password)
	repo.ClearCalls()

	// Act
	user, err := userService.AuthenticateUser(email, password)

	// Assert - 外部振る舞い
	if err != nil {
		t.Errorf("AuthenticateUser should not return error: %v", err)
	}
	if user == nil {
		t.Fatalf("User should not be nil")
	}

	// Assert - 内部実装詳細
	getCalls := repo.GetGetCalls()
	if len(getCalls) != 1 {
		t.Errorf("Expected 1 GetByEmail call, got %d", len(getCalls))
	}
	if getCalls[0] != "EMAIL:"+email {
		t.Errorf("Expected GetByEmail call with %s, got %s", email, getCalls[0])
	}
}

func TestUserService_UpdateUser_CallSequence(t *testing.T) {
	// Arrange
	repo := NewDetailedMockUserRepository()
	userService := NewUserService(repo)
	email := "test@example.com"
	password := "password123"
	
	createdUser, _ := userService.CreateUser(email, password)
	repo.ClearCalls()
	
	newEmail := "updated@example.com"
	newPassword := "newpassword456"

	// Act
	updatedUser, err := userService.UpdateUser(createdUser.ID, newEmail, newPassword)

	// Assert - 外部振る舞い
	if err != nil {
		t.Errorf("UpdateUser should not return error: %v", err)
	}
	if updatedUser == nil {
		t.Fatalf("UpdatedUser should not be nil")
	}

	// Assert - 内部実装詳細
	getCalls := repo.GetGetCalls()
	if len(getCalls) != 1 {
		t.Errorf("Expected 1 GetByID call, got %d", len(getCalls))
	}
	if getCalls[0] != "ID:"+createdUser.ID {
		t.Errorf("Expected GetByID call with %s, got %s", createdUser.ID, getCalls[0])
	}
	
	updateCalls := repo.GetUpdateCalls()
	if len(updateCalls) != 1 {
		t.Errorf("Expected 1 Update call, got %d", len(updateCalls))
	}
	if updateCalls[0] != createdUser.ID {
		t.Errorf("Expected Update call with ID %s, got %s", createdUser.ID, updateCalls[0])
	}
}

func TestUserService_UpdateUser_EmptyFieldHandling(t *testing.T) {
	// Arrange
	repo := NewDetailedMockUserRepository()
	userService := NewUserService(repo)
	email := "test@example.com"
	password := "password123"
	
	createdUser, _ := userService.CreateUser(email, password)
	originalEmail := createdUser.Email
	originalPasswordHash := createdUser.PasswordHash
	repo.ClearCalls()

	// Act - Update with empty fields
	updatedUser, err := userService.UpdateUser(createdUser.ID, "", "")

	// Assert - 外部振る舞い
	if err != nil {
		t.Errorf("UpdateUser should not return error: %v", err)
	}
	if updatedUser.Email != originalEmail {
		t.Errorf("Email should remain unchanged when empty")
	}
	if updatedUser.PasswordHash != originalPasswordHash {
		t.Errorf("Password should remain unchanged when empty")
	}

	// Assert - 内部実装詳細（Updateは常に呼ばれる）
	updateCalls := repo.GetUpdateCalls()
	if len(updateCalls) != 1 {
		t.Errorf("Expected 1 Update call even for empty fields, got %d", len(updateCalls))
	}
}

func TestUserService_DeleteUser_CallSequence(t *testing.T) {
	// Arrange
	repo := NewDetailedMockUserRepository()
	userService := NewUserService(repo)
	email := "test@example.com"
	password := "password123"
	
	createdUser, _ := userService.CreateUser(email, password)
	repo.ClearCalls()

	// Act
	err := userService.DeleteUser(createdUser.ID)

	// Assert - 外部振る舞い
	if err != nil {
		t.Errorf("DeleteUser should not return error: %v", err)
	}

	// Assert - 内部実装詳細
	deleteCalls := repo.GetDeleteCalls()
	if len(deleteCalls) != 1 {
		t.Errorf("Expected 1 Delete call, got %d", len(deleteCalls))
	}
	if deleteCalls[0] != createdUser.ID {
		t.Errorf("Expected Delete call with ID %s, got %s", createdUser.ID, deleteCalls[0])
	}
}

func TestUserService_RepositoryDependency(t *testing.T) {
	// Arrange
	repo := NewDetailedMockUserRepository()
	userService := NewUserService(repo)

	// Act & Assert - userRepoフィールドがプライベートなので直接アクセステストは不可
	// 代わりに動作確認でDependency Injectionをテスト
	email := "test@example.com"
	user, err := userService.CreateUser(email, "password123")

	if err != nil {
		t.Errorf("Service should work with injected repository")
	}
	if user == nil {
		t.Errorf("Service should return user when repository works correctly")
	}

	// 異なるリポジトリインスタンスでテスト
	repo2 := NewDetailedMockUserRepository()
	userService2 := NewUserService(repo2)
	
	// 最初のサービスで作ったユーザーは2つ目のサービスでは見えない
	_, err = userService2.GetUser(user.ID)
	if err == nil {
		t.Errorf("Different service instances should have separate repository instances")
	}
}
