package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tadasy/mytodo202507/server/services/user/internal/domain/entity"
	"github.com/tadasy/mytodo202507/server/services/user/internal/domain/repository"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(email, password string) (*entity.User, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Generate UUID for user
	userID := uuid.New().String()

	// Create new user
	user, err := entity.NewUser(userID, email, password)
	if err != nil {
		return nil, err
	}

	// Save user
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUser(id string) (*entity.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) AuthenticateUser(email, password string) (*entity.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.CheckPassword(password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *UserService) UpdateUser(id, email, password string) (*entity.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if email != "" {
		user.UpdateEmail(email)
	}

	if password != "" {
		if err := user.UpdatePassword(password); err != nil {
			return nil, err
		}
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(id string) error {
	return s.userRepo.Delete(id)
}
