package grpc

import (
	"context"
	"time"

	pb "github.com/tadasy/todo-app/proto"
	"github.com/tadasy/todo-app/server/services/user/internal/domain/service"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	userService *service.UserService
}

func NewUserServer(userService *service.UserService) *UserServer {
	return &UserServer{
		userService: userService,
	}
}

func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user, err := s.userService.CreateUser(req.Email, req.Password)
	if err != nil {
		return &pb.CreateUserResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.CreateUserResponse{
		User: &pb.User{
			Id:           user.ID,
			Email:        user.Email,
			PasswordHash: user.PasswordHash,
			CreatedAt:    user.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    user.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.userService.GetUser(req.Id)
	if err != nil {
		return &pb.GetUserResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			Id:           user.ID,
			Email:        user.Email,
			PasswordHash: user.PasswordHash,
			CreatedAt:    user.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    user.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *UserServer) AuthenticateUser(ctx context.Context, req *pb.AuthenticateUserRequest) (*pb.AuthenticateUserResponse, error) {
	user, err := s.userService.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		return &pb.AuthenticateUserResponse{
			Error: err.Error(),
		}, nil
	}

	// TODO: Generate JWT token
	token := "jwt-token-placeholder"

	return &pb.AuthenticateUserResponse{
		User: &pb.User{
			Id:           user.ID,
			Email:        user.Email,
			PasswordHash: user.PasswordHash,
			CreatedAt:    user.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    user.UpdatedAt.Format(time.RFC3339),
		},
		Token: token,
	}, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	user, err := s.userService.UpdateUser(req.Id, req.Email, req.Password)
	if err != nil {
		return &pb.UpdateUserResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.UpdateUserResponse{
		User: &pb.User{
			Id:           user.ID,
			Email:        user.Email,
			PasswordHash: user.PasswordHash,
			CreatedAt:    user.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    user.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *UserServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := s.userService.DeleteUser(req.Id)
	if err != nil {
		return &pb.DeleteUserResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &pb.DeleteUserResponse{
		Success: true,
	}, nil
}
