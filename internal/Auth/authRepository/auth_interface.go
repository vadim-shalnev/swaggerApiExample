package authRepository

import (
	"context"
	"database/sql"
	"github.com/vadim-shalnev/swaggerApiExample/Models"
)

type Authrepository struct {
	DB *sql.DB
}

type AuthRepository interface {
	CreateUser(ctx context.Context, user Models.NewUserRequest) error
	GetByEmail(ctx context.Context, email string) (Models.User, error)
}
