package Repository

import (
	"database/sql"
	"errors"
	mod "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Models"
	"log"
)

type RepositoryImpl struct {
	DB *sql.DB
}

type Repository interface {
	CreateUser(user mod.NewUserResponse) error
	CheckEmail(email string) bool
	CheckPassword(password string) bool
	CheckToken(token string) bool
	RefreshToken(email, password, newToken string) error
	CacheChecker(query mod.RequestUser) (bool, mod.RequestAddress, error)
	Select(query mod.RequestUser) error
	Insert(query mod.RequestUser) error
	Update(query mod.RequestUser) error
	Delete(query mod.RequestUser) error
	GetAll(query mod.RequestUser) ([]mod.RequestAddress, error)
	GetById(query mod.RequestUser) (mod.RequestAddress, error)
	GetByEmail(query mod.RequestUser) (mod.RequestAddress, error)
	List(query mod.RequestUser) ([]mod.RequestAddress, error)
}

func NewRepositoryImpl(db *sql.DB) *RepositoryImpl {
	return &RepositoryImpl{DB: db}
}

func (r *RepositoryImpl) CreateUser(user mod.NewUserResponse) error {
	_, err := r.DB.Exec("INSERT INTO users (email, password, token) VALUES ($1, $2, $3)", user.Email, user.Password, user.TokenString.Token)
	if err != nil {
		log.Println("Error creating user:", err)
		return errors.New("failed to create user")
	}
	return nil
}

func (r *RepositoryImpl) CheckEmail(email string) bool {
	var count int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
	if err != nil {
		log.Println("Error checking email:", err)
		return false
	}
	return count > 0
}

func (r *RepositoryImpl) CheckPassword(password string) bool {
	// Implement password check logic here if necessary
	return true
}

func (r *RepositoryImpl) CheckToken(token string) bool {
	var count int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM users WHERE token = $1", token).Scan(&count)
	if err != nil {
		log.Println("Error checking token:", err)
		return false
	}
	return count > 0
}

func (r *RepositoryImpl) RefreshToken(email, password, newToken string) error {
	_, err := r.DB.Exec("UPDATE users SET token = $1 WHERE email = $2 AND password = $3", newToken, email, password)
	if err != nil {
		log.Println("Error refreshing token:", err)
		return errors.New("failed to refresh token")
	}
	return nil
}

func (r *RepositoryImpl) CacheChecker(query mod.RequestUser) (bool, mod.RequestAddress, error) {
	// Implement cache checking logic here
	return true, mod.RequestAddress{}, nil
}

func (r *RepositoryImpl) Select(query mod.RequestUser) error {
	// Implement select logic here
	return nil
}
func (r *RepositoryImpl) Insert(query mod.RequestUser) error {
	//TODO implement me
	panic("implement me")
}

func (r *RepositoryImpl) Update(query mod.RequestAddress) error {
	//TODO implement me
	panic("implement me")
}

func (r *RepositoryImpl) Delete(query mod.RequestAddress) error {
	//TODO implement me
	panic("implement me")
}

func (r *RepositoryImpl) GetAll(query mod.RequestAddress) ([]mod.RequestAddress, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RepositoryImpl) GetById(query mod.RequestAddress) (mod.RequestAddress, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RepositoryImpl) GetByEmail(query mod.RequestAddress) (mod.RequestAddress, error) {
	//TODO implement me
	panic("implement me")
}
