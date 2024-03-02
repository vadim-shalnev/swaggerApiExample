package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/ekomobile/dadata/v2"
	"github.com/ekomobile/dadata/v2/api/model"
	"github.com/ekomobile/dadata/v2/client"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github/vadim-shalnev/docs"
	"io/ioutil"
	"net/http"
	"strings"
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
type TokenString struct {
	T string `json:"token"`
}
type NewUser struct {
	Username string `json:"user_name"`
	Password string `json:"password"`
}

var AuthUser map[string]NewUser
var UserToken map[string]TokenString
var tokenAuth *jwtauth.JWTAuth

// @title Todo geocode API
// @version 1.0
// @description API Server for search GEOinfo

// @host localhost:8080
// @BasePath /api/address

func main() {
	r := chi.NewRouter()
	r.Get("/api/login", Login)
	r.Post("/api/register", Register)
	r.Route("/api/address", func(r chi.Router) {
		r.Use(AuthMiddleware)
		r.Post("/search", HandleSearch)
		r.Post("/geocode", HandleGeo)
	})
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // Укажите путь к файлу swagger.json
	))
	http.ListenAndServe(":8080", r)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	var regData NewUser
	err = json.Unmarshal(bodyJSON, &regData)
	if err != nil {
		fmt.Println(err)
	}
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"Username": regData.Username, "Password": regData.Password})
	if err != nil {
		fmt.Println(err)
	}
	var tokenStr TokenString
	tokenStr.T = tokenString
	AuthUser = make(map[string]NewUser)
	AuthUser[regData.Username] = regData
	UserToken = make(map[string]TokenString)
	UserToken[regData.Username] = tokenStr

	tokenJSON, err := json.Marshal(tokenStr)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(tokenJSON)
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
	fmt.Println(tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверка подписи методом HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("secret"), nil
	})

	if err != nil {
		fmt.Println("Ошибка разбора токена:", err)
		return
	}
	CheckForUsername := ""
	CheckForPassword := ""
	// Проверка, успешно ли разобран токен
	if token.Valid {
		// Получение клеймов из токена
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Println("Ошибка получения клеймов из токена")
			return
		}

		// Вывод клеймов
		for key, value := range claims {

			if key == "Username" {
				username, _ := value.(string)
				CheckForUsername = username
			}
			if key == "Password" {
				password, _ := value.(string)
				CheckForPassword = password
			}
		}

	} else {
		fmt.Println("Некорректный токен")
	}
	user, _ := AuthUser[CheckForUsername]
	if user.Username != CheckForUsername || user.Password != CheckForPassword {
		w.Write([]byte("Неправильный пароль или имя пользователя"))
		return
	}
	userToken, _ := UserToken[CheckForUsername]
	if userToken.T != tokenString {
		w.Write([]byte("ц-ц-ц, неправильный токен"))
		return
	}
	w.Write([]byte("Вы успешно авторизованы"))
}
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		fmt.Println(tokenString)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Проверка подписи методом HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte("secret"), nil
		})

		if err != nil {
			fmt.Println("Ошибка разбора токена:", err)
			return
		}
		CheckForUsername := ""
		CheckForPassword := ""
		// Проверка, успешно ли разобран токен
		if token.Valid {
			// Получение клеймов из токена
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				fmt.Println("Ошибка получения клеймов из токена")
				return
			}

			// Вывод клеймов
			for key, value := range claims {

				if key == "Username" {
					username, _ := value.(string)
					CheckForUsername = username
				}
				if key == "Password" {
					password, _ := value.(string)
					CheckForPassword = password
				}
			}

		} else {
			fmt.Println("Некорректный токен")
		}
		user, _ := AuthUser[CheckForUsername]
		if user.Username != CheckForUsername || user.Password != CheckForPassword {
			w.Write([]byte("Неправильный пароль или имя пользователя"))
			return
		}
		userToken, _ := UserToken[CheckForUsername]
		if userToken.T != tokenString {
			w.Write([]byte("ц-ц-ц, неправильный токен"))
			return
		}
		next.ServeHTTP(w, r)
	})
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
