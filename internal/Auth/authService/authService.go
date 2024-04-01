package authService

import (
	"errors"
	mod "github.com/vadim-shalnev/swaggerApiExample/Models"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Auth/authRepository"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Cryptografi"
	"github.com/vadim-shalnev/swaggerApiExample/internal/middleware"
	"log"
	"strings"
)

func NewAuthService(repository authRepository.AuthRepository, md middleware.TokenManager) *Authservice {
	return &Authservice{Repo: repository, Tokenmanager: md}
}

func (f *Authservice) Register(regData mod.NewUserRequest) (string, error) {
	tokenAuth, err := f.Tokenmanager.TokenGenerate(regData.Email, regData.Password)
	if err != nil {
		return "", err
	}
	// Хэшируем пароль и добавляем его в запрос к БД
	hashPass, err := Cryptografi.HashPassword(regData)
	if err != nil {
		return "", errors.New("failed to hash password")
	}
	err = f.Repo.CreateUser(hashPass)
	if err != nil {
		return "", errors.New("failed to add new userController to the database")
	}
	return tokenAuth, nil
}

func (f *Authservice) Login(loginData mod.NewUserRequest) (string, error) {
	emailValid, passwordValid := f.UserInfoChecker(loginData.Email, loginData.Password)
	if !emailValid || !passwordValid {
		return "", errors.New("invalid email or password")
	}
	freshToken := f.Tokenmanager.RefreshToken(loginData.Email, loginData.Password)

	return freshToken, nil
}
func (f *Authservice) Logout() (string, error) {
	return "", nil
}

func (f *Authservice) UserInfoChecker(email, password string) (bool, bool) {
	user, _ := f.Repo.GetByEmail(email)

	if strings.TrimSpace(user.Email) != strings.TrimSpace(email) {
		return false, false
	}
	if err := Cryptografi.CheckPassword(user.Password, password); err != nil {
		log.Println("check pass is false")
		return false, false
	}
	return true, true
}
