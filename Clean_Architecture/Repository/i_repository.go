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
	CheckEmail(ctx context.Context, email string) bool
	CheckPassword(ctx context.Context, password string) bool
	CacheChecker(ctx context.Context, userID, historyCount int) (mod.RequestAddress, error)
	Select(ctx context.Context, query mod.UserRequest) error
	Insert(ctx context.Context, query mod.RequestUser) error
	Update(ctx context.Context, query mod.RequestUser) error
	Delete(ctx context.Context, query mod.RequestUser) error
	GetAll(ctx context.Context, query mod.RequestUser, colomsCount int) ([]mod.RequestAddress, error)
	GetById(ctx context.Context, query mod.RequestUser) (mod.RequestAddress, error)
	GetByEmail(ctx context.Context, email string) (mod.RequestAddress, error)
	List(ctx context.Context, query mod.RequestUser) ([]mod.RequestAddress, error)
	GetUser(ctx context.Context, id int) (mod.NewUserResponse, error)
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

// получаем пользователя по id из базы данных
func (r *RepositoryImpl) GetUser(ctx context.Context, id int) (mod.NewUserResponse, error) {
	var user mod.NewUserResponse
	var isDel sql.NullTime
	err := r.DB.QueryRow("SELECT email, role, deleted_at FROM users WHERE id = $1", id).Scan(&user.Email, &user.Role, &isDel)
	if err != nil {
		log.Println("Error getting user:", err)
		return mod.NewUserResponse{}, errors.New("failed to get user")
	}
	if isDel.Valid {
		return user, nil
	} else {
		return mod.NewUserResponse{}, errors.New("user not found")
	}
}

func (r *RepositoryImpl) CheckEmail(ctx context.Context, email string) bool {
	var count int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
	if err != nil {
		log.Println("Error checking email:", err)
		return false
	}
	return count > 0
}

func (r *RepositoryImpl) CheckPassword(ctx context.Context, password string) bool {
	// Implement password check logic here if necessary
	return true
}

func (r *RepositoryImpl) CacheChecker(ctx context.Context, query mod.RequestQuery) (mod.RequestAddress, error) {
	// ходим в базу поисковых запросов и базу ответов дадата
	// если есть ответ то используем его и отдаем дальше
	return mod.RequestAddress{}, nil
}

func (r *RepositoryImpl) Select(ctx context.Context, query mod.UserRequest) error {
	// Implement select logic here
	return nil
}
func (r *RepositoryImpl) Insert(ctx context.Context, query mod.RequestUser) error {
	//TODO implement me
	panic("implement me")
}

func (r *RepositoryImpl) Update(ctx context.Context, query mod.RequestUser) error {
	//TODO implement me
	panic("implement me")
}

func (r *RepositoryImpl) Delete(ctx context.Context, query mod.RequestUser) error {
	//TODO implement me
	panic("implement me")
}

// Ходит по базам и собирает указанное количество строк
func (r *RepositoryImpl) GetAll(ctx context.Context, query mod.RequestUser, colomsCount int) ([]mod.RequestAddress, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RepositoryImpl) GetById(ctx context.Context, query mod.RequestUser) (mod.RequestAddress, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RepositoryImpl) GetByEmail(ctx context.Context, query mod.RequestUser) (mod.RequestAddress, error) {
	//TODO implement me
	panic("implement me")
}
func (r *RepositoryImpl) List(ctx context.Context, query mod.RequestUser) ([]mod.RequestAddress, error) {
	//TODO implement me
	panic("implement me")
}
