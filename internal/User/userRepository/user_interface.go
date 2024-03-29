package userRepository

import (
	"context"
	"database/sql"
	"github.com/vadim-shalnev/swaggerApiExample/Models"
)

type Userrepository struct {
	DB *sql.DB
}

type UserRepository interface {
	GetByID(ctx context.Context, id int) (Models.User, error)
	List(ctx context.Context) ([]Models.User, error)
	Delete(ctx context.Context, userID int) error
}
