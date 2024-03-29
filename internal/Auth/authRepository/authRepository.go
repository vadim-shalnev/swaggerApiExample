package authRepository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/vadim-shalnev/swaggerApiExample/Models"
	"log"
)

// Таблица ответов dаdata содержит внешний ключ к таблице запросов, которая в свою очередь хранит внешний ключ таблицы пользователей
func NewAuthrepository(db *sql.DB) *Authrepository {
	return &Authrepository{DB: db}
}

func (r *Authrepository) CreateUser(ctx context.Context, user Models.NewUserRequest) error {
	_, err := r.DB.Exec("INSERT INTO users (email, password, role, created_at, deleted_at) VALUES ($1, $2, $3, CURRENT_TIMESTAMP, NULL)", user.Email, user.Password, user.Role)
	if err != nil {
		log.Println("Error creating User:", err)
		return errors.New("failed to create User")
	}
	return nil
}
func (r *Authrepository) GetByEmail(ctx context.Context, email string) (Models.User, error) {
	var user Models.User

	err := r.DB.QueryRowContext(ctx, "SELECT id, email, password, role, created_at, deleted_at FROM users WHERE email = $1", email).Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("User not found", email)
			return Models.User{}, err
		}
		log.Println("Error getting User by email:", err)
		return Models.User{}, err
	}
	if user.DeletedAt != nil {
		log.Println("User deleted", email)
		return Models.User{}, err
	}

	return user, nil
}
