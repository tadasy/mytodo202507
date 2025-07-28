package clients

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/tadasy/todo-app/proto"
	"github.com/tadasy/todo-app/server/bff/internal/models"
)

type UserServiceClient struct {
	client pb.UserServiceClient
	conn   *grpc.ClientConn
}

func NewUserServiceClient(address string) (*UserServiceClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %v", err)
	}

	client := pb.NewUserServiceClient(conn)
	return &UserServiceClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *UserServiceClient) CreateUser(ctx context.Context, email, password string) (*models.User, error) {
	resp, err := c.client.CreateUser(ctx, &pb.CreateUserRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, fmt.Errorf(resp.Error)
	}

	return c.protoUserToModel(resp.User), nil
}

func (c *UserServiceClient) AuthenticateUser(ctx context.Context, email, password string) (*models.User, string, error) {
	resp, err := c.client.AuthenticateUser(ctx, &pb.AuthenticateUserRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, "", err
	}

	if resp.Error != "" {
		return nil, "", fmt.Errorf(resp.Error)
	}

	return c.protoUserToModel(resp.User), resp.Token, nil
}

func (c *UserServiceClient) GetUser(ctx context.Context, id string) (*models.User, error) {
	resp, err := c.client.GetUser(ctx, &pb.GetUserRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, fmt.Errorf(resp.Error)
	}

	return c.protoUserToModel(resp.User), nil
}

func (c *UserServiceClient) Close() error {
	return c.conn.Close()
}

func (c *UserServiceClient) protoUserToModel(pbUser *pb.User) *models.User {
	createdAt, _ := time.Parse(time.RFC3339, pbUser.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, pbUser.UpdatedAt)

	return &models.User{
		ID:        pbUser.Id,
		Email:     pbUser.Email,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

// TodoServiceClient represents a client for the Todo service
type TodoServiceClient struct {
	client pb.TodoServiceClient
	conn   *grpc.ClientConn
}

func NewTodoServiceClient(address string) (*TodoServiceClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to todo service: %v", err)
	}

	client := pb.NewTodoServiceClient(conn)
	return &TodoServiceClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *TodoServiceClient) CreateTodo(ctx context.Context, userID, title, description string) (*models.Todo, error) {
	resp, err := c.client.CreateTodo(ctx, &pb.CreateTodoRequest{
		UserId:      userID,
		Title:       title,
		Description: description,
	})
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, fmt.Errorf(resp.Error)
	}

	return c.protoTodoToModel(resp.Todo), nil
}

func (c *TodoServiceClient) GetTodo(ctx context.Context, id, userID string) (*models.Todo, error) {
	resp, err := c.client.GetTodo(ctx, &pb.GetTodoRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, fmt.Errorf(resp.Error)
	}

	return c.protoTodoToModel(resp.Todo), nil
}

func (c *TodoServiceClient) ListTodos(ctx context.Context, userID string) ([]*models.Todo, error) {
	resp, err := c.client.ListTodos(ctx, &pb.ListTodosRequest{
		UserId:        userID,
		CompletedOnly: false,
	})
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, fmt.Errorf(resp.Error)
	}

	var todos []*models.Todo
	for _, pbTodo := range resp.Todos {
		todos = append(todos, c.protoTodoToModel(pbTodo))
	}

	return todos, nil
}

func (c *TodoServiceClient) ListCompletedTodos(ctx context.Context, userID string) ([]*models.Todo, error) {
	resp, err := c.client.ListCompletedTodos(ctx, &pb.ListCompletedTodosRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, fmt.Errorf(resp.Error)
	}

	var todos []*models.Todo
	for _, pbTodo := range resp.Todos {
		todos = append(todos, c.protoTodoToModel(pbTodo))
	}

	return todos, nil
}

func (c *TodoServiceClient) UpdateTodo(ctx context.Context, id, userID, title, description string) (*models.Todo, error) {
	resp, err := c.client.UpdateTodo(ctx, &pb.UpdateTodoRequest{
		Id:          id,
		UserId:      userID,
		Title:       title,
		Description: description,
	})
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, fmt.Errorf(resp.Error)
	}

	return c.protoTodoToModel(resp.Todo), nil
}

func (c *TodoServiceClient) MarkTodoComplete(ctx context.Context, id, userID string, completed bool) (*models.Todo, error) {
	resp, err := c.client.MarkTodoComplete(ctx, &pb.MarkTodoCompleteRequest{
		Id:        id,
		UserId:    userID,
		Completed: completed,
	})
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, fmt.Errorf(resp.Error)
	}

	return c.protoTodoToModel(resp.Todo), nil
}

func (c *TodoServiceClient) DeleteTodo(ctx context.Context, id, userID string) error {
	resp, err := c.client.DeleteTodo(ctx, &pb.DeleteTodoRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		return err
	}

	if resp.Error != "" {
		return fmt.Errorf(resp.Error)
	}

	return nil
}

func (c *TodoServiceClient) Close() error {
	return c.conn.Close()
}

func (c *TodoServiceClient) protoTodoToModel(pbTodo *pb.Todo) *models.Todo {
	createdAt, _ := time.Parse(time.RFC3339, pbTodo.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, pbTodo.UpdatedAt)

	todo := &models.Todo{
		ID:          pbTodo.Id,
		UserID:      pbTodo.UserId,
		Title:       pbTodo.Title,
		Description: pbTodo.Description,
		Completed:   pbTodo.Completed,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	if pbTodo.CompletedAt != "" {
		completedAt, _ := time.Parse(time.RFC3339, pbTodo.CompletedAt)
		todo.CompletedAt = &completedAt
	}

	return todo
}
