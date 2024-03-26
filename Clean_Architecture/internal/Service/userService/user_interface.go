package userService

import (
	"context"
	mod "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Models"
	repository "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/internal/Repository"
)

type UserServiceImpl struct {
	repo repository.Repository
}

type UserService interface {
	GetUser(ctx context.Context, id string) (mod.NewUserResponse, error)
	DelUser(ctx context.Context, id string) error
}
