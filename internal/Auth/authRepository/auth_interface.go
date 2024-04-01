package authRepository

import (
	"database/sql"
	"github.com/vadim-shalnev/swaggerApiExample/Models"
)

type Authrepository struct {
	DB *sql.DB
}

type AuthRepository interface {
	CreateUser(user Models.NewUserRequest) error
	GetByEmail(email string) (Models.User, error)
}
