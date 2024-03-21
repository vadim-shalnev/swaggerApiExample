package Route

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Controller"
	"net/http"
)

func New_route() http.Handler {
	r := chi.NewRouter()
	r.Get("/api/login", Controller.Login)
	r.Post("/api/register", Controller.Register)
	r.Route("/api/address", func(r chi.Router) {
		r.Use(Controller.AuthMiddleware)
		r.Post("/search", Controller.HandleSearch)
		r.Post("/geocode", Controller.HandleGeo)
	})
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // Укажите путь к файлу swagger.json
	))
	http.ListenAndServe(":8080", r)
	return r
}
