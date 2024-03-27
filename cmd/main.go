package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Auth/authController"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Auth/authService"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Geocoder/geocodController"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Geocoder/geocodService"
	repository "github.com/vadim-shalnev/swaggerApiExample/internal/Repository"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Router"
	"github.com/vadim-shalnev/swaggerApiExample/internal/User/userController"
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
	// добавил задержку т.к. постгрес в контейнере шалил и не хотел запускаться
	time.Sleep(time.Second * 5)
	db, err := sql.Open("postgres", "host=db port=5432 User=postgresuser password=userpassword dbname=postgres sslmode=disable")
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
