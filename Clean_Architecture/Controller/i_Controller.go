package Controller

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type Controller struct {
	Service UserService
}

type UserService interface {
	Register(body io.ReadCloser) (NewUserResponse, error)
	Login(token string, body io.ReadCloser) (NewUserResponse, error)
	Search(body io.ReadCloser) (interface{}, error)
	Address(body io.ReadCloser) (interface{}, error)
	UserInfoChecker(email, password, token string) (bool, NewUserResponse, bool)
}

type TokenString struct {
	Token string `json:"auth"`
}
type NewUserResponse struct {
	Email       string      `json:"email"`
	Password    string      `json:"password"`
	TokenString TokenString `json:"token"`
}

func NewController(service UserService) *Controller {
	return &Controller{Service: service}
}

func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	UserInfo, err := c.Service.Register(r.Body)
	if err != nil {
		handleError(w, err)
		return
	}
	sendJSONResponse(w, UserInfo)
}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	UserInfo, err := c.Service.Login(Usertoken, r.Body)
	if err != nil {
		handleError(w, err)
		return
	}
	sendJSONResponse(w, UserInfo)
}

func (c *Controller) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		_, _, token := c.Service.UserInfoChecker("", "", Usertoken)
		if !token {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (c *Controller) HandleSearch(w http.ResponseWriter, r *http.Request) {
	respSearch, err := c.Service.Search(r.Body)
	if err != nil {
		handleError(w, err)
		return
	}
	sendJSONResponse(w, respSearch)
}

func (c *Controller) HandleGeo(w http.ResponseWriter, r *http.Request) {
	respGeo, err := c.Service.Address(r.Body)
	if err != nil {
		handleError(w, err)
		return
	}
	sendJSONResponse(w, respGeo)
}

func handleError(w http.ResponseWriter, err error) {
	var status int
	switch err.Error() {
	case "не удалось прочитать запрос", "не удалось дессериализировать JSON":
		status = http.StatusBadRequest
	case "не верный логин", "не верный пароль", "вы успешно вышли из сервиса":
		status = http.StatusUnauthorized
	case "ошибка в работе dadata", "ошибка запроса Select":
		status = http.StatusInternalServerError
	default:
		status = http.StatusInternalServerError
	}
	http.Error(w, err.Error(), status)
}

func sendJSONResponse(w http.ResponseWriter, data interface{}) {
	respJSON, err := json.Marshal(data)
	if err != nil {
		handleError(w, err)
		return
	}
	w.Write(respJSON)
}
