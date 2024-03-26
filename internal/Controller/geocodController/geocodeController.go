package geocodController

import (
	"context"
	"encoding/json"
	mod "github.com/vadim-shalnev/swaggerApiExample/Models"
	responder "github.com/vadim-shalnev/swaggerApiExample/internal/Responder"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Service/geocodService"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type GeocodControllerImpl struct {
	Geocoder geocodService.GeocodeWorker
}

type Geocoder interface {
	HandleSearch(w http.ResponseWriter, r *http.Request)
	HandleGeo(w http.ResponseWriter, r *http.Request)
}

func NewGeocodController(geocoder geocodService.GeocodeWorker) *GeocodControllerImpl {
	return &GeocodControllerImpl{Geocoder: geocoder}
}

func (c *GeocodControllerImpl) HandleSearch(w http.ResponseWriter, r *http.Request) {
	Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	log.Println("searchToken", Usertoken)
	ctx := context.WithValue(r.Context(), "jwt_token", Usertoken)
	bodyJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responder.HandleError(w, err)
		return
	}

	var searchRequest mod.RequestQuery

	err = json.Unmarshal(bodyJSON, &searchRequest)
	if err != nil {
		responder.HandleError(w, err)
		return
	}

	respSearch, err := c.Geocoder.Search(ctx, searchRequest)
	if err != nil {
		responder.HandleError(w, err)
		return
	}
	responder.SendJSONResponse(w, respSearch)
}

func (c *GeocodControllerImpl) HandleGeo(w http.ResponseWriter, r *http.Request) {
	Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	ctx := context.WithValue(r.Context(), "jwt_token", Usertoken)
	bodyJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responder.HandleError(w, err)
		return
	}

	var searchRequest mod.RequestQuery

	err = json.Unmarshal(bodyJSON, &searchRequest)
	if err != nil {
		responder.HandleError(w, err)
		return
	}
	respGeo, err := c.Geocoder.Address(ctx, searchRequest)
	if err != nil {
		responder.HandleError(w, err)
		return
	}
	responder.SendJSONResponse(w, respGeo)
}
