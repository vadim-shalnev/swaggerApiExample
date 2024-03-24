package geocodController

import (
	"context"
	"net/http"
	"strings"
)

type GeocodControllerImpl struct {
	Geocoder  service.Geocoder
	Responder responder.Responder
}

type Geocoder interface {
	HandleSearch(w http.ResponseWriter, r *http.Request)
	HandleGeo(w http.ResponseWriter, r *http.Request)
}

func NewGeocodController(geocoder service.Geocoder, responder responder.Responder) GeocodControllerImpl {
	return GeocodControllerImpl{Geocoder: geocoder, Responder: responder}
}

func (c *GeocodControllerImpl) HandleSearch(w http.ResponseWriter, r *http.Request) {
	Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	ctx := context.WithValue(r.Context(), "jwt_token", Usertoken)
	respSearch, err := c.Service.Search(ctx, r.Body)
	if err != nil {
		handleError(w, err)
		return
	}
	sendJSONResponse(w, respSearch)
}

func (c *GeocodControllerImpl) HandleGeo(w http.ResponseWriter, r *http.Request) {
	Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	ctx := context.WithValue(r.Context(), "jwt_token", Usertoken)
	respGeo, err := c.Service.Address(ctx, r.Body)
	if err != nil {
		handleError(w, err)
		return
	}
	sendJSONResponse(w, respGeo)
}
