package run

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/vadim-shalnev/swaggerApiExample/config"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Modules"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Router"
	"github.com/vadim-shalnev/swaggerApiExample/internal/infrastructures/Cache"
	"github.com/vadim-shalnev/swaggerApiExample/internal/infrastructures/Cryptografi"
	"github.com/vadim-shalnev/swaggerApiExample/internal/infrastructures/Responder"
	"github.com/vadim-shalnev/swaggerApiExample/internal/infrastructures/components"
	"github.com/vadim-shalnev/swaggerApiExample/internal/infrastructures/middleware"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type App struct {
	Conf   config.AppConf
	DB     *sql.DB
	Logger *zap.Logger
}

type Apper interface {
	Run(controllers *Modules.Controllers, component components.Components)
	Boostrap()
}

func NewApp(conf config.AppConf, logger *zap.Logger) *App {
	return &App{
		Conf:   conf,
		Logger: logger,
	}
}
func (a *App) Run(controllers *Modules.Controllers, component *components.Components) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Создаем контекст с таймаутом 5 секунд
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Создаем HTTP-сервер
	server := &http.Server{
		Addr:    ":8080",
		Handler: Router.New_router(controllers, component),
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Ожидаем сигнала завершения от операционной системы или отмены контекста
	select {
	case <-stop:
		fmt.Println("Получен сигнал завершения")
	case <-ctx.Done():
		fmt.Println("Завершение по таймауту")
	}

	// Выполняем graceful shutdown сервера
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Ошибка при graceful shutdown:", err)
	}
}

func (a *App) Boostrap() (*Modules.Controllers, *components.Components) {
	// Создаем менеджер токенов
	tokenManager := middleware.NewTokenManager(a.Conf.MD)
	// Создаем обработчики ответов и ошибок
	respondManager := Responder.NewResponder(a.Logger)
	// Создаем хэш
	hash := Cryptografi.NewHasher()
	// Создаем кэш
	cache := Cache.NewRedisCache(a.Conf.Cache.Address, a.Conf.Cache.Password, a.Logger)
	// Создаем компоненты для роутера
	newComponents := components.NewComponents(a.Conf, tokenManager, hash, cache, respondManager, a.Logger)
	// Подключаемся к PostgreSQL
	db := a.ConnectionDB(a.Conf)
	// Создаем таблицы если их нет
	a.CreateTablesIfNotExist(db)
	// Создаем слои
	repos := Modules.NewStorages(db, cache)
	services := Modules.NewServices(repos, newComponents)
	controllers := Modules.NewControllers(services, newComponents)

	return controllers, newComponents
}
