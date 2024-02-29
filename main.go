package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ekomobile/dadata/v2"
	"github.com/ekomobile/dadata/v2/api/model"
	"github.com/ekomobile/dadata/v2/client"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github/vadim-shalnev/docs"
	"io/ioutil"
	"net/http"
)

const ApiKey string = "22d3fa86b8743e497b32195cbc690abc06b42436"
const SecretKey string = "adf07bdd63b240ae60087efd2e72269b9c65cc91"

type RequestAddressGeocode struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}
type RequestAddressInfo struct {
	Addres string `json:"addres"`
}
type RequestAddressSearch struct {
	Query string `json:"query"`
}

// @title Todo geocode API
// @version 1.0
// @description API Server for search GEOinfo

// @host localhost:8080
// @BasePath /api/address

func main() {
	r := chi.NewRouter()
	Query(r)
	http.ListenAndServe(":8080", r)

}

func Query(r *chi.Mux) {
	r.Route("/api", func(r chi.Router) {
		r.Route("/address", func(r chi.Router) {
			r.Post("/search", HandleSearch)
			r.Post("/geocode", HandleGeo)

		})

	})
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // Укажите путь к файлу swagger.json
	))

}

// HandleSearch @HandleSearch
// @Summary QueryGeocode
// @Tags geocode
// @Description create a search query
// @Accept json
// @Produce json
// @Param input body RequestAddressSearch true "query"
// @Success 200 {integer} integer 1
// @Failure 404 {error} http.Error
// @Failure 500 {error} http.Error
// @Router /geocode [post]
func HandleSearch(w http.ResponseWriter, r *http.Request) {

	bodyJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var SearchResp RequestAddressSearch

	err = json.Unmarshal(bodyJSON, &SearchResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	geocodeResponse, err := Geocode(SearchResp)
	if err != nil {
		fmt.Println(err)
	}

	var clientResponse RequestAddressSearch
	var typeResponseInfo RequestAddressInfo
	var typeResponseGeocode RequestAddressGeocode

	url := r.URL

	if url.Path == "/api/address/search" {
		for _, v := range geocodeResponse {
			typeResponseInfo.Addres = v.Result
			clientResponse.Query = v.Result
		}
	}
	if url.Path == "/api/address/geocode" {
		for _, v := range geocodeResponse {
			typeResponseGeocode.Lat = v.GeoLat
			typeResponseGeocode.Lng = v.GeoLon
			clientResponse.Query += fmt.Sprintf("Широта: %s Долгота %s", v.GeoLat, v.GeoLon)
		}
	}

	bodyRespone, err := json.Marshal(clientResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Write(bodyRespone)
}

// HandleGeo @HandleGeo
// @Summary QuerySearch
// @Tags search
// @Description create a search query
// @Accept json
// @Produce json
// @Param input body RequestAddressSearch true "query"
// @Success 200 {integer} integer 1
// @Failure 404 {error} http.Error
// @Failure 500 {error} http.Error
// @Router /search [post]
func HandleGeo(w http.ResponseWriter, r *http.Request) {

	bodyJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var SearchResp RequestAddressSearch

	err = json.Unmarshal(bodyJSON, &SearchResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	geocodeResponse, err := Geocode(SearchResp)
	if err != nil {
		fmt.Println(err)
	}

	var clientResponse RequestAddressSearch
	var typeResponseInfo RequestAddressInfo
	var typeResponseGeocode RequestAddressGeocode

	url := r.URL

	if url.Path == "/api/address/search" {
		for _, v := range geocodeResponse {
			typeResponseInfo.Addres = v.Result
			clientResponse.Query = v.Result
		}
	}
	if url.Path == "/api/address/geocode" {
		for _, v := range geocodeResponse {
			typeResponseGeocode.Lat = v.GeoLat
			typeResponseGeocode.Lng = v.GeoLon
			clientResponse.Query += fmt.Sprintf("Широта: %s Долгота %s", v.GeoLat, v.GeoLon)
		}
	}

	bodyRespone, err := json.Marshal(clientResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Write(bodyRespone)
}

func Geocode(Querys RequestAddressSearch) ([]*model.Address, error) {

	creds := client.Credentials{
		ApiKeyValue:    ApiKey,
		SecretKeyValue: SecretKey,
	}

	api := dadata.NewCleanApi(client.WithCredentialProvider(&creds))

	result, err := api.Address(context.Background(), Querys.Query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return result, nil

}
