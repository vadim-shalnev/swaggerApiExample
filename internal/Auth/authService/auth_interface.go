package authService

import (
	mod "github.com/vadim-shalnev/swaggerApiExample/Models"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Auth/authRepository"
	"github.com/vadim-shalnev/swaggerApiExample/internal/infrastructures/Cryptografi"
	"github.com/vadim-shalnev/swaggerApiExample/internal/infrastructures/middleware"
)

type Authservice struct {
	Repo         authRepository.AuthRepository
	Tokenmanager middleware.TokenManager
	Hash         Cryptografi.Hasher
}

type AuthService interface {
	Register(regData mod.NewUserRequest) (string, error)
	Login(loginData mod.NewUserRequest) (string, error)
	Logout() (string, error)
	UserInfoChecker(email, password string) (bool, bool)
}

type Facade struct {
	AuthService *Authservice
}

func NewAuthFacade(auth *Authservice) *Facade {
	return &Facade{AuthService: auth}
}
func (f Facade) Register(regData mod.NewUserRequest) (string, error) {
	return f.AuthService.Register(regData)
}
func (f Facade) Login(loginData mod.NewUserRequest) (string, error) {
	return f.AuthService.Login(loginData)
}
func (f Facade) Logout() (string, error) {
	return f.AuthService.Logout()
}
func (f Facade) UserInfoChecker(email, password string) (bool, bool) {
	return f.AuthService.UserInfoChecker(email, password)
}
