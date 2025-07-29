package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/tadasy/mytodo202507/proto"
	"github.com/tadasy/mytodo202507/server/services/user/internal/domain/service"
	"github.com/tadasy/mytodo202507/server/services/user/internal/infrastructure/database"
	grpcServer "github.com/tadasy/mytodo202507/server/services/user/internal/infrastructure/grpc"
)

func main() {
	// Initialize database
	userRepo, err := database.NewSQLiteUserRepository("./users.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer userRepo.Close()

	// Initialize domain service
	userService := service.NewUserService(userRepo)

	// Initialize gRPC server
	userGRPCServer := grpcServer.NewUserServer(userService)

	// Create gRPC server
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, userGRPCServer)

	// Listen on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("User service starting on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
