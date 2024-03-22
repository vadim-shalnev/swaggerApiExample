package Route

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Controller"
	"net/http"
)

func New_router(controllers *Controller.Controller) http.Handler {
	r := chi.NewRouter()
	controller := controllers.Auth
	r.Get("/api/login", controller.Login)
	r.Post("/api/register", controller.Register)
	r.Route("/api/address", func(r chi.Router) {
		r.Use(controller.AuthMiddleware)
		r.Post("/search", controller.HandleSearch)
		r.Post("/geocode", controller.HandleGeo)
	})
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // Укажите путь к файлу swagger.json
	))
	http.ListenAndServe(":8080", r)
	return r
}
