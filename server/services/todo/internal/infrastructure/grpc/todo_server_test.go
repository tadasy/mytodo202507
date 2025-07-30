package grpc_test

import (
	"context"
	"testing"

	pb "github.com/tadasy/mytodo202507/proto"
	"github.com/tadasy/mytodo202507/server/services/todo/internal/domain/entity"
	"github.com/tadasy/mytodo202507/server/services/todo/internal/domain/repository"
	"github.com/tadasy/mytodo202507/server/services/todo/internal/domain/service"
	grpcServer "github.com/tadasy/mytodo202507/server/services/todo/internal/infrastructure/grpc"
)

// SimpleMockRepository for external testing
type SimpleMockRepository struct {
	todos map[string]*entity.Todo
}

func NewSimpleMockRepository() *SimpleMockRepository {
	return &SimpleMockRepository{
		todos: make(map[string]*entity.Todo),
	}
}

func (r *SimpleMockRepository) Create(todo *entity.Todo) error {
	r.todos[todo.ID] = todo
	return nil
}

func (r *SimpleMockRepository) GetByID(id, userID string) (*entity.Todo, error) {
	todo, exists := r.todos[id]
	if !exists || todo.UserID != userID {
		return nil, nil
	}
	return todo, nil
}

func (r *SimpleMockRepository) ListByUserID(userID string) ([]*entity.Todo, error) {
	var userTodos []*entity.Todo
	for _, todo := range r.todos {
		if todo.UserID == userID {
			userTodos = append(userTodos, todo)
		}
	}
	return userTodos, nil
}

func (r *SimpleMockRepository) ListCompletedByUserID(userID string) ([]*entity.Todo, error) {
	var userTodos []*entity.Todo
	for _, todo := range r.todos {
		if todo.UserID == userID && todo.Completed {
			userTodos = append(userTodos, todo)
		}
	}
	return userTodos, nil
}

func (r *SimpleMockRepository) Update(todo *entity.Todo) error {
	r.todos[todo.ID] = todo
	return nil
}

func (r *SimpleMockRepository) Delete(id, userID string) error {
	todo, exists := r.todos[id]
	if exists && todo.UserID == userID {
		delete(r.todos, id)
	}
	return nil
}

func createTodoServer(repo repository.TodoRepository) *grpcServer.TodoServer {
	todoService := service.NewTodoService(repo)
	return grpcServer.NewTodoServer(todoService)
}

func TestTodoServer_CreateTodo_Behavior(t *testing.T) {
	ctx := context.Background()
	repo := NewSimpleMockRepository()
	server := createTodoServer(repo)

	req := &pb.CreateTodoRequest{
		Title:       "Test Todo",
		Description: "Test Description",
		UserId:      "user123",
	}

	resp, err := server.CreateTodo(ctx, req)
	if err != nil {
		t.Fatalf("CreateTodo failed: %v", err)
	}

	if resp.Todo.Title != "Test Todo" {
		t.Errorf("Expected title 'Test Todo', got %s", resp.Todo.Title)
	}
	if resp.Todo.Description != "Test Description" {
		t.Errorf("Expected description 'Test Description', got %s", resp.Todo.Description)
	}
	if resp.Todo.UserId != "user123" {
		t.Errorf("Expected user ID 'user123', got %s", resp.Todo.UserId)
	}
	if resp.Todo.Completed {
		t.Error("Expected completed to be false")
	}
}

func TestTodoServer_GetTodo_Behavior(t *testing.T) {
	ctx := context.Background()
	repo := NewSimpleMockRepository()
	server := createTodoServer(repo)

	// First create a todo
	createReq := &pb.CreateTodoRequest{
		Title:       "Test Todo",
		Description: "Test Description",
		UserId:      "user123",
	}

	createResp, err := server.CreateTodo(ctx, createReq)
	if err != nil {
		t.Fatalf("CreateTodo failed: %v", err)
	}

	// Now get the todo
	getReq := &pb.GetTodoRequest{
		Id:     createResp.Todo.Id,
		UserId: "user123",
	}

	getResp, err := server.GetTodo(ctx, getReq)
	if err != nil {
		t.Fatalf("GetTodo failed: %v", err)
	}

	if getResp.Todo.Title != "Test Todo" {
		t.Errorf("Expected title 'Test Todo', got %s", getResp.Todo.Title)
	}
}

func TestTodoServer_ListTodos_UserIsolation(t *testing.T) {
	ctx := context.Background()
	repo := NewSimpleMockRepository()
	server := createTodoServer(repo)

	// Create todos for different users
	todos := []struct {
		title  string
		userID string
	}{
		{"User1 Todo1", "user1"},
		{"User1 Todo2", "user1"},
		{"User2 Todo1", "user2"},
	}

	for _, todo := range todos {
		req := &pb.CreateTodoRequest{
			Title:  todo.title,
			UserId: todo.userID,
		}
		_, err := server.CreateTodo(ctx, req)
		if err != nil {
			t.Fatalf("CreateTodo failed: %v", err)
		}
	}

	// List todos for user1
	listReq := &pb.ListTodosRequest{
		UserId: "user1",
	}

	listResp, err := server.ListTodos(ctx, listReq)
	if err != nil {
		t.Fatalf("ListTodos failed: %v", err)
	}

	if len(listResp.Todos) != 2 {
		t.Errorf("Expected 2 todos for user1, got %d", len(listResp.Todos))
	}

	for _, todo := range listResp.Todos {
		if todo.UserId != "user1" {
			t.Errorf("Expected user ID 'user1', got %s", todo.UserId)
		}
	}
}

func TestTodoServer_UpdateTodo_Behavior(t *testing.T) {
	ctx := context.Background()
	repo := NewSimpleMockRepository()
	server := createTodoServer(repo)

	// Create a todo
	createReq := &pb.CreateTodoRequest{
		Title:       "Original Title",
		Description: "Original Description",
		UserId:      "user123",
	}

	createResp, err := server.CreateTodo(ctx, createReq)
	if err != nil {
		t.Fatalf("CreateTodo failed: %v", err)
	}

	// Update the todo
	updateReq := &pb.UpdateTodoRequest{
		Id:          createResp.Todo.Id,
		UserId:      "user123",
		Title:       "Updated Title",
		Description: "Updated Description",
	}

	updateResp, err := server.UpdateTodo(ctx, updateReq)
	if err != nil {
		t.Fatalf("UpdateTodo failed: %v", err)
	}

	if updateResp.Todo.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got %s", updateResp.Todo.Title)
	}
	if updateResp.Todo.Description != "Updated Description" {
		t.Errorf("Expected description 'Updated Description', got %s", updateResp.Todo.Description)
	}
}

func TestTodoServer_DeleteTodo_Behavior(t *testing.T) {
	ctx := context.Background()
	repo := NewSimpleMockRepository()
	server := createTodoServer(repo)

	// Create a todo
	createReq := &pb.CreateTodoRequest{
		Title:  "Todo to Delete",
		UserId: "user123",
	}

	createResp, err := server.CreateTodo(ctx, createReq)
	if err != nil {
		t.Fatalf("CreateTodo failed: %v", err)
	}

	// Delete the todo
	deleteReq := &pb.DeleteTodoRequest{
		Id:     createResp.Todo.Id,
		UserId: "user123",
	}

	_, err = server.DeleteTodo(ctx, deleteReq)
	if err != nil {
		t.Fatalf("DeleteTodo failed: %v", err)
	}

	// Try to get the deleted todo
	getReq := &pb.GetTodoRequest{
		Id:     createResp.Todo.Id,
		UserId: "user123",
	}

	getResp, err := server.GetTodo(ctx, getReq)
	if err != nil {
		t.Fatalf("GetTodo failed: %v", err)
	}

	// After deletion, the todo should not be found (implementation returns nil Todo)
	if getResp.Todo != nil {
		t.Error("Expected todo to be nil after deletion")
	}
}
