package userService

import (
	"context"
	mod "github.com/vadim-shalnev/swaggerApiExample/Models"
	"github.com/vadim-shalnev/swaggerApiExample/internal/User/userRepository"
)

type Userservice struct {
	repo userRepository.UserRepository
}

type UserService interface {
	GetUser(ctx context.Context, id string) (mod.NewUserResponse, error)
	DelUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context) ([]mod.User, error)
}
