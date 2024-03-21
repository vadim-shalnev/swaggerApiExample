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
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func Register(User io.ReadCloser) (s.NewUserResponse, error) {

	var regData s.NewUserResponse

	bodyJSON, err := ioutil.ReadAll(User)
	if err != nil {
		return s.NewUserResponse{}, errors.New("не удалось прочитать запрос")
	}
	err = json.Unmarshal(bodyJSON, &regData)
	if err != nil {
		return s.NewUserResponse{}, errors.New("не удалось дессериализировать JSON")
	}
	tokenAuth, err := TokenGenerate(regData.Email, regData.Password)
	if err != nil {
		return s.NewUserResponse{}, err
	}
	// Устанавливаем токен и добавляем пользователя в БД
	regData.TokenString.Token = tokenAuth

	err = Repository.Create(regData)
	if err != nil {
		return regData, err
	}

	return regData, nil

}
func Login(Usertoken string, User io.ReadCloser) (s.NewUserResponse, error) {

	var regData s.NewUserResponse

	bodyJSON, err := ioutil.ReadAll(User)
	if err != nil {
		return s.NewUserResponse{}, errors.New("не удалось прочитать запрос")
	}
	err = json.Unmarshal(bodyJSON, &regData)
	if err != nil {
		return s.NewUserResponse{}, errors.New("не удалось дессериализировать JSON")
	}
	// Проверяем логин, пароль и токен
	Email, Password, Token := UserInfo_Checker(regData.Email, regData.Password, Usertoken)
	if !Email {
		return s.NewUserResponse{}, errors.New("не верный логин")
	}
	if !Password {
		return s.NewUserResponse{}, errors.New("не верный пароль")
	}
	if !Token {
		freshToken := RefreshToken(regData.Email, regData.Password)
		regData.TokenString.Token = freshToken
		return s.NewUserResponse{}, errors.New("вы успешно вышли из сервиса")
	}

	return regData, nil
}

func TokenGenerate(email, password string) (string, error) {
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"Username": email, "Password": password})
	if err != nil {
		return "regData", errors.New("ошибка генерации токена")
	}
	return tokenString, nil
}
func RefreshToken(email, pawwsord string) string {
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"Username": email, "Password": pawwsord})
	if err != nil {
		log.Println(err)
	}
	err := Repository.RefreshToken(email, pawwsord, tokenString)
	if err != nil {
		log.Println(err)
	}
	return tokenString
}

func UserInfo_Checker(email, password, token string) (bool, bool, bool) {
	var Email, Password, Token bool
	Email = Repository.CheckEmail(email)
	Password = Repository.CheckPassword(password)
	Token = Repository.CheckToken(token)
	return Email, Password, Token
}

func HandleSearch(UserRequest io.ReadCloser) (s.requestQuery, error) {

	bodyJSON, err := ioutil.ReadAll(UserRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var SearchResp s.RequestUser
	var requestSearch s.RequestAddress
	var requestQuery s.RequestQuery

	err = json.Unmarshal(bodyJSON, &SearchResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	resp, err := HandleWorker(SearchResp.RequestQuery.Query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	requestQuery.Query = resp.RequestQuery

	requestJSON, err := json.Marshal(requestQuery)
	if err != nil {
		log.Println(err)
	}
	return requestJSON, nil
}

func HandleAddress(UserRequest io.ReadCloser) (s.requestQuery, error) {

	bodyJSON, err := ioutil.ReadAll(UserRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var SearchResp s.RequestUser

	err = json.Unmarshal(bodyJSON, &SearchResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var requestAddress s.RequestQuery

	geocodeResponse, err := Geocode(SearchResp)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range geocodeResponse {
		requestAddress.Query = v.Result
	}

	requestJSON, err := json.Marshal(requestAddress)
	if err != nil {
		log.Println(err)
	}
	return requestJSON, nil
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
func HandleWorker() {
	cache, err := Repository.CacheChecker(SearchResp.Email, SearchResp.RequestQuery.Query)
	if err != nil {
		log.Println(err)
	}
	if cache {
		requestQuery.Query = cache
		return requestQuery, nil
	}

	geocodeResponse, err := Geocode(SearchResp)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range geocodeResponse {
		requestSearch.Lat = v.GeoLat
		requestSearch.Lng = v.GeoLon
		requestQuery.Query += fmt.Sprintf("Широта: %s Долгота %s", v.GeoLat, v.GeoLon)
	}
}
func CacheChecker(Email, Query string) (string, error) {
	cache, err := Repository.CacheChecker(SearchResp.Email, SearchResp.RequestQuery.Query)
	if err != nil {
		log.Println(err)
	}
	if !cache {
		err = Repository.Select(SearchResp.Email, SearchResp.RequestQuery.Query)
		if err != nil {
			log.Println(err)
		}
	} else {
		var request s.RequestQuery
		request.Query = cache
		return request, nil
	}
}
