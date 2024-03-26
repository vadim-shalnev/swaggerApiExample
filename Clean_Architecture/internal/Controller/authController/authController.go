package authController

import (
	"context"
	"encoding/json"
	mod "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Models"
	responder "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/internal/Responder"
	"github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/internal/Service/authService"
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

// @Summary Регистрация нового пользователя
// @Description Регистрация нового пользователя с указанным email и паролем
// @Tags users
// @Accept json
// @Produce json
// @Param   user body mod.NewUserRequest true "Данные пользователя"
// @Success 200 {object} mod.NewUserResponse "Успешная регистрация"
// @Router /api/register [post]

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

// @Summary Вход в систему
// @Description Вход в систему с указанным email и паролем
// @Tags users
// @Accept json
// @Produce json
// @Param   user body mod.NewUserRequest true "Данные пользователя"
// @Success 200 {object} mod.NewUserResponse "Успешный вход в систему"
// @Router /api/login [post]

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
		_, _, token := c.Auth.VerifyToken(Usertoken)
		if !token {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
