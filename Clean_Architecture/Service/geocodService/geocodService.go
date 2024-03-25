package geocodService

import (
	"context"
	"errors"
	"fmt"
	"github.com/ekomobile/dadata/v2"
	"github.com/ekomobile/dadata/v2/api/model"
	"github.com/ekomobile/dadata/v2/client"
	"github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Cryptografi"
	mod "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Models"
	repository "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Repository"
	"github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Service/authService"
	"log"
)

const (
	ApiKey    = "22d3fa86b8743e497b32195cbc690abc06b42436"
	SecretKey = "adf07bdd63b240ae60087efd2e72269b9c65cc91"
)

func NewgeocodeService(repository repository.Repository, aothorisation authService.AuthService) *GeocodeWorkerImpl {
	return &GeocodeWorkerImpl{repo: repository, auth: aothorisation}
}

func (d *GeocodeWorkerImpl) Search(ctx context.Context, userRequest mod.RequestQuery) (interface{}, error) {
	var responseQuery mod.RequestQuery
	resp, err := d.HandleWorker(ctx, userRequest)
	if err != nil {
		return mod.RequestQuery{}, err
	}
	responseQuery.Query = fmt.Sprintf("Latitude: %s, Longitude: %s", resp.RequestSearch.Lng, resp.RequestSearch.Lat)

	return responseQuery, nil
}

func (d *GeocodeWorkerImpl) Address(ctx context.Context, userRequest mod.RequestQuery) (interface{}, error) {
	var responseQuery mod.RequestQuery
	resp, err := d.HandleWorker(ctx, userRequest)
	if err != nil {
		return mod.RequestQuery{}, err
	}
	responseQuery.Query = fmt.Sprintf("Formatted address: %s", resp.Addres)

	return responseQuery, nil
}

func (d *GeocodeWorkerImpl) HandleWorker(ctx context.Context, query mod.RequestQuery) (mod.RequestAddress, error) {
	var requestQuery mod.RequestAddress
	ok, cache, email, err := d.CacheChecker(ctx, query, 5)
	if err != nil {
		log.Println("ошибка проверки кэша", err)
	}
	if ok {
		requestQuery.Addres = cache.Addres
		requestQuery.RequestSearch.Lat = cache.RequestSearch.Lat
		requestQuery.RequestSearch.Lng = cache.RequestSearch.Lng
		return requestQuery, nil
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
		return mod.RequestAddress{}, errors.New("Select query error")
	}
	return requestQuery, nil
}

func (d *GeocodeWorkerImpl) CacheChecker(ctx context.Context, query mod.RequestQuery, ttl int) (bool, mod.RequestAddress, string, error) {
	userToken := ctx.Value("jwt_token").(string)
	email, _, _ := d.auth.VerifyToken(userToken)
	// идем в репо за последними запросами
	searchHistory, err := d.repo.CacheChecker(ctx, email, ttl)
	if err != nil {
		return false, mod.RequestAddress{}, "", err
	}
	levanshtain, ok := Cryptografi.Levanshtain(searchHistory, query.Query)
	if ok {
		return true, levanshtain, email, nil
	}
	return false, mod.RequestAddress{}, "", nil
}

func (d *GeocodeWorkerImpl) Geocode(query mod.RequestQuery) ([]*model.Address, error) {
	creds := client.Credentials{
		ApiKeyValue:    ApiKey,
		SecretKeyValue: SecretKey,
	}
	api := dadata.NewCleanApi(client.WithCredentialProvider(&creds))
	result, err := api.Address(context.Background(), query.Query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, nil
}
