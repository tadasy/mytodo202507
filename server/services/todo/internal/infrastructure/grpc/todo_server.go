package grpc

import (
	"context"
	"time"

	pb "github.com/tadasy/mytodo202507/proto"
	"github.com/tadasy/mytodo202507/server/services/todo/internal/domain/entity"
	"github.com/tadasy/mytodo202507/server/services/todo/internal/domain/service"
)

type TodoServer struct {
	pb.UnimplementedTodoServiceServer
	todoService *service.TodoService
}

func NewTodoServer(todoService *service.TodoService) *TodoServer {
	return &TodoServer{
		todoService: todoService,
	}
}

func (s *TodoServer) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.CreateTodoResponse, error) {
	todo, err := s.todoService.CreateTodo(req.UserId, req.Title, req.Description)
	if err != nil {
		return &pb.CreateTodoResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.CreateTodoResponse{
		Todo: s.todoToProto(todo),
	}, nil
}

func (s *TodoServer) GetTodo(ctx context.Context, req *pb.GetTodoRequest) (*pb.GetTodoResponse, error) {
	todo, err := s.todoService.GetTodo(req.Id, req.UserId)
	if err != nil {
		return &pb.GetTodoResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.GetTodoResponse{
		Todo: s.todoToProto(todo),
	}, nil
}

func (s *TodoServer) ListTodos(ctx context.Context, req *pb.ListTodosRequest) (*pb.ListTodosResponse, error) {
	var todos []*pb.Todo

	if req.CompletedOnly {
		todoEntities, err := s.todoService.ListCompletedTodos(req.UserId)
		if err != nil {
			return &pb.ListTodosResponse{
				Error: err.Error(),
			}, nil
		}

		for _, todo := range todoEntities {
			todos = append(todos, s.todoToProto(todo))
		}
	} else {
		todoEntities, err := s.todoService.ListTodos(req.UserId)
		if err != nil {
			return &pb.ListTodosResponse{
				Error: err.Error(),
			}, nil
		}

		for _, todo := range todoEntities {
			todos = append(todos, s.todoToProto(todo))
		}
	}

	return &pb.ListTodosResponse{
		Todos: todos,
	}, nil
}

func (s *TodoServer) UpdateTodo(ctx context.Context, req *pb.UpdateTodoRequest) (*pb.UpdateTodoResponse, error) {
	todo, err := s.todoService.UpdateTodo(req.Id, req.UserId, req.Title, req.Description)
	if err != nil {
		return &pb.UpdateTodoResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.UpdateTodoResponse{
		Todo: s.todoToProto(todo),
	}, nil
}

func (s *TodoServer) DeleteTodo(ctx context.Context, req *pb.DeleteTodoRequest) (*pb.DeleteTodoResponse, error) {
	err := s.todoService.DeleteTodo(req.Id, req.UserId)
	if err != nil {
		return &pb.DeleteTodoResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &pb.DeleteTodoResponse{
		Success: true,
	}, nil
}

func (s *TodoServer) MarkTodoComplete(ctx context.Context, req *pb.MarkTodoCompleteRequest) (*pb.MarkTodoCompleteResponse, error) {
	todo, err := s.todoService.MarkTodoComplete(req.Id, req.UserId, req.Completed)
	if err != nil {
		return &pb.MarkTodoCompleteResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.MarkTodoCompleteResponse{
		Todo: s.todoToProto(todo),
	}, nil
}

func (s *TodoServer) ListCompletedTodos(ctx context.Context, req *pb.ListCompletedTodosRequest) (*pb.ListCompletedTodosResponse, error) {
	todoEntities, err := s.todoService.ListCompletedTodos(req.UserId)
	if err != nil {
		return &pb.ListCompletedTodosResponse{
			Error: err.Error(),
		}, nil
	}

	var todos []*pb.Todo
	for _, todo := range todoEntities {
		todos = append(todos, s.todoToProto(todo))
	}

	return &pb.ListCompletedTodosResponse{
		Todos: todos,
	}, nil
}

func (s *TodoServer) todoToProto(todo *entity.Todo) *pb.Todo {
	pbTodo := &pb.Todo{
		Id:          todo.ID,
		UserId:      todo.UserID,
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
		CreatedAt:   todo.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   todo.UpdatedAt.Format(time.RFC3339),
	}

	if todo.CompletedAt != nil {
		pbTodo.CompletedAt = todo.CompletedAt.Format(time.RFC3339)
	}

	return pbTodo
}
