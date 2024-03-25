package Repository

import (
	"context"
	"database/sql"
	"errors"
	mod "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Models"
	"log"
)

type RepositoryImpl struct {
	DB *sql.DB
}

type Repository interface {
	CreateUser(ctx context.Context, user mod.NewUserRequest) error
	GetByEmail(ctx context.Context, email string) (mod.User, error)
	GetByID(ctx context.Context, id int) (mod.User, error)
	Insert(ctx context.Context, email string, query mod.RequestQuery, requestQuery mod.RequestAddress) error
	Delete(ctx context.Context, userID int) error
	CacheChecker(ctx context.Context, email string, historyCount int) ([]mod.SearchHistory, error)
	//Select(ctx context.Context, query mod.UserRequest) error
	//Update(ctx context.Context, query mod.RequestUser) error
	//CheckEmail(ctx context.Context, email string) bool
	//CheckPassword(ctx context.Context, password string) bool
	//GetAll(ctx context.Context, query mod.RequestUser, colomsCount int) ([]mod.RequestAddress, error)
	//List(ctx context.Context, query mod.RequestUser) ([]mod.RequestAddress, error)
}

func NewRepositoryImpl(db *sql.DB) *RepositoryImpl {
	return &RepositoryImpl{DB: db}
}

func (r *RepositoryImpl) CreateUser(ctx context.Context, user mod.NewUserRequest) error {
	_, err := r.DB.Exec("INSERT INTO users (email, password, role, created_at, deleted_at) VALUES ($1, $2, $3, CURRENT_TIMESTAMP, NULL)", user.Email, user.Password, user.Role)
	if err != nil {
		log.Println("Error creating user:", err)
		return errors.New("failed to create user")
	}
	return nil
}
func (r *RepositoryImpl) GetByEmail(ctx context.Context, email string) (mod.User, error) {
	var user mod.User

	err := r.DB.QueryRowContext(ctx, "SELECT id, email, password, role, created_at, deleted_at FROM users WHERE email = $1", email).Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return mod.User{}, errors.New("user not found")
		}
		log.Println("Error getting user by email:", err)
		return mod.User{}, err
	}
	if user.DeletedAt != nil {
		return mod.User{}, errors.New("user not found")
	}

	return user, nil
}

func (r *RepositoryImpl) GetByID(ctx context.Context, id int) (mod.User, error) {
	var user mod.User

	err := r.DB.QueryRow("SELECT id, email, password, role, created_at, deleted_at FROM users WHERE id = $1", id).Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return mod.User{}, errors.New("user not found")
		}
		log.Println("Error getting user by email:", err)
		return mod.User{}, err
	}
	if user.DeletedAt != nil {
		return mod.User{}, errors.New("user not found")
	}
	return user, nil
}
func (r *RepositoryImpl) Insert(ctx context.Context, email string, query mod.RequestQuery, requestQuery mod.RequestAddress) error {
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
	err = tx.QueryRow("INSERT INTO qwery_history (user_key, query, created_date, deleted_date) VALUES ($1, $2, CURRENT_TIMESTAMP, NULL) RETURNING id", userID.ID, query.Query).Scan(&queryID)
	if err != nil {
		log.Println("Error adding query:", err)
		return err
	}

	// Добавление результата в response_history
	_, err = tx.Exec("INSERT INTO response_history (qwery_key, address, lng, lat, created_date, deleted_date) VALUES ($1, $2, CURRENT_TIMESTAMP, NULL)", queryID, requestQuery.Addres, requestQuery.RequestSearch.Lng, requestQuery.RequestSearch.Lat)
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
func (r *RepositoryImpl) Delete(ctx context.Context, userID int) error {
	_, err := r.DB.ExecContext(ctx, "UPDATE users SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1", userID)
	if err != nil {
		log.Println("Error deleting user:", err)
		return err
	}
	return nil
}

func (r *RepositoryImpl) CacheChecker(ctx context.Context, email string, historyCount int) ([]mod.SearchHistory, error) {
	// ходим в базу поисковых запросов и базу ответов дадата
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return []mod.SearchHistory{}, err
	}
	defer tx.Rollback()
	// получаем id пользователя из базы данных
	userID, err := r.GetByEmail(ctx, email)
	if err != nil {
		return []mod.SearchHistory{}, err
	}
	// получаем последние n поисковых запросов из базы данных
	query := `
        SELECT q.id, q.query, r.address, r.lng, r.lat
        FROM qwery_history q
        LEFT JOIN response_history r ON q.id = r.qwery_key
        WHERE q.user_key = $1
        ORDER BY q.created_date DESC
        LIMIT $2
    `
	rows, err := tx.QueryContext(ctx, query, userID, historyCount)
	if err != nil {
		log.Println("Error querying search history:", err)
		return []mod.SearchHistory{}, err
	}
	defer rows.Close()

	// Сбор результатов
	var searchHistories []mod.SearchHistory
	for rows.Next() {
		var sh mod.SearchHistory
		err := rows.Scan(&sh.ID, &sh.Query, &sh.Address, &sh.Lng, &sh.Lat)
		if err != nil {
			log.Println("Error scanning search history:", err)
			return []mod.SearchHistory{}, err
		}
		searchHistories = append(searchHistories, sh)
	}

	// Подтверждение транзакции
	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction:", err)
		return []mod.SearchHistory{}, err
	}

	return searchHistories, nil
}
