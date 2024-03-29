package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Auth/authController"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Auth/authRepository"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Auth/authService"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Geocoder/geocodController"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Geocoder/geocodService"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Geocoder/geocodeRepository"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Router"
	"github.com/vadim-shalnev/swaggerApiExample/internal/User/userController"
	"github.com/vadim-shalnev/swaggerApiExample/internal/User/userRepository"
	"github.com/vadim-shalnev/swaggerApiExample/internal/User/userService"
	"log"
	"net/http"
	"time"
)

// @title Swagger Example API
// @version 1.0
// @description This is a geocode api server.

// @host localhost:8080
// @BasePath /api/
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	time.Sleep(time.Second * 5)

	db, err := sql.Open("postgres", "host=db port=5432 user=postgresuser password=userpassword dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	createTablesIfNotExist(db)

	repoauth := authRepository.NewAuthrepository(db)
	repouser := userRepository.NewUserRepository(db)
	repogeo := geocodeRepository.NewGeocodeRepository(db)
	servauth := authService.NewAuthService(repoauth)
	serveuser := userService.NewAuthService(repouser)
	servegeocode := geocodService.NewgeocodeService(repogeo, servauth)
	cAuth := authController.NewAuthController(servauth)
	cUser := userController.NewUserController(serveuser)
	cGeo := geocodController.NewGeocodController(servegeocode)
	router := Router.New_router(cAuth, cUser, cGeo)

	http.ListenAndServe(":8080", router)

}
