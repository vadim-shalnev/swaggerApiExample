package geocoder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ekomobile/dadata/v2"
	"github.com/ekomobile/dadata/v2/api/model"
	"github.com/ekomobile/dadata/v2/client"
	mod "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Models"
	repository "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Repository"
	"io"
	"io/ioutil"
	"log"
)

const (
	ApiKey    = "22d3fa86b8743e497b32195cbc690abc06b42436"
	SecretKey = "adf07bdd63b240ae60087efd2e72269b9c65cc91"
)

func NewdadataWorkerService(repository repository.Repository) *DadataWorkerImpl {
	return &DadataWorkerImpl{repo: repository}
}

func (d *DadataWorkerImpl) Search(ctx context.Context, userRequest io.ReadCloser) (interface{}, error) {
	bodyJSON, err := ioutil.ReadAll(userRequest)
	if err != nil {
		return mod.RequestQuery{}, errors.New("failed to read request")
	}

	var searchRequest mod.RequestQuery
	var responseQuery mod.RequestQuery

	err = json.Unmarshal(bodyJSON, &searchRequest)
	if err != nil {
		return mod.RequestQuery{}, errors.New("failed to deserialize JSON")
	}

	resp, err := d.HandleWorker(ctx, searchRequest)
	if err != nil {
		return mod.RequestQuery{}, err
	}
	responseQuery.Query = fmt.Sprintf("Latitude: %s, Longitude: %s", resp.RequestSearch.Lng, resp.RequestSearch.Lat)

	return responseQuery, nil
}

func (d *DadataWorkerImpl) Address(ctx context.Context, userRequest io.ReadCloser) (interface{}, error) {
	bodyJSON, err := ioutil.ReadAll(userRequest)
	if err != nil {
		return mod.RequestQuery{}, errors.New("failed to read request")
	}

	var searchRequest mod.RequestQuery
	var responseQuery mod.RequestQuery

	err = json.Unmarshal(bodyJSON, &searchRequest)
	if err != nil {
		return mod.RequestQuery{}, errors.New("failed to deserialize JSON")
	}

	resp, err := d.HandleWorker(ctx, searchRequest)
	if err != nil {
		return mod.RequestQuery{}, err
	}
	responseQuery.Query = fmt.Sprintf("Formatted address: %s", resp.Addres)

	return responseQuery, nil
}

func (d *DadataWorkerImpl) HandleWorker(ctx context.Context, query mod.RequestQuery) (mod.RequestAddress, error) {
	var requestQuery mod.RequestAddress
	ok, cache, err := d.CacheChecker(ctx, query, 5)
	if err != nil {
		log.Println(err)
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
	err = d.repo.Insert(ctx, query)
	if err != nil {
		return mod.RequestAddress{}, errors.New("Select query error")
	}
	return requestQuery, nil
}

func (d *DadataWorkerImpl) CacheChecker(ctx context.Context, query mod.RequestQuery, ttl int) (bool, mod.RequestAddress, error) {
	userToken := ctx.Value("jwt_token").(string)
	email, _ := d.auth.VerifyToken(userToken, "email")
	// получаем id пользователя
	userID, err := d.repo.GetByEmail(ctx, email)
	// идем в репо за последними запросами
	searchHistory, err := d.repo.CacheChecker(ctx, userID, ttl)
}

func (d *DadataWorkerImpl) Geocode(query mod.RequestQuery) ([]*model.Address, error) {
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
