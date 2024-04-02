package main

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/vadim-shalnev/swaggerApiExample/config"
	"github.com/vadim-shalnev/swaggerApiExample/internal/infrastructures/logs"
	"github.com/vadim-shalnev/swaggerApiExample/run"
	"os"
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

	// Создаем конфигурацию приложения
	conf := config.NewAppConf()

	logger := logs.NewLogger(conf, os.Stdout)
	if err != nil {
		logger.Fatal("error loading .env file")
	}

	app := run.NewApp(conf, logger)

	boot, component := app.Boostrap()

	app.Run(boot, component)

}
