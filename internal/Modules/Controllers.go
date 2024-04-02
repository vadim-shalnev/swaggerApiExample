package Modules

import (
	"github.com/vadim-shalnev/swaggerApiExample/internal/Auth/authController"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Geo/geocodController"
	"github.com/vadim-shalnev/swaggerApiExample/internal/User/userController"
	"github.com/vadim-shalnev/swaggerApiExample/internal/infrastructures/components"
)

type Controllers struct {
	Auth authController.AuthController
	User userController.UserController
	Geo  geocodController.Geocoder
}

func NewControllers(services *Services, components *components.Components) *Controllers {
	return &Controllers{
		Auth: authController.NewAuthController(services.Auth, components.Responder),
		User: userController.NewUserController(services.User, components.Responder),
		Geo:  geocodController.NewGeocodController(services.Geo, components.Responder),
	}
}
