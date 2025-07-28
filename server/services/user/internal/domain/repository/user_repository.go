package repository

import (
	"github.com/tadasy/todo-app/server/services/user/internal/domain/entity"
)

type UserRepository interface {
	Create(user *entity.User) error
	GetByID(id string) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	Update(user *entity.User) error
	Delete(id string) error
}
