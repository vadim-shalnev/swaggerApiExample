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
	"github.com/golang-jwt/jwt/v5"
	mod "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Models"
	repository "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Repository"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

const (
	ApiKey    = "22d3fa86b8743e497b32195cbc690abc06b42436"
	SecretKey = "adf07bdd63b240ae60087efd2e72269b9c65cc91"
)

type UserServiceImpl struct {
	repo repository.Repository
}

type UserService interface {
	Register(ctx context.Context, userBody io.ReadCloser) (mod.NewUserResponse, error)
	Login(ctx context.Context, userBody io.ReadCloser) (mod.NewUserResponse, error)
	Search(ctx context.Context, userRequest io.ReadCloser) (interface{}, error)
	Address(ctx context.Context, userRequest io.ReadCloser) (interface{}, error)
	UserInfoChecker(ctx context.Context, email, password, token string) (bool, bool, bool)
	HandleWorker(ctx context.Context, query mod.RequestQuery) (mod.RequestAddress, error)
	RefreshToken(ctx context.Context, email, password string) string
	Geocode(query mod.RequestQuery) ([]*model.Address, error)
	GetUser(ctx context.Context, id string) (mod.NewUserResponse, error)
}

func NewUserServiceImpl(repository repository.Repository) *UserServiceImpl {
	return &UserServiceImpl{repo: repository}
}

func (s *UserServiceImpl) Register(ctx context.Context, userBody io.ReadCloser) (mod.NewUserResponse, error) {
	var regData mod.NewUserRequest
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
	var userResponse mod.NewUserResponse
	userResponse.Email = regData.Email
	userResponse.Token.Token = tokenAuth
	userResponse.Role = regData.Role

	// Хэшируем пароль и добавляем его в запрос к БД
	hashPassword, err := s.crypto.HashPassword(regData.Password)
	if err != nil {
		return mod.NewUserResponse{}, errors.New("failed to hash password")
	}
	regData.Password = hashPassword

	err = s.repo.CreateUser(ctx, regData)
	if err != nil {
		return mod.NewUserResponse{}, errors.New("failed to add new user to the database")
	}

	return userResponse, nil
}

func (s *UserServiceImpl) Login(ctx context.Context, userBody io.ReadCloser) (mod.NewUserResponse, error) {
	var regData mod.UserRequest

	bodyJSON, err := ioutil.ReadAll(userBody)
	if err != nil {
		return mod.NewUserResponse{}, errors.New("failed to read request")
	}
	err = json.Unmarshal(bodyJSON, &regData)
	if err != nil {
		return mod.NewUserResponse{}, errors.New("failed to deserialize JSON")
	}
	var userResponse mod.NewUserResponse
	userResponse.Email = regData.Email

	userToken := ctx.Value("jwt_token").(string)
	emailValid, passwordValid, tokenValid := s.UserInfoChecker(ctx, regData.Email, regData.Password, userToken)
	if !emailValid {
		return userResponse, errors.New("invalid email")
	}
	if !passwordValid {
		return userResponse, errors.New("invalid password")
	}
	if !tokenValid {
		freshToken := s.RefreshToken(ctx, regData.Email, regData.Password)
		userResponse.Token.Token = freshToken
		return userResponse, errors.New("you have successfully logged out of the service")
	}

	return userResponse, nil
}
func (s *UserServiceImpl) GetUser(ctx context.Context, id string) (mod.NewUserResponse, error) {
	var userResponse mod.NewUserResponse
	userID, _ := strconv.Atoi(id)
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		return userResponse, err
	}
	userToken := ctx.Value("jwt_token").(string)
	userResponse.Email = user.Email
	userResponse.Role = user.Role
	userResponse.Token.Token = userToken

	return userResponse, nil
}

func (s *UserServiceImpl) UserInfoChecker(ctx context.Context, email, password, token string) (bool, bool, bool) {
	emailValid := s.repo.CheckEmail(ctx, email)
	passwordValid := s.repo.CheckPassword(ctx, password)
	_, tokenValid := VerifyToken(token, "exp")
	return emailValid, passwordValid, tokenValid
}
func VerifyToken(tokenString, searchIntoken string) (string, bool) {
	// Парсим токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Возвращаем секретный ключ для проверки подписи
		return []byte("secret"), nil
	})
	if err != nil {
		return "", false
	}
	if !token.Valid {
		return "", false
	}

	var search string
	if searchIntoken == "exp" {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			exp := int64(claims[searchIntoken].(float64))
			if time.Now().Unix() > exp {
				return "", false
			}
		}
	} else {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			search = claims[searchIntoken].(string)
		}
	}

	return search, true
}

func (s *UserServiceImpl) Search(ctx context.Context, userRequest io.ReadCloser) (interface{}, error) {
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

	resp, err := s.HandleWorker(ctx, searchRequest)
	if err != nil {
		return mod.RequestQuery{}, err
	}
	responseQuery.Query = fmt.Sprintf("Latitude: %s, Longitude: %s", resp.RequestSearch.Lng, resp.RequestSearch.Lat)

	return responseQuery, nil
}

func (s *UserServiceImpl) Address(ctx context.Context, userRequest io.ReadCloser) (interface{}, error) {
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

	resp, err := s.HandleWorker(ctx, searchRequest)
	if err != nil {
		return mod.RequestQuery{}, err
	}
	responseQuery.Query = fmt.Sprintf("Formatted address: %s", resp.Addres)

	return responseQuery, nil
}

func (s *UserServiceImpl) HandleWorker(ctx context.Context, query mod.RequestQuery) (mod.RequestAddress, error) {
	var requestQuery mod.RequestAddress
	ok, cache, err := s.CacheChecker(ctx, query, 5)
	if err != nil {
		log.Println(err)
	}
	if ok {
		requestQuery.Addres = cache.Addres
		requestQuery.RequestSearch.Lat = cache.RequestSearch.Lat
		requestQuery.RequestSearch.Lng = cache.RequestSearch.Lng
		return requestQuery, nil
	}

	geocodeResponse, err := s.Geocode(query)
	if err != nil {
		return mod.RequestAddress{}, errors.New("error in dadata operation")
	}
	for _, v := range geocodeResponse {
		requestQuery.RequestSearch.Lat = v.GeoLat
		requestQuery.RequestSearch.Lng = v.GeoLon
		requestQuery.Addres = v.Result
	}
	err = s.repo.Insert(ctx, query)
	if err != nil {
		return mod.RequestAddress{}, errors.New("Select query error")
	}
	return requestQuery, nil
}

func (s *UserServiceImpl) CacheChecker(ctx context.Context, query mod.RequestQuery, ttl int) (bool, mod.RequestAddress, error) {
	userToken := ctx.Value("jwt_token").(string)
	email, _ := VerifyToken(userToken, "email")
	// получаем id пользователя
	userID, err := s.repo.GetByEmail(ctx, email)
	// идем в репо за последними запросами
	searchHistory, err := s.repo.CacheChecker(ctx, userID, ttl)
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
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"Username": email,
		"Password": password,
		"Exp":      time.Now().Add(time.Second * 60).Unix(),
	})
	if err != nil {
		return "", errors.New("token generation error")
	}
	return tokenString, nil
}

func (s *UserServiceImpl) RefreshToken(ctx context.Context, email, password string) string {
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"Username": email,
		"Password": password,
		"Exp":      time.Now().Add(time.Second * 60).Unix(),
	})
	if err != nil {
		log.Println(err)
	}

	return tokenString
}
