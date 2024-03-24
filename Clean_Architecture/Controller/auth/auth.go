package auth

import (
	"context"
	"encoding/json"
	mod "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Models"
	"io/ioutil"
	"net/http"
	"strings"
)

type AuthControllerImpl struct {
	Auth      service.auth
	Responder responder.Responder
}

type AuthController interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

func NewAuthController(auth service.auth) Auth {
	return &AuthImpl{auth: auth}
}

func (c *AuthControllerImpl) Register(w http.ResponseWriter, r *http.Request) {
	var regData mod.NewUserRequest
	bodyJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(w, err)
		return
	}
	err = json.Unmarshal(bodyJSON, &regData)
	if err != nil {
		handleError(w, err)
		return
	}

	UserInfo, err := c.Service.Register(r.Context(), regData)
	if err != nil {
		handleError(w, err)
		return
	}
	sendJSONResponse(w, UserInfo)
}

func (c *AuthControllerImpl) Login(w http.ResponseWriter, r *http.Request) {
	var loginData mod.NewUserRequest
	bodyJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(w, err)
		return
	}
	err = json.Unmarshal(bodyJSON, &loginData)
	if err != nil {
		handleError(w, err)
		return
	}

	Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	ctx := context.WithValue(r.Context(), "jwt_token", Usertoken)
	UserInfo, err := c.Service.Login(ctx, loginData)
	if err != nil {
		handleError(w, err)
		return
	}
	sendJSONResponse(w, UserInfo)
}
