package userController

import (
	"context"
	"github.com/go-chi/chi/v5"
	responder "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/internal/Responder"
	"github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/internal/Service/userService"
	"log"
	"net/http"
	"strings"
)

type UserControllerImpl struct {
	Service userService.UserService
}
type UserController interface {
	GetUser(w http.ResponseWriter, r *http.Request)
	DelUser(w http.ResponseWriter, r *http.Request)
}

func NewUserController(service userService.UserService) *UserControllerImpl {
	return &UserControllerImpl{Service: service}
}

func (c *UserControllerImpl) GetUser(w http.ResponseWriter, r *http.Request) {
	Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	ctx := context.WithValue(r.Context(), "jwt_token", Usertoken)
	userID := chi.URLParam(r, "id")
	log.Println("get id is", userID)
	UserInfo, err := c.Service.GetUser(ctx, userID)
	if err != nil {
		responder.HandleError(w, err)
		return
	}
	log.Println(UserInfo)
	responder.SendJSONResponse(w, UserInfo)
}

func (c *UserControllerImpl) DelUser(w http.ResponseWriter, r *http.Request) {
	Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	ctx := context.WithValue(r.Context(), "jwt_token", Usertoken)
	userID := chi.URLParam(r, "id")
	log.Println("del id is", userID)
	err := c.Service.DelUser(ctx, userID)
	if err != nil {
		responder.HandleError(w, err)
		return
	}
	responder.SendJSONResponse(w, "Succsec")
}
