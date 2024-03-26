package userService

import (
	"context"
	mod "github.com/vadim-shalnev/swaggerApiExample/Models"
	repository "github.com/vadim-shalnev/swaggerApiExample/internal/Repository"
)

type UserServiceImpl struct {
	repo repository.Repository
}

type UserService interface {
	GetUser(ctx context.Context, id string) (mod.NewUserResponse, error)
	DelUser(ctx context.Context, id string) error
}
