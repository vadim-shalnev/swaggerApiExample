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
	CheckToken(ctx context.Context, token string) bool
	RefreshToken(ctx context.Context, email, password, newToken string) error
	CacheChecker(ctx context.Context, query mod.RequestQuery) (bool, mod.RequestAddress, error)
	Select(ctx context.Context, query mod.UserRequest) error
	Insert(ctx context.Context, query mod.RequestUser) error
	Update(ctx context.Context, query mod.RequestUser) error
	Delete(ctx context.Context, query mod.RequestUser) error
	GetAll(ctx context.Context, query mod.RequestUser) ([]mod.RequestAddress, error)
	GetById(ctx context.Context, query mod.RequestUser) (mod.RequestAddress, error)
	GetByEmail(ctx context.Context, query mod.RequestUser) (mod.RequestAddress, error)
	List(ctx context.Context, query mod.RequestUser) ([]mod.RequestAddress, error)
}

func NewRepositoryImpl(db *sql.DB) *RepositoryImpl {
	return &RepositoryImpl{DB: db}
}

func (r *RepositoryImpl) CreateUser(ctx context.Context, user mod.NewUserRequest) error {
	_, err := r.DB.Exec("INSERT INTO users (email, password, token) VALUES ($1, $2, $3)", user.Email, user.Password, user.TokenString.Token)
	if err != nil {
		log.Println("Error creating user:", err)
		return errors.New("failed to create user")
	}
	return nil
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

func (r *RepositoryImpl) CheckToken(ctx context.Context, token string) bool {
	var count int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM users WHERE token = $1", token).Scan(&count)
	if err != nil {
		log.Println("Error checking token:", err)
		return false
	}
	return count > 0
}

func (r *RepositoryImpl) RefreshToken(ctx context.Context, email, password, newToken string) error {
	_, err := r.DB.Exec("UPDATE users SET token = $1 WHERE email = $2 AND password = $3", newToken, email, password)
	if err != nil {
		log.Println("Error refreshing token:", err)
		return errors.New("failed to refresh token")
	}
	return nil
}

func (r *RepositoryImpl) CacheChecker(ctx context.Context, query mod.RequestQuery) (bool, mod.RequestAddress, error) {
	// Implement cache checking logic here
	return true, mod.RequestAddress{}, nil
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

func (r *RepositoryImpl) GetAll(ctx context.Context, query mod.RequestUser) ([]mod.RequestAddress, error) {
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
