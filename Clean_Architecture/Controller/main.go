package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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
type TokenString struct {
	T string `json:"token"`
}
type NewUser struct {
	Username string `json:"user_name"`
	Password string `json:"password"`
}

// @title Todo geocode API
// @version 1.0
// @description API Server for search GEOinfo

// @host localhost:8080
// @BasePath /api/address

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/api/login", Login)
	r.Post("/api/register", Register)
	r.Route("/api/address", func(r chi.Router) {
		r.Use(AuthMiddleware)
		r.Post("/search", Handle)
		r.Post("/geocode", Handle)
	})
	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
	http.ListenAndServe(":8080", r)
}

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
func Login(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	valid := PrivaseCheker(tokenString)
	w.Write([]byte(valid))

}

func PrivaseCheker(Usertoken string) string {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://Repository:8070/api/login", nil)
	if err != nil {
		log.Fatal("Ошибка при логине", err)
	}

	req.Header.Set("Authorization", "Bearer "+Usertoken)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Ошибка логина к сервису", err)
	}
	defer resp.Body.Close()

	bodyJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Ошибка чтения ответа сервиса ", err)
	}
	return string(bodyJSON)

}

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
func Register(w http.ResponseWriter, r *http.Request) {
	bodyJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Не удалось прочитать запрос", http.StatusBadRequest)
	}
	var regData NewUser
	err = json.Unmarshal(bodyJSON, &regData)
	if err != nil {
		http.Error(w, "Не удалось дессериализировать JSON", http.StatusBadRequest)
	}

	tokenString := TokenReqGenerate(bodyJSON)

	var tokenStr TokenString
	tokenStr.T = tokenString

	tokenJSON, err := json.Marshal(tokenStr)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(tokenJSON)
}
func TokenReqGenerate(User []byte) string {
	req, err := http.NewRequest("POST", "http://Repository:8070/api/register", bytes.NewReader(User))
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Неверный ответ от сервиса регистрации", err)
	}

	var tokenstr TokenString

	err = json.Unmarshal(bodyJSON, &tokenstr)
	if err != nil {
		log.Fatal("Анмарш токена сервиса реги", err)
	}

	tokenToUser := tokenstr.T
	return tokenToUser

}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		fmt.Println("tests")
		client := &http.Client{}

		req, err := http.NewRequest("GET", "http://Repository:8070/api/login", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		req.Header.Set("Authorization", "Bearer "+Usertoken)

		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()

		bodyJ, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(string(bodyJ))
		if resp.StatusCode != http.StatusOK {
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
