package grpc_test

import (
	"context"
	"testing"
	"time"

	pb "github.com/tadasy/mytodo202507/proto"
	"github.com/tadasy/mytodo202507/server/services/user/internal/domain/entity"
	"github.com/tadasy/mytodo202507/server/services/user/internal/domain/repository"
	"github.com/tadasy/mytodo202507/server/services/user/internal/domain/service"
	"github.com/tadasy/mytodo202507/server/services/user/internal/infrastructure/grpc"
)

// ========================================
// 外部振る舞いテスト（Black-box Testing）
// gRPCプロトコル経由での外部振る舞いテスト
// 実装詳細に依存しない
// ========================================

// MockUserRepository は外部振る舞いテスト用のシンプルなモック
type MockUserRepository struct {
	users  map[string]*entity.User
	emails map[string]*entity.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:  make(map[string]*entity.User),
		emails: make(map[string]*entity.User),
	}
}

func (m *MockUserRepository) Create(user *entity.User) error {
	m.users[user.ID] = user
	m.emails[user.Email] = user
	return nil
}

func (m *MockUserRepository) GetByID(id string) (*entity.User, error) {
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, repository.ErrUserNotFound
}

func (m *MockUserRepository) GetByEmail(email string) (*entity.User, error) {
	if user, exists := m.emails[email]; exists {
		return user, nil
	}
	return nil, repository.ErrUserNotFound
}

func (m *MockUserRepository) Update(user *entity.User) error {
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

func (m *MockUserRepository) Delete(id string) error {
	if user, exists := m.users[id]; exists {
		delete(m.users, id)
		delete(m.emails, user.Email)
		return nil
	}
	return repository.ErrUserNotFound
}

func TestUserServer_CreateUser(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository()
	userService := service.NewUserService(mockRepo)
	server := grpc.NewUserServer(userService)
	ctx := context.Background()
	
	req := &pb.CreateUserRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Act
	resp, err := server.CreateUser(ctx, req)

	// Assert
	if err != nil {
		t.Errorf("CreateUser should not return gRPC error: %v", err)
	}
	if resp.Error != "" {
		t.Errorf("Response should not contain error: %s", resp.Error)
	}
	if resp.User == nil {
		t.Fatalf("User should not be nil")
	}
	if resp.User.Email != req.Email {
		t.Errorf("Expected Email %s, got %s", req.Email, resp.User.Email)
	}
	if resp.User.Id == "" {
		t.Errorf("User ID should be generated")
	}
	if resp.User.CreatedAt == "" {
		t.Errorf("CreatedAt should be set")
	}
	if resp.User.UpdatedAt == "" {
		t.Errorf("UpdatedAt should be set")
	}
}

func TestUserServer_GetUser(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository()
	userService := service.NewUserService(mockRepo)
	server := grpc.NewUserServer(userService)
	ctx := context.Background()
	
	// Create a user first
	user, _ := userService.CreateUser("test@example.com", "password123")
	
	req := &pb.GetUserRequest{
		Id: user.ID,
	}

	// Act
	resp, err := server.GetUser(ctx, req)

	// Assert
	if err != nil {
		t.Errorf("GetUser should not return gRPC error: %v", err)
	}
	if resp.Error != "" {
		t.Errorf("Response should not contain error: %s", resp.Error)
	}
	if resp.User == nil {
		t.Fatalf("User should not be nil")
	}
	if resp.User.Id != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, resp.User.Id)
	}
	if resp.User.Email != user.Email {
		t.Errorf("Expected Email %s, got %s", user.Email, resp.User.Email)
	}
}

func TestUserServer_GetUser_NotFound(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository()
	userService := service.NewUserService(mockRepo)
	server := grpc.NewUserServer(userService)
	ctx := context.Background()
	
	req := &pb.GetUserRequest{
		Id: "nonexistent-id",
	}

	// Act
	resp, err := server.GetUser(ctx, req)

	// Assert
	if err != nil {
		t.Errorf("GetUser should not return gRPC error: %v", err)
	}
	if resp.Error == "" {
		t.Errorf("Response should contain error for nonexistent user")
	}
	if resp.User != nil {
		t.Errorf("User should be nil when not found")
	}
}

func TestUserServer_AuthenticateUser(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository()
	userService := service.NewUserService(mockRepo)
	server := grpc.NewUserServer(userService)
	ctx := context.Background()
	
	email := "test@example.com"
	password := "password123"
	userService.CreateUser(email, password)
	
	req := &pb.AuthenticateUserRequest{
		Email:    email,
		Password: password,
	}

	// Act
	resp, err := server.AuthenticateUser(ctx, req)

	// Assert
	if err != nil {
		t.Errorf("AuthenticateUser should not return gRPC error: %v", err)
	}
	if resp.Error != "" {
		t.Errorf("Response should not contain error: %s", resp.Error)
	}
	if resp.User == nil {
		t.Fatalf("User should not be nil")
	}
	if resp.User.Email != email {
		t.Errorf("Expected Email %s, got %s", email, resp.User.Email)
	}
	if resp.Token == "" {
		t.Errorf("Token should be provided")
	}
}

func TestUserServer_TimeFormatting(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository()
	userService := service.NewUserService(mockRepo)
	server := grpc.NewUserServer(userService)
	ctx := context.Background()
	
	req := &pb.CreateUserRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Act
	resp, err := server.CreateUser(ctx, req)

	// Assert
	if err != nil {
		t.Errorf("CreateUser should not return gRPC error: %v", err)
	}
	
	// CreatedAt と UpdatedAt が RFC3339 フォーマットかチェック
	if resp.User.CreatedAt == "" {
		t.Errorf("CreatedAt should be formatted as RFC3339")
	}
	if resp.User.UpdatedAt == "" {
		t.Errorf("UpdatedAt should be formatted as RFC3339")
	}
	
	// 実際にパースできるかテスト
	_, err = time.Parse(time.RFC3339, resp.User.CreatedAt)
	if err != nil {
		t.Errorf("CreatedAt should be valid RFC3339 format: %v", err)
	}
	
	_, err = time.Parse(time.RFC3339, resp.User.UpdatedAt)
	if err != nil {
		t.Errorf("UpdatedAt should be valid RFC3339 format: %v", err)
	}
}
