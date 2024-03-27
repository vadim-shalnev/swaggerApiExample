package Repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/vadim-shalnev/swaggerApiExample/Models"
	"log"
)

type RepositoryDB struct {
	DB *sql.DB
}

type Repository interface {
	CreateUser(ctx context.Context, user Models.NewUserRequest) error
	GetByEmail(ctx context.Context, email string) (Models.User, error)
	GetByID(ctx context.Context, id int) (Models.User, error)
	Insert(ctx context.Context, email string, query Models.RequestQuery, requestQuery Models.RequestAddress) error
	Delete(ctx context.Context, userID int) error
	CacheChecker(ctx context.Context, email string, historyCount int) ([]Models.SearchHistory, error)
	List(ctx context.Context) ([]Models.User, error)
	//Select(ctx context.Context, query mod.UserRequest) error
	//Update(ctx context.Context, query mod.RequestUser) error
	//CheckEmail(ctx context.Context, email string) bool
	//CheckPassword(ctx context.Context, password string) bool
	//GetAll(ctx context.Context, query mod.RequestUser, colomsCount int) ([]mod.RequestAddress, error)
}

// Таблица ответов dаdata содержит внешний ключ к таблице запросов, которая в свою очередь хранит внешний ключ таблицы пользователей
func NewRepositoryImpl(db *sql.DB) *RepositoryDB {
	return &RepositoryDB{DB: db}
}

func (r *RepositoryDB) CreateUser(ctx context.Context, user Models.NewUserRequest) error {
	_, err := r.DB.Exec("INSERT INTO users (email, password, role, created_at, deleted_at) VALUES ($1, $2, $3, CURRENT_TIMESTAMP, NULL)", user.Email, user.Password, user.Role)
	if err != nil {
		log.Println("Error creating User:", err)
		return errors.New("failed to create User")
	}
	return nil
}
func (r *RepositoryDB) GetByEmail(ctx context.Context, email string) (Models.User, error) {
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

func (r *RepositoryDB) GetByID(ctx context.Context, id int) (Models.User, error) {
	var user Models.User

	err := r.DB.QueryRow("SELECT id, email, password, role, created_at, deleted_at FROM users WHERE id = $1", id).Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return Models.User{}, errors.New("User not found")
		}
		log.Println("Error getting User by email:", err)
		return Models.User{}, err
	}
	if user.DeletedAt != nil {
		return Models.User{}, errors.New("User not found")
	}
	return user, nil
}
func (r *RepositoryDB) Insert(ctx context.Context, email string, query Models.RequestQuery, requestQuery Models.RequestAddress) error {
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
func (r *RepositoryDB) Delete(ctx context.Context, userID int) error {
	_, err := r.DB.ExecContext(ctx, "UPDATE users SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1", userID)
	if err != nil {
		log.Println("Error deleting User:", err)
		return err
	}
	return nil
}

func (r *RepositoryDB) CacheChecker(ctx context.Context, email string, historyCount int) ([]Models.SearchHistory, error) {
	// ходим в базу поисковых запросов и базу ответов дадата
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return []Models.SearchHistory{}, err
	}
	defer tx.Rollback()
	// получаем id пользователя из базы данных
	user, err := r.GetByEmail(ctx, email)
	if err != nil {
		return []Models.SearchHistory{}, err
	}
	userID := user.ID
	// получаем последние n поисковых запросов из базы данных
	query := `
        SELECT q.id, q.query, r.address, r.lng, r.lat
        FROM qwery_history q
        LEFT JOIN response_history r ON q.id = r.qwery_key
        WHERE q.user_key = $1
        ORDER BY q.created_at DESC
        LIMIT $2
    `
	rows, err := tx.QueryContext(ctx, query, userID, historyCount)
	if err != nil {
		log.Println("Error querying search history:", err)
		return []Models.SearchHistory{}, err
	}
	defer rows.Close()

	// Сбор результатов
	var searchHistories []Models.SearchHistory
	for rows.Next() {
		var sh Models.SearchHistory
		err := rows.Scan(&sh.ID, &sh.Query, &sh.Address, &sh.Lng, &sh.Lat)
		if err != nil {
			log.Println("Error scanning search history:", err)
			return []Models.SearchHistory{}, err
		}
		searchHistories = append(searchHistories, sh)
	}

	// Подтверждение транзакции
	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction:", err)
		return []Models.SearchHistory{}, err
	}

	return searchHistories, nil
}
func (r *RepositoryDB) List(ctx context.Context) ([]Models.User, error) {
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
