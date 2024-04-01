package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/vadim-shalnev/swaggerApiExample/Models/controller"
	"github.com/vadim-shalnev/swaggerApiExample/config"
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
	"github.com/vadim-shalnev/swaggerApiExample/internal/middleware"
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
	// Загружаем переменные окружения из файла .env
	err := godotenv.Load("/app/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}
	// Создаем конфигурацию приложения
	conf := config.NewAppConf()

	db := ConnectionDB(conf)
	createTablesIfNotExist(db)

	st := HandleInit(db, conf)
	router := Router.New_router(st)

	defer db.Close()
	http.ListenAndServe(":8080", router)

}

// ConnectionDB Подключаемся к бд
func ConnectionDB(conf config.AppConf) *sql.DB {
	time.Sleep(time.Second * 2)
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", conf.DB.Host, conf.DB.Port, conf.DB.User, conf.DB.Password, conf.DB.Name))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	return db
}

func HandleInit(db *sql.DB, conf config.AppConf) controller.Controllers {
	// DB
	repoauth := authRepository.NewAuthrepository(db)
	repouser := userRepository.NewUserRepository(db)
	repogeo := geocodeRepository.NewGeocodeRepository(db)

	midleware := middleware.NewTokenManager(conf.MD)
	servauth := authService.NewAuthService(repoauth, midleware)
	facade := authService.NewAuthFacade(servauth)
	serveuser := userService.NewAuthService(repouser)
	servegeocode := geocodService.NewgeocodeService(repogeo, servauth)

	cAuth := authController.NewAuthController(*facade)
	cUser := userController.NewUserController(serveuser)
	cGeo := geocodController.NewGeocodController(servegeocode)
	return controller.Controllers{
		Auth: cAuth,
		User: cUser,
		Geo:  cGeo,
	}
}
