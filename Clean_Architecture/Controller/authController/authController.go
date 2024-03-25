package authController

import (
	"context"
	"encoding/json"
	mod "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Models"
	responder "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Responder"
	"github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Service/authService"
	"io/ioutil"
	"net/http"
	"strings"
)

type AuthControllerImpl struct {
	Auth authService.AuthService
}

type AuthController interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

func NewAuthController(auth authService.AuthService) *AuthControllerImpl {
	return &AuthControllerImpl{Auth: auth}
}

func (c *AuthControllerImpl) Register(w http.ResponseWriter, r *http.Request) {
	var regData mod.NewUserRequest
	bodyJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responder.HandleError(w, err)
		return
	}
	err = json.Unmarshal(bodyJSON, &regData)
	if err != nil {
		responder.HandleError(w, err)
		return
	}

	UserInfo, err := c.Auth.Register(r.Context(), regData)
	if err != nil {
		responder.HandleError(w, err)
		return
	}
	responder.SendJSONResponse(w, UserInfo)
}

func (c *AuthControllerImpl) Login(w http.ResponseWriter, r *http.Request) {
	var loginData mod.NewUserRequest
	bodyJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responder.HandleError(w, err)
		return
	}
	err = json.Unmarshal(bodyJSON, &loginData)
	if err != nil {
		responder.HandleError(w, err)
		return
	}

	Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	ctx := context.WithValue(r.Context(), "jwt_token", Usertoken)
	UserInfo, err := c.Auth.Login(ctx, loginData)
	if err != nil {
		responder.HandleError(w, err)
		return
	}
	responder.SendJSONResponse(w, UserInfo)
}

func (c *AuthControllerImpl) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		_, _, token := c.Auth.UserInfoChecker(r.Context(), "", "", Usertoken)
		if !token {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
