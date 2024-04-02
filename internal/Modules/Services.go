package Modules

import (
	"github.com/vadim-shalnev/swaggerApiExample/internal/Auth/authService"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Geo/geocodService"
	"github.com/vadim-shalnev/swaggerApiExample/internal/User/userService"
	"github.com/vadim-shalnev/swaggerApiExample/internal/infrastructures/components"
)

type Services struct {
	Auth authService.AuthService
	User userService.UserService
	Geo  geocodService.GeocodeService
}

func NewServices(repos *Storages, components *components.Components) *Services {
	return &Services{
		Auth: authService.NewAuthService(repos.Auth, components.TokenManager, components.Hash),
		User: userService.NewUserService(repos.User),
		Geo:  geocodService.NewgeocodeService(repos.Geo, components.TokenManager, components.Conf),
	}
}
