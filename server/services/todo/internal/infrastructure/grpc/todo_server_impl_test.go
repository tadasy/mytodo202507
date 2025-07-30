package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	pb "github.com/tadasy/mytodo202507/proto"
	"github.com/tadasy/mytodo202507/server/services/todo/internal/domain/entity"
	"github.com/tadasy/mytodo202507/server/services/todo/internal/domain/service"
)

// DetailedMockRepository for implementation testing
type DetailedMockRepository struct {
	CreateCalled bool
	CreateInput  *entity.Todo
	CreateError  error

	GetByIDCalled bool
	GetByIDInput  []string // id, userID
	GetByIDReturn *entity.Todo
	GetByIDError  error

	ListByUserIDCalled bool
	ListByUserIDInput  string
	ListByUserIDReturn []*entity.Todo
	ListByUserIDError  error

	ListCompletedByUserIDCalled bool
	ListCompletedByUserIDInput  string
	ListCompletedByUserIDReturn []*entity.Todo
	ListCompletedByUserIDError  error

	UpdateCalled bool
	UpdateInput  *entity.Todo
	UpdateError  error

	DeleteCalled bool
	DeleteInput  []string // id, userID
	DeleteError  error

	// Internal state for verification
	todos map[string]*entity.Todo
}

func NewDetailedMockRepository() *DetailedMockRepository {
	return &DetailedMockRepository{
		todos: make(map[string]*entity.Todo),
	}
}

func (r *DetailedMockRepository) Create(todo *entity.Todo) error {
	r.CreateCalled = true
	r.CreateInput = todo
	if r.CreateError != nil {
		return r.CreateError
	}
	r.todos[todo.ID] = todo
	return nil
}

func (r *DetailedMockRepository) GetByID(id, userID string) (*entity.Todo, error) {
	r.GetByIDCalled = true
	r.GetByIDInput = []string{id, userID}
	if r.GetByIDError != nil {
		return nil, r.GetByIDError
	}
	if r.GetByIDReturn != nil {
		return r.GetByIDReturn, nil
	}
	todo, exists := r.todos[id]
	if !exists || todo.UserID != userID {
		return nil, nil
	}
	return todo, nil
}

func (r *DetailedMockRepository) ListByUserID(userID string) ([]*entity.Todo, error) {
	r.ListByUserIDCalled = true
	r.ListByUserIDInput = userID
	if r.ListByUserIDError != nil {
		return nil, r.ListByUserIDError
	}
	if r.ListByUserIDReturn != nil {
		return r.ListByUserIDReturn, nil
	}

	var todos []*entity.Todo
	for _, todo := range r.todos {
		if todo.UserID == userID {
			todos = append(todos, todo)
		}
	}
	return todos, nil
}

func (r *DetailedMockRepository) ListCompletedByUserID(userID string) ([]*entity.Todo, error) {
	r.ListCompletedByUserIDCalled = true
	r.ListCompletedByUserIDInput = userID
	if r.ListCompletedByUserIDError != nil {
		return nil, r.ListCompletedByUserIDError
	}
	if r.ListCompletedByUserIDReturn != nil {
		return r.ListCompletedByUserIDReturn, nil
	}

	var todos []*entity.Todo
	for _, todo := range r.todos {
		if todo.UserID == userID && todo.Completed {
			todos = append(todos, todo)
		}
	}
	return todos, nil
}

func (r *DetailedMockRepository) Update(todo *entity.Todo) error {
	r.UpdateCalled = true
	r.UpdateInput = todo
	if r.UpdateError != nil {
		return r.UpdateError
	}
	r.todos[todo.ID] = todo
	return nil
}

func (r *DetailedMockRepository) Delete(id, userID string) error {
	r.DeleteCalled = true
	r.DeleteInput = []string{id, userID}
	if r.DeleteError != nil {
		return r.DeleteError
	}
	delete(r.todos, id)
	return nil
}

func (r *DetailedMockRepository) Reset() {
	r.CreateCalled = false
	r.CreateInput = nil
	r.CreateError = nil

	r.GetByIDCalled = false
	r.GetByIDInput = nil
	r.GetByIDReturn = nil
	r.GetByIDError = nil

	r.ListByUserIDCalled = false
	r.ListByUserIDInput = ""
	r.ListByUserIDReturn = nil
	r.ListByUserIDError = nil

	r.ListCompletedByUserIDCalled = false
	r.ListCompletedByUserIDInput = ""
	r.ListCompletedByUserIDReturn = nil
	r.ListCompletedByUserIDError = nil

	r.UpdateCalled = false
	r.UpdateInput = nil
	r.UpdateError = nil

	r.DeleteCalled = false
	r.DeleteInput = nil
	r.DeleteError = nil

	r.todos = make(map[string]*entity.Todo)
}

func createTodoServerWithMockRepo(repo *DetailedMockRepository) *TodoServer {
	todoService := service.NewTodoService(repo)
	return NewTodoServer(todoService)
}

func TestTodoServer_CreateTodo_Implementation(t *testing.T) {
	repo := NewDetailedMockRepository()
	server := createTodoServerWithMockRepo(repo)

	req := &pb.CreateTodoRequest{
		Title:       "Test Todo",
		Description: "Test Description",
		UserId:      "user123",
	}

	resp, err := server.CreateTodo(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateTodo failed: %v", err)
	}

	// Verify repository was called
	if !repo.CreateCalled {
		t.Error("Expected Create to be called on repository")
	}

	// Verify input passed to repository
	if repo.CreateInput == nil {
		t.Fatal("Expected CreateInput to be set")
	}
	if repo.CreateInput.Title != "Test Todo" {
		t.Errorf("Expected input title 'Test Todo', got %s", repo.CreateInput.Title)
	}
	if repo.CreateInput.UserID != "user123" {
		t.Errorf("Expected input user ID 'user123', got %s", repo.CreateInput.UserID)
	}
	if repo.CreateInput.Description != "Test Description" {
		t.Errorf("Expected input description 'Test Description', got %s", repo.CreateInput.Description)
	}

	// Verify response
	if resp.Todo.Title != "Test Todo" {
		t.Errorf("Expected response title 'Test Todo', got %s", resp.Todo.Title)
	}
	if resp.Todo.UserId != "user123" {
		t.Errorf("Expected response user ID 'user123', got %s", resp.Todo.UserId)
	}
}

func TestTodoServer_CreateTodo_RepositoryError(t *testing.T) {
	repo := NewDetailedMockRepository()
	repo.CreateError = errors.New("repository error")
	server := createTodoServerWithMockRepo(repo)

	req := &pb.CreateTodoRequest{
		Title:  "Test Todo",
		UserId: "user123",
	}

	resp, err := server.CreateTodo(context.Background(), req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.Error == "" {
		t.Error("Expected error to be set in response")
	}

	if !repo.CreateCalled {
		t.Error("Expected Create to be called on repository")
	}
}

func TestTodoServer_GetTodo_Implementation(t *testing.T) {
	repo := NewDetailedMockRepository()
	server := createTodoServerWithMockRepo(repo)

	// Pre-setup a todo in repository
	existingTodo := &entity.Todo{
		ID:          "todo123",
		Title:       "Test Todo",
		Description: "Test Description",
		UserID:      "user123",
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	repo.GetByIDReturn = existingTodo

	req := &pb.GetTodoRequest{
		Id:     "todo123",
		UserId: "user123",
	}

	resp, err := server.GetTodo(context.Background(), req)
	if err != nil {
		t.Fatalf("GetTodo failed: %v", err)
	}

	// Verify repository was called with correct ID
	if !repo.GetByIDCalled {
		t.Error("Expected GetByID to be called on repository")
	}
	if len(repo.GetByIDInput) != 2 || repo.GetByIDInput[0] != "todo123" {
		t.Errorf("Expected input ID 'todo123', got %v", repo.GetByIDInput)
	}

	// Verify response conversion
	if resp.Todo.Id != "todo123" {
		t.Errorf("Expected response ID 'todo123', got %s", resp.Todo.Id)
	}
	if resp.Todo.Title != "Test Todo" {
		t.Errorf("Expected response title 'Test Todo', got %s", resp.Todo.Title)
	}
}

func TestTodoServer_ListTodos_Implementation(t *testing.T) {
	repo := NewDetailedMockRepository()
	server := createTodoServerWithMockRepo(repo)

	// Setup mock return
	repo.ListByUserIDReturn = []*entity.Todo{
		{ID: "todo1", Title: "Todo 1", UserID: "user123", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "todo2", Title: "Todo 2", UserID: "user123", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	req := &pb.ListTodosRequest{
		UserId: "user123",
	}

	resp, err := server.ListTodos(context.Background(), req)
	if err != nil {
		t.Fatalf("ListTodos failed: %v", err)
	}

	// Verify repository was called with correct user ID
	if !repo.ListByUserIDCalled {
		t.Error("Expected ListByUserID to be called on repository")
	}
	if repo.ListByUserIDInput != "user123" {
		t.Errorf("Expected input user ID 'user123', got %s", repo.ListByUserIDInput)
	}

	// Verify response conversion
	if len(resp.Todos) != 2 {
		t.Errorf("Expected 2 todos in response, got %d", len(resp.Todos))
	}
}

func TestTodoServer_UpdateTodo_Implementation(t *testing.T) {
	repo := NewDetailedMockRepository()
	server := createTodoServerWithMockRepo(repo)

	// Pre-setup existing todo
	existingTodo := &entity.Todo{
		ID:          "todo123",
		Title:       "Original Title",
		Description: "Original Description",
		UserID:      "user123",
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	repo.todos["todo123"] = existingTodo

	req := &pb.UpdateTodoRequest{
		Id:          "todo123",
		UserId:      "user123",
		Title:       "Updated Title",
		Description: "Updated Description",
	}

	resp, err := server.UpdateTodo(context.Background(), req)
	if err != nil {
		t.Fatalf("UpdateTodo failed: %v", err)
	}

	// Verify repository was called
	if !repo.UpdateCalled {
		t.Error("Expected Update to be called on repository")
	}

	// Verify input passed to repository
	if repo.UpdateInput == nil {
		t.Fatal("Expected UpdateInput to be set")
	}
	if repo.UpdateInput.ID != "todo123" {
		t.Errorf("Expected input ID 'todo123', got %s", repo.UpdateInput.ID)
	}
	if repo.UpdateInput.Title != "Updated Title" {
		t.Errorf("Expected input title 'Updated Title', got %s", repo.UpdateInput.Title)
	}
	if repo.UpdateInput.Description != "Updated Description" {
		t.Errorf("Expected input description 'Updated Description', got %s", repo.UpdateInput.Description)
	}

	// Verify response conversion
	if resp.Todo.Title != "Updated Title" {
		t.Errorf("Expected response title 'Updated Title', got %s", resp.Todo.Title)
	}
}

func TestTodoServer_DeleteTodo_Implementation(t *testing.T) {
	repo := NewDetailedMockRepository()
	server := createTodoServerWithMockRepo(repo)

	req := &pb.DeleteTodoRequest{
		Id:     "todo123",
		UserId: "user123",
	}

	resp, err := server.DeleteTodo(context.Background(), req)
	if err != nil {
		t.Fatalf("DeleteTodo failed: %v", err)
	}

	// Verify repository was called with correct ID
	if !repo.DeleteCalled {
		t.Error("Expected Delete to be called on repository")
	}
	if len(repo.DeleteInput) != 2 || repo.DeleteInput[0] != "todo123" {
		t.Errorf("Expected input ID 'todo123', got %v", repo.DeleteInput)
	}

	// Verify response
	if !resp.Success {
		t.Error("Expected success to be true")
	}
}

func TestTodoServer_todoToProto_Implementation(t *testing.T) {
	server := &TodoServer{}

	createdAt := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC)

	entity := &entity.Todo{
		ID:          "todo123",
		Title:       "Test Todo",
		Description: "Test Description",
		UserID:      "user123",
		Completed:   true,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	proto := server.todoToProto(entity)

	if proto.Id != entity.ID {
		t.Errorf("Expected ID %s, got %s", entity.ID, proto.Id)
	}
	if proto.Title != entity.Title {
		t.Errorf("Expected title %s, got %s", entity.Title, proto.Title)
	}
	if proto.Description != entity.Description {
		t.Errorf("Expected description %s, got %s", entity.Description, proto.Description)
	}
	if proto.UserId != entity.UserID {
		t.Errorf("Expected user ID %s, got %s", entity.UserID, proto.UserId)
	}
	if proto.Completed != entity.Completed {
		t.Errorf("Expected completed %v, got %v", entity.Completed, proto.Completed)
	}
	if proto.CreatedAt != entity.CreatedAt.Format(time.RFC3339) {
		t.Errorf("Expected created at %s, got %s", entity.CreatedAt.Format(time.RFC3339), proto.CreatedAt)
	}
	if proto.UpdatedAt != entity.UpdatedAt.Format(time.RFC3339) {
		t.Errorf("Expected updated at %s, got %s", entity.UpdatedAt.Format(time.RFC3339), proto.UpdatedAt)
	}
}
