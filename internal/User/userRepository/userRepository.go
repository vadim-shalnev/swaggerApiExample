package userRepository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/vadim-shalnev/swaggerApiExample/Models"
	"log"
)

func NewUserRepository(db *sql.DB) *Userrepository {
	return &Userrepository{DB: db}
}

func (r *Userrepository) GetByID(ctx context.Context, id int) (Models.User, error) {
	var user Models.User

	err := r.DB.QueryRow("SELECT id, email, password, role, created_at, deleted_at FROM users WHERE id = $1", id).Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return Models.User{}, errors.New("User not found")
		}
		log.Println("Error getting User by id:", err)
		return Models.User{}, err
	}
	if user.DeletedAt != nil {
		return Models.User{}, errors.New("User not found")
	}
	log.Println("User found:", user)
	return user, nil
}

func (r *Userrepository) List(ctx context.Context) ([]Models.User, error) {
	var users []Models.User
	rows, err := r.DB.QueryContext(ctx, "SELECT id, email, password, role, created_at, deleted_at FROM users")
	if err != nil {
		log.Println("Error querying users:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user Models.User
		err = rows.Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.DeletedAt)
		if err != nil {
			log.Println("Error scanning User:", err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *Userrepository) Delete(ctx context.Context, userID int) error {
	_, err := r.DB.ExecContext(ctx, "UPDATE users SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1", userID)
	if err != nil {
		log.Println("Error deleting User:", err)
		return err
	}
	return nil
}
