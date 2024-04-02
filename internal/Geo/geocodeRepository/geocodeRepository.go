package geocodeRepository

import (
	"context"
	"database/sql"
	"github.com/vadim-shalnev/swaggerApiExample/Models"
	"github.com/vadim-shalnev/swaggerApiExample/internal/infrastructures/Cache"
	"log"
)

func NewGeocodeRepository(db *sql.DB, cache Cache.Cache) *Geocoderepository {
	return &Geocoderepository{DB: db, Cache: cache}
}
func (r *Geocoderepository) GetByEmail(ctx context.Context, email string) (Models.User, error) {
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

func (r *Geocoderepository) CacheChecker(ctx context.Context, query Models.RequestQuery) (Models.RequestAddress, error) {
	var cache Models.RequestAddress
	err := r.Cache.Get(ctx, query.Query, &cache)
	if err != nil {
		return Models.RequestAddress{}, err
	}
	return cache, nil
}

func (r *Geocoderepository) Insert(ctx context.Context, email string, query Models.RequestQuery, requestQuery Models.RequestAddress) error {
	// Начало транзакции
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // Откат транзакции в случае ошибки

	userID, err := r.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	var queryID int
	err = tx.QueryRow("INSERT INTO qwery_history (user_key, query, created_at, deleted_at) VALUES ($1, $2, CURRENT_TIMESTAMP, NULL) RETURNING id", userID.ID, query.Query).Scan(&queryID)
	if err != nil {
		log.Println("Error adding query:", err)
		return err
	}

	// Добавление результата в response_history
	_, err = tx.Exec("INSERT INTO response_history (qwery_key, address, lng, lat, created_at, deleted_at) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, NULL)", queryID, requestQuery.Addres, requestQuery.RequestSearch.Lng, requestQuery.RequestSearch.Lat)
	if err != nil {
		log.Println("Error adding result:", err)
		return err
	}

	// Подтверждение транзакции
	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction:", err)
		return err
	}

	return nil
}
