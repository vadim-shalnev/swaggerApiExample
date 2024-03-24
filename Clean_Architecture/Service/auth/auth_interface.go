package auth

import (
	"context"
	mod "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Models"
	repository "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Repository"
	"io"
)

type AuthServiceImpl struct {
	repo repository.Repository
}

type AuthService interface {
	Register(ctx context.Context, userBody io.ReadCloser) (mod.NewUserResponse, error)
	Login(ctx context.Context, userBody io.ReadCloser) (mod.NewUserResponse, error)
	UserInfoChecker(ctx context.Context, email, password, token string) (bool, bool, bool)
	TokenGenerate(ctx context.Context, email, password string) (string, error)
	VerifyToken(tokenString, searchIntoken string) (string, bool)
	RefreshToken(ctx context.Context, email, password string) string
}
