package authController

import (
	"context"
	"encoding/json"
	mod "github.com/vadim-shalnev/swaggerApiExample/Models"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Auth/authService"
	responder "github.com/vadim-shalnev/swaggerApiExample/internal/Responder"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Authcontroller struct {
	Auth authService.AuthService
}

type AuthController interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

func NewAuthController(auth authService.AuthService) *Authcontroller {
	return &Authcontroller{Auth: auth}
}

// Register @Summary Регистрация нового пользователя
// @Description Регистрация нового пользователя с указанным email и паролем
// @Tags reg
// @Accept json
// @Produce json
// @Param   User body mod.NewUserRequest true "Данные пользователя"
// @Success 200 {object} mod.NewUserResponse "Успешная регистрация"
// @Router /api/register [post]
func (c *Authcontroller) Register(w http.ResponseWriter, r *http.Request) {
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

// Login @Summary Вход в систему
// @Description Вход в систему с указанным email и паролем
// @Tags reg
// @Accept json
// @Produce json
// @Param   User body mod.NewUserRequest true "Данные пользователя"
// @Success 200 {object} mod.NewUserResponse "Успешный вход в систему"
// @Router /api/login [post]
func (c *Authcontroller) Login(w http.ResponseWriter, r *http.Request) {
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

func (c *Authcontroller) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		_, _, token := c.Auth.VerifyToken(Usertoken)
		log.Println("midleware")
		if !token {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
