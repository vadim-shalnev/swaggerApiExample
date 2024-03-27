package geocodController

import (
	"context"
	"encoding/json"
	"github.com/vadim-shalnev/swaggerApiExample/Models"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Geocoder/geocodService"
	responder "github.com/vadim-shalnev/swaggerApiExample/internal/Responder"
	"io/ioutil"
	"net/http"
	"strings"
)

type Geocodcontroller struct {
	Geocoder geocodService.GeocodeWorker
}

type Geocoder interface {
	HandleSearch(w http.ResponseWriter, r *http.Request)
	HandleGeo(w http.ResponseWriter, r *http.Request)
}

func NewGeocodController(geocoder geocodService.GeocodeWorker) *Geocodcontroller {
	return &Geocodcontroller{Geocoder: geocoder}
}

// HandleSearch @Summary Поиск полного адреса
// @Description Поиск полного адреса
// @Tags geocode
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer"
// @Param request body Models.RequestQuery true "request"
// @Success 200 {object} Models.RequestQuery
// @Router /api/address/search [post]
func (c *Geocodcontroller) HandleSearch(w http.ResponseWriter, r *http.Request) {
	Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	ctx := context.WithValue(r.Context(), "jwt_token", Usertoken)
	bodyJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responder.HandleError(w, err)
		return
	}

	var searchRequest Models.RequestQuery

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

// HandleGeo @Summary Поиск координат по адресу
// @Description Поиск координат по адресу
// @Tags geocode
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer"
// @Param request body Models.RequestQuery true "request"
// @Success 200 {object} Models.RequestQuery
// @Router /api/address/geocode [post]
func (c *Geocodcontroller) HandleGeo(w http.ResponseWriter, r *http.Request) {
	Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	ctx := context.WithValue(r.Context(), "jwt_token", Usertoken)
	bodyJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responder.HandleError(w, err)
		return
	}

	var searchRequest Models.RequestQuery

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
