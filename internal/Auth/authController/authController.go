package authController

import (
	"encoding/json"
	mod "github.com/vadim-shalnev/swaggerApiExample/Models"
	"github.com/vadim-shalnev/swaggerApiExample/internal/Auth/authService"
	"github.com/vadim-shalnev/swaggerApiExample/internal/infrastructures/Responder"
	"net/http"
)

type Authcontroller struct {
	Auth      authService.AuthService
	responder Responder.Responder
}

type AuthController interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

func NewAuthController(auth authService.AuthService, responder Responder.Responder) *Authcontroller {
	return &Authcontroller{Auth: auth, responder: responder}
}

// Register @Summary Регистрация нового пользователя
// @Description Регистрация нового пользователя с указанным email и паролем
// @Tags reg
// @Accept json
// @Produce json
// @Param   User body mod.NewUserRequest true "Данные пользователя"
// @Success 200 {object} string "Успешная регистрация"
// @Router /api/register [post]
func (c *Authcontroller) Register(w http.ResponseWriter, r *http.Request) {
	var regData mod.NewUserRequest
	err := json.NewDecoder(r.Body).Decode(&regData)
	if err != nil {
		c.responder.HandleError(w, err)
		return
	}
	token, err := c.Auth.Register(regData)
	if err != nil {
		c.responder.HandleError(w, err)
		return
	}
	c.responder.SendJSONResponse(w, token)
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
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		c.responder.HandleError(w, err)
		return
	}
	UserInfo, err := c.Auth.Login(loginData)
	if err != nil {
		c.responder.HandleError(w, err)
		return
	}
	c.responder.SendJSONResponse(w, UserInfo)
}

func (c *Authcontroller) Logout(w http.ResponseWriter, r *http.Request) {
	niltoken, err := c.Auth.Logout()
	if err != nil {
		c.responder.HandleError(w, err)
		return
	}
	c.responder.SendJSONResponse(w, niltoken)
}
