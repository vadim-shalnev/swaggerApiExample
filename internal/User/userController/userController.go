package userController

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/vadim-shalnev/swaggerApiExample/internal/User/userService"
	responder "github.com/vadim-shalnev/swaggerApiExample/internal/infrastructures/Responder"
	"log"
	"net/http"
	"strings"
)

type Usercontroller struct {
	Service userService.UserService
}
type UserController interface {
	GetUser(w http.ResponseWriter, r *http.Request)
	DelUser(w http.ResponseWriter, r *http.Request)
	ListUsers(w http.ResponseWriter, r *http.Request)
}

func NewUserController(service userService.UserService) *Usercontroller {
	return &Usercontroller{Service: service}
}

// GetUser @Summary Получить информацию о пользователе
// @Description Получить информацию о пользователе по его ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} Models.NewUserResponse "Пользователь"
// @Router /users/get/{id} [get]
func (c *Usercontroller) GetUser(w http.ResponseWriter, r *http.Request) {
	Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	ctx := context.WithValue(r.Context(), "jwt_token", Usertoken)
	userID := chi.URLParam(r, "id")
	UserInfo, err := c.Service.GetUser(ctx, userID)
	if err != nil {
		responder.HandleError(w, err)
		return
	}
	log.Println(UserInfo)
	responder.SendJSONResponse(w, UserInfo)
}

// DelUser @Summary Удалить пользователя
// @Description Удалить пользователя по его ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {string} string "Succsec"
// @Router /users/get/{id} [delete]
func (c *Usercontroller) DelUser(w http.ResponseWriter, r *http.Request) {
	Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	ctx := context.WithValue(r.Context(), "jwt_token", Usertoken)
	userID := chi.URLParam(r, "id")
	err := c.Service.DelUser(ctx, userID)
	if err != nil {
		responder.HandleError(w, err)
		return
	}
	responder.SendJSONResponse(w, "Succsec")
}

// не отражаю это в свагере, т.к. по идее это админская функция, а не для клиента
func (c *Usercontroller) ListUsers(w http.ResponseWriter, r *http.Request) {
	Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	ctx := context.WithValue(r.Context(), "jwt_token", Usertoken)
	UserInfo, err := c.Service.ListUsers(ctx)
	if err != nil {
		responder.HandleError(w, err)
		return
	}
	responder.SendJSONResponse(w, UserInfo)
}
