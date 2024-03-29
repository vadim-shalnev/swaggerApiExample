package Router

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/vadim-shalnev/swaggerApiExample/docs"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Auth/authController"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Geocoder/geocodController"
	"github.com/vadim-shalnev/swaggerApiExample/internal/User/userController"
	"net/http"
)

func New_router(controllerAuth *authController.Authcontroller, controllerUser *userController.Usercontroller, controllerGeocode *geocodController.Geocodcontroller) http.Handler {
	r := chi.NewRouter()
	controller := controllerAuth
	r.Post("/api/register", controller.Register)
	r.Post("/api/login", controller.Login)
	r.Route("/api/user", func(r chi.Router) {
		r.Use(controller.AuthMiddleware)
		r.Get("/get/{id}", controllerUser.GetUser)
		r.Get("/list/", controllerUser.ListUsers)
		r.Delete("/del/{id}", controllerUser.DelUser)
	})
	r.Route("/api/address", func(r chi.Router) {
		r.Use(controller.AuthMiddleware)
		r.Post("/search", controllerGeocode.HandleSearch)
		r.Post("/geocode", controllerGeocode.HandleGeo)
	})
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	return r
}
