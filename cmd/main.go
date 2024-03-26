package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Controller/authController"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Controller/geocodController"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Controller/userController"
	repository "github.com/vadim-shalnev/swaggerApiExample/internal/Repository"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Router"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Service/authService"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Service/geocodService"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Service/userService"
	"log"
	"net/http"
)

// @title Geocode Api Server
// @version 1.0
// @description This is a geocode api server.

// @host localhost:8080
// @BasePath /api/
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=vubuntu password=qwerty dbname=vubuntu sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	createTablesIfNotExist(db)

	repo := repository.NewRepositoryImpl(db)
	servauth := authService.NewAuthService(repo)
	serveuser := userService.NewAuthService(repo)
	servegeocode := geocodService.NewgeocodeService(repo, servauth)
	cAuth := authController.NewAuthController(servauth)
	cUser := userController.NewUserController(serveuser)
	cGeo := geocodController.NewGeocodController(servegeocode)
	router := Router.New_router(cAuth, cUser, cGeo)

	http.ListenAndServe(":8080", router)

}
