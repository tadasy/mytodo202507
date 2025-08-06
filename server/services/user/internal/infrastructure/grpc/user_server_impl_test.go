package grpc

import (
	"context"
	"testing"

	pb "github.com/tadasy/mytodo202507/proto"
	"github.com/tadasy/mytodo202507/server/services/user/internal/domain/entity"
	"github.com/tadasy/mytodo202507/server/services/user/internal/domain/repository"
	"github.com/tadasy/mytodo202507/server/services/user/internal/domain/service"
)

// ========================================
// 内部実装テスト（White-box Testing）
// gRPCサーバーの内部実装詳細テスト
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

func TestUserServer_CreateUser_ServiceIntegration(t *testing.T) {
	// Arrange
	mockRepo := NewDetailedMockUserRepository()
	userService := service.NewUserService(mockRepo)
	server := NewUserServer(userService)
	ctx := context.Background()
	
	req := &pb.CreateUserRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Act
	resp, err := server.CreateUser(ctx, req)

	// Assert - 外部振る舞い
	if err != nil {
		t.Errorf("CreateUser should not return gRPC error: %v", err)
	}
	if resp.Error != "" {
		t.Errorf("Response should not contain error: %s", resp.Error)
	}

	// Assert - 内部実装詳細
	getCalls := mockRepo.GetGetCalls()
	if len(getCalls) != 1 {
		t.Errorf("Expected 1 GetByEmail call for duplicate check, got %d", len(getCalls))
	}
	if getCalls[0] != "EMAIL:test@example.com" {
		t.Errorf("Expected GetByEmail call with correct email")
	}

	createCalls := mockRepo.GetCreateCalls()
	if len(createCalls) != 1 {
		t.Errorf("Expected 1 Create call, got %d", len(createCalls))
	}
}

func TestUserServer_GetUser_CallSequence(t *testing.T) {
	// Arrange
	mockRepo := NewDetailedMockUserRepository()
	userService := service.NewUserService(mockRepo)
	server := NewUserServer(userService)
	ctx := context.Background()
	
	// Create user first
	user, _ := userService.CreateUser("test@example.com", "password123")
	mockRepo.ClearCalls()
	
	req := &pb.GetUserRequest{
		Id: user.ID,
	}

	// Act
	resp, err := server.GetUser(ctx, req)

	// Assert - 外部振る舞い
	if err != nil {
		t.Errorf("GetUser should not return gRPC error: %v", err)
	}
	if resp.Error != "" {
		t.Errorf("Response should not contain error: %s", resp.Error)
	}

	// Assert - 内部実装詳細
	getCalls := mockRepo.GetGetCalls()
	if len(getCalls) != 1 {
		t.Errorf("Expected 1 GetByID call, got %d", len(getCalls))
	}
	if getCalls[0] != "ID:"+user.ID {
		t.Errorf("Expected GetByID call with correct ID")
	}
}

func TestUserServer_AuthenticateUser_CallSequence(t *testing.T) {
	// Arrange
	mockRepo := NewDetailedMockUserRepository()
	userService := service.NewUserService(mockRepo)
	server := NewUserServer(userService)
	ctx := context.Background()
	
	email := "test@example.com"
	password := "password123"
	userService.CreateUser(email, password)
	mockRepo.ClearCalls()
	
	req := &pb.AuthenticateUserRequest{
		Email:    email,
		Password: password,
	}

	// Act
	resp, err := server.AuthenticateUser(ctx, req)

	// Assert - 外部振る舞い
	if err != nil {
		t.Errorf("AuthenticateUser should not return gRPC error: %v", err)
	}
	if resp.Error != "" {
		t.Errorf("Response should not contain error: %s", resp.Error)
	}
	if resp.Token == "" {
		t.Errorf("Token should be provided")
	}

	// Assert - 内部実装詳細
	getCalls := mockRepo.GetGetCalls()
	if len(getCalls) != 1 {
		t.Errorf("Expected 1 GetByEmail call, got %d", len(getCalls))
	}
	if getCalls[0] != "EMAIL:"+email {
		t.Errorf("Expected GetByEmail call with correct email")
	}
}

func TestUserServer_ServiceDependency(t *testing.T) {
	// Arrange
	mockRepo := NewDetailedMockUserRepository()
	userService := service.NewUserService(mockRepo)
	server := NewUserServer(userService)

	// Act & Assert - userServiceフィールドがプライベートなので直接アクセステストは不可
	// 代わりに動作確認でDependency Injectionをテスト
	ctx := context.Background()
	req := &pb.CreateUserRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	
	resp, err := server.CreateUser(ctx, req)

	if err != nil {
		t.Errorf("Server should work with injected service")
	}
	if resp.Error != "" {
		t.Errorf("Server should return success when service works correctly")
	}

	// 異なるサービスインスタンスでテスト
	repo2 := NewDetailedMockUserRepository()
	service2 := service.NewUserService(repo2)
	server2 := NewUserServer(service2)
	
	// 最初のサーバーで作ったユーザーは2つ目のサーバーでは見えない
	getReq := &pb.GetUserRequest{Id: resp.User.Id}
	getResp, err := server2.GetUser(ctx, getReq)
	if err != nil {
		t.Errorf("GetUser should not return gRPC error")
	}
	if getResp.Error == "" {
		t.Errorf("Different server instances should have separate service instances")
	}
}

func TestUserServer_ErrorHandling(t *testing.T) {
	// Arrange
	mockRepo := NewDetailedMockUserRepository()
	userService := service.NewUserService(mockRepo)
	server := NewUserServer(userService)
	ctx := context.Background()

	// Act & Assert - Service層のエラーがgRPC応答に正しく変換されることを確認
	req := &pb.GetUserRequest{
		Id: "nonexistent-id",
	}
	
	resp, err := server.GetUser(ctx, req)

	// gRPCエラーではなく、アプリケーションレベルのエラーとして返されること
	if err != nil {
		t.Errorf("Should not return gRPC error for application-level errors: %v", err)
	}
	if resp.Error == "" {
		t.Errorf("Should return application error in response")
	}
	if resp.User != nil {
		t.Errorf("User should be nil when error occurs")
	}
}

func TestUserServer_TokenPlaceholder(t *testing.T) {
	// Arrange
	mockRepo := NewDetailedMockUserRepository()
	userService := service.NewUserService(mockRepo)
	server := NewUserServer(userService)
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

	// Assert - 現在の実装ではプレースホルダーが返される
	if err != nil {
		t.Errorf("AuthenticateUser should not return gRPC error: %v", err)
	}
	if resp.Token != "jwt-token-placeholder" {
		t.Errorf("Expected placeholder token, got %s", resp.Token)
	}
	// TODO: 実際のJWT実装時にこのテストを更新
}
