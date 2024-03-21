package Service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ekomobile/dadata/v2"
	"github.com/ekomobile/dadata/v2/api/model"
	"github.com/ekomobile/dadata/v2/client"
	"github.com/go-chi/jwtauth/v5"
	mod "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Models"
	"io"
	"io/ioutil"
	"log"
)

const (
	ApiKey    = "someapi"
	SecretKey = "somekey"
)

type UserServiceImpl struct {
	repo repository.Repository
}

func NewUserServiceImpl(repo repository.Repository) *UserServiceImpl {
	return &UserServiceImpl{repo: repo}
}

func (s *UserServiceImpl) Register(userBody io.ReadCloser) (mod.NewUserResponse, error) {
	var regData mod.NewUserResponse

	bodyJSON, err := ioutil.ReadAll(userBody)
	if err != nil {
		return mod.NewUserResponse{}, errors.New("failed to read request")
	}
	err = json.Unmarshal(bodyJSON, &regData)
	if err != nil {
		return mod.NewUserResponse{}, errors.New("failed to deserialize JSON")
	}

	tokenAuth, err := TokenGenerate(regData.Email, regData.Password)
	if err != nil {
		return mod.NewUserResponse{}, err
	}
	regData.TokenString.Token = tokenAuth

	err = s.repo.CreateUser(regData)
	if err != nil {
		return mod.NewUserResponse{}, errors.New("failed to add new user to the database")
	}

	return regData, nil
}

func (s *UserServiceImpl) Login(userToken string, userBody io.ReadCloser) (mod.NewUserResponse, error) {
	var regData mod.NewUserResponse

	bodyJSON, err := ioutil.ReadAll(userBody)
	if err != nil {
		return mod.NewUserResponse{}, errors.New("failed to read request")
	}
	err = json.Unmarshal(bodyJSON, &regData)
	if err != nil {
		return mod.NewUserResponse{}, errors.New("failed to deserialize JSON")
	}

	emailValid, passwordValid, tokenValid := s.UserInfoChecker(regData.Email, regData.Password, userToken)
	if !emailValid {
		return mod.NewUserResponse{}, errors.New("invalid email")
	}
	if !passwordValid {
		return mod.NewUserResponse{}, errors.New("invalid password")
	}
	if !tokenValid {
		freshToken := RefreshToken(regData.Email, regData.Password)
		regData.TokenString.Token = freshToken
		return mod.NewUserResponse{}, errors.New("you have successfully logged out of the service")
	}

	return regData, nil
}

func (s *UserServiceImpl) UserInfoChecker(email, password, token string) (bool, bool, bool) {
	emailValid := s.repo.CheckEmail(email)
	passwordValid := s.repo.CheckPassword(password)
	tokenValid := s.repo.CheckToken(token)
	return emailValid, passwordValid, tokenValid
}

func (s *UserServiceImpl) Search(userRequest io.ReadCloser) (mod.RequestQuery, error) {
	bodyJSON, err := ioutil.ReadAll(userRequest)
	if err != nil {
		return mod.RequestQuery{}, errors.New("failed to read request")
	}

	var searchResp mod.RequestUser
	var requestQuery mod.RequestQuery

	err = json.Unmarshal(bodyJSON, &searchResp)
	if err != nil {
		return mod.RequestQuery{}, errors.New("failed to deserialize JSON")
	}

	resp, err := s.HandleWorker(searchResp)
	if err != nil {
		return mod.RequestQuery{}, err
	}
	requestQuery.Query = fmt.Sprintf("Latitude: %s, Longitude: %s", resp.RequestSearch.Lng, resp.RequestSearch.Lat)

	return requestQuery, nil
}

func (s *UserServiceImpl) Address(userRequest io.ReadCloser) (mod.RequestQuery, error) {
	bodyJSON, err := ioutil.ReadAll(userRequest)
	if err != nil {
		return mod.RequestQuery{}, errors.New("failed to read request")
	}

	var searchResp mod.RequestUser
	var requestQuery mod.RequestQuery

	err = json.Unmarshal(bodyJSON, &searchResp)
	if err != nil {
		return mod.RequestQuery{}, errors.New("failed to deserialize JSON")
	}

	resp, err := s.HandleWorker(searchResp)
	if err != nil {
		return mod.RequestQuery{}, err
	}
	requestQuery.Query = fmt.Sprintf("Formatted address: %s", resp.Addres)

	return requestQuery, nil
}

func (s *UserServiceImpl) HandleWorker(query mod.RequestUser) (mod.RequestAddress, error) {
	var requestQuery mod.RequestAddress
	cache, err := s.repo.CacheChecker(query)
	if err != nil {
		log.Println(err)
	}
	if cache {
		requestQuery.Addres = cache.Addres
		requestQuery.RequestSearch.Lat = cache.RequestSearch.Lat
		requestQuery.RequestSearch.Lng = cache.RequestSearch.Lng
		return requestQuery, nil
	}

	geocodeResponse, err := s.Geocode(query.RequestQuery)
	if err != nil {
		return mod.RequestAddress{}, errors.New("error in dadata operation")
	}
	for _, v := range geocodeResponse {
		requestQuery.RequestSearch.Lat = v.GeoLat
		requestQuery.RequestSearch.Lng = v.GeoLon
		requestQuery.Addres = v.Result
	}
	err = s.repo.Select(query)
	if err != nil {
		return mod.RequestAddress{}, errors.New("Select query error")
	}
	return requestQuery, nil
}

func (s *UserServiceImpl) Geocode(query mod.RequestQuery) ([]*model.Address, error) {
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

func TokenGenerate(email, password string) (string, error) {
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"Username": email, "Password": password})
	if err != nil {
		return "regData", errors.New("token generation error")
	}
	return tokenString, nil
}

func RefreshToken(email, password string) string {
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"Username": email, "Password": password})
	if err != nil {
		log.Println(err)
	}
	err = repository.RefreshToken(email, password, tokenString)
	if err != nil {
		log.Println(err)
	}
	return tokenString
}
