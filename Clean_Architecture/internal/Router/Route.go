package Router

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/internal/Controller/authController"
	"github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/internal/Controller/geocodController"
	"github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/internal/Controller/userController"
	"net/http"
)

func New_router(controllerAuth *authController.AuthControllerImpl, controllerUser *userController.UserControllerImpl, controllerGeocode *geocodController.GeocodControllerImpl) http.Handler {
	r := chi.NewRouter()
	controller := controllerAuth
	r.Post("/api/register", controller.Register)
	r.Post("/api/login", controller.Login)
	r.Route("/api/user", func(r chi.Router) {
		r.Use(controller.AuthMiddleware)
		r.Get("/get/{id}", controllerUser.GetUser)
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
