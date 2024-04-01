package Router

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/vadim-shalnev/swaggerApiExample/Models/controller"
	_ "github.com/vadim-shalnev/swaggerApiExample/docs"
	"net/http"
)

func New_router(controllers controller.Controllers) http.Handler {
	r := chi.NewRouter()
	md := controllers.Auth.Auth.AuthService.Tokenmanager
	r.Post("/api/register", controllers.Auth.Register)
	r.Post("/api/login", controllers.Auth.Login)
	r.Route("/api/user", func(r chi.Router) {
		r.Use(md.AuthMiddleware)
		r.Get("/get/{id}", controllers.User.GetUser)
		r.Get("/list/", controllers.User.ListUsers)
		r.Delete("/del/{id}", controllers.User.DelUser)
	})
	r.Route("/api/address", func(r chi.Router) {
		r.Use(md.AuthMiddleware)
		r.Post("/search", controllers.Geo.HandleSearch)
		r.Post("/geocode", controllers.Geo.HandleGeo)
	})
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	return r
}
