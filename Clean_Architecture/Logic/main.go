package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ekomobile/dadata/v2"
	"github.com/ekomobile/dadata/v2/api/model"
	"github.com/ekomobile/dadata/v2/client"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io/ioutil"
	"net/http"
)

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

const ApiKey string = "22d3fa86b8743e497b32195cbc690abc06b42436"
const SecretKey string = "adf07bdd63b240ae60087efd2e72269b9c65cc91"

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Route("/api/address", func(r chi.Router) {
		r.Post("/search", Handle)
		r.Post("/geocode", Handle)
	})
	http.ListenAndServe(":8090", r)
}

func Handle(w http.ResponseWriter, r *http.Request) {

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
