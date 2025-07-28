package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/tadasy/todo-app/proto"
	"github.com/tadasy/todo-app/server/services/todo/internal/domain/service"
	"github.com/tadasy/todo-app/server/services/todo/internal/infrastructure/database"
	grpcServer "github.com/tadasy/todo-app/server/services/todo/internal/infrastructure/grpc"
)

func main() {
	// Initialize database
	todoRepo, err := database.NewSQLiteTodoRepository("./todos.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer todoRepo.Close()

	// Initialize domain service
	todoService := service.NewTodoService(todoRepo)

	// Initialize gRPC server
	todoGRPCServer := grpcServer.NewTodoServer(todoService)

	// Create gRPC server
	s := grpc.NewServer()
	pb.RegisterTodoServiceServer(s, todoGRPCServer)

	// Listen on port 50052
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Todo service starting on port 50052...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
