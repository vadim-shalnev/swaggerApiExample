package geocodService

import (
	"context"
	"errors"
	"fmt"
	"github.com/ekomobile/dadata/v2"
	"github.com/ekomobile/dadata/v2/api/model"
	"github.com/ekomobile/dadata/v2/client"
	mod "github.com/vadim-shalnev/swaggerApiExample/Models"
	"github.com/vadim-shalnev/swaggerApiExample/config"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Geo/geocodeRepository"
	"github.com/vadim-shalnev/swaggerApiExample/internal/infrastructures/middleware"
	"log"
	"time"
)

/*
const (
	ApiKey    = "22d3fa86b8743e497b32195cbc690abc06b42436"
	SecretKey = "adf07bdd63b240ae60087efd2e72269b9c65cc91"
)
*/

func NewgeocodeService(repository geocodeRepository.GeocodeRepository, Tokenmanager middleware.TokenManager, conf config.AppConf) *Geocodeservice {
	return &Geocodeservice{repo: repository, Tokenmanager: Tokenmanager, Config: conf}
}

func (d *Geocodeservice) Search(ctx context.Context, userRequest mod.RequestQuery) (mod.RequestQuery, error) {
	var responseQuery mod.RequestQuery
	resp, err := d.HandleWorker(ctx, userRequest)
	if err != nil {
		return mod.RequestQuery{}, err
	}
	responseQuery.Query = fmt.Sprintf("Latitude: %s, Longitude: %s", resp.RequestSearch.Lng, resp.RequestSearch.Lat)

	return responseQuery, nil
}

func (d *Geocodeservice) Address(ctx context.Context, userRequest mod.RequestQuery) (mod.RequestQuery, error) {
	var responseQuery mod.RequestQuery
	resp, err := d.HandleWorker(ctx, userRequest)
	if err != nil {
		return mod.RequestQuery{}, err
	}
	responseQuery.Query = fmt.Sprintf("Formatted address: %s", resp.Addres)

	return responseQuery, nil
}

func (d *Geocodeservice) HandleWorker(ctx context.Context, query mod.RequestQuery) (mod.RequestAddress, error) {
	var requestQuery mod.RequestAddress
	userToken := ctx.Value("jwt_token").(string)
	claims, err := d.Tokenmanager.GetClaims(userToken)
	if err != nil {
		return mod.RequestAddress{}, err
	}
	email := claims["Email"].(string)
	timeout, cancel := context.WithTimeout(ctx, time.Duration(500)*time.Millisecond)
	defer cancel()
	cache, err := d.repo.CacheChecker(timeout, query)
	if err == nil {
		log.Println("Запрос из кэша")
		return cache, nil
	}

	geocodeResponse, err := d.Geocode(query)
	if err != nil {
		return mod.RequestAddress{}, errors.New("error in dadata operation")
	}
	for _, v := range geocodeResponse {
		requestQuery.RequestSearch.Lat = v.GeoLat
		requestQuery.RequestSearch.Lng = v.GeoLon
		requestQuery.Addres = v.Result
	}
	err = d.repo.Insert(ctx, email, query, requestQuery)
	if err != nil {
		log.Println("ошибка записи в репо", err)
		return mod.RequestAddress{}, errors.New("Select query error")
	}
	return requestQuery, nil
}

/*
	searchHistory, err := d.repo.CacheChecker(ctx, email, ttl)
	if err != nil {
		return false, mod.RequestAddress{}, "", err
	}
	if searchHistory != nil || len(searchHistory) > 0 {
		levenshtein, ok := Cryptografi.Levanshtain(searchHistory, query.Query)
		if ok {
			return true, levenshtein, email, nil
		}
	}
	return false, mod.RequestAddress{}, email, nil

*/

func (d *Geocodeservice) Geocode(query mod.RequestQuery) ([]*model.Address, error) {
	creds := client.Credentials{
		ApiKeyValue:    d.Config.GEO.APIKey,
		SecretKeyValue: d.Config.GEO.GEOKey,
	}
	api := dadata.NewCleanApi(client.WithCredentialProvider(&creds))
	result, err := api.Address(context.Background(), query.Query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, nil
}
