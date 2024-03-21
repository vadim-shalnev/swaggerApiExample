package Controller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	Service "logic"
	"net/http"
	"strings"
)

// @title Todo geocode API
// @version 1.0
// @description API Server for search GEOinfo

// @host localhost:8080
// @BasePath /api/address

// Login @Login
// @Summary User login
// @Description User login with JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT token"
// @Success 200 {string} string "User successfully logged in"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Router /api/login [get]

// Register @Register
// @Summary Register
// @Tags Reg in service
// @Description Register a new user
// @Accept json
// @Produce json
// @Param input body NewUser true "User object for registration"
// @Success 200 {integer} integer 1
// @Failure 404 {error} http.Error
// @Failure 500 {error} http.Error
// @Router /api/register [post]

const ApiKey string = "22d3fa86b8743e497b32195cbc690abc06b42436"
const SecretKey string = "adf07bdd63b240ae60087efd2e72269b9c65cc91"

func Register(w http.ResponseWriter, r *http.Request) {

	UserInfo := Service.Register(r.Body)

	tokenJSON, err := json.Marshal(UserInfo)
	if err != nil {
		if err.Error() == "не удалось прочитать запрос" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		if err.Error() == "не удалось дессериализировать JSON" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		if err.Error() == "ошибка генерации токена" {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	w.Write(tokenJSON)
}

func Login(w http.ResponseWriter, r *http.Request) {
	Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	UserInfo := Service.Login(Usertoken, r.Body)

	tokenJSON, err := json.Marshal(UserInfo)
	if err != nil {
		if err.Error() == "не удалось прочитать запрос" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		if err.Error() == "не удалось дессериализировать JSON" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		if err.Error() == "не верный логин" {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		if err.Error() == "не верный пароль" {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		if err.Error() == "вы успешно вышли из сервиса" {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
	}
	w.Write(tokenJSON)
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		_, _, token := Service.UserInfo_Checker("", "", Usertoken)
		if !token {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Handle @Controller
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
// @Router /search [post]

func HandleSearch(w http.ResponseWriter, r *http.Request) {
	respSearch := Service.Search(r.Body)

	respSearchJSON, err := json.Marshal(respSearch)
	if err != nil {
		///
	}
	w.Write(respSearchJSON)
}

func Handle(w http.ResponseWriter, r *http.Request) {
	bodyJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Ошибка чтения запроса пользователя", err)
	}
	Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	client := &http.Client{}
	url := "http://Service:8090"
	url += r.URL.Path
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		log.Fatal("Ошибка в ответе сервиса поиска", err)
	}
	req.Header.Set("Authorization", "Bearer "+Usertoken)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyJSON, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(bodyJSON)

}
