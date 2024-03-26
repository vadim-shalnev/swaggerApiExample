package authService

import (
	"context"
	mod "github.com/vadim-shalnev/swaggerApiExample/Models"
	repository "github.com/vadim-shalnev/swaggerApiExample/internal/Repository"
)

type Authservice struct {
	repo repository.Repository
}

type AuthService interface {
	Register(ctx context.Context, regData mod.NewUserRequest) (mod.NewUserResponse, error)
	Login(ctx context.Context, loginData mod.NewUserRequest) (mod.NewUserResponse, error)
	UserInfoChecker(ctx context.Context, email, password, token string) (bool, bool, bool)
	TokenGenerate(ctx context.Context, email, password string) (string, error)
	VerifyToken(tokenString string) (string, string, bool)
	RefreshToken(ctx context.Context, email, password string) string
}
