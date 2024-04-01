package controller

import (
	"github.com/vadim-shalnev/swaggerApiExample/internal/Auth/authController"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Geocoder/geocodController"
	"github.com/vadim-shalnev/swaggerApiExample/internal/User/userController"
)

type Controllers struct {
	Auth *authController.Authcontroller
	User *userController.Usercontroller
	Geo  *geocodController.Geocodcontroller
}
