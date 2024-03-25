package userController

import (
	"context"
	"github.com/gorilla/mux"
	responder "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Responder"
	"github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Service/userService"
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
	vars := mux.Vars(r)
	userID := vars["id"]
	UserInfo, err := c.Service.GetUser(ctx, userID)
	if err != nil {
		responder.HandleError(w, err)
		return
	}
	responder.SendJSONResponse(w, UserInfo)
}

func (c *UserControllerImpl) DelUser(w http.ResponseWriter, r *http.Request) {
	Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	ctx := context.WithValue(r.Context(), "jwt_token", Usertoken)
	vars := mux.Vars(r)
	userID := vars["id"]
	err := c.Service.DelUser(ctx, userID)
	if err != nil {
		responder.HandleError(w, err)
		return
	}
	responder.SendJSONResponse(w, "Succsec")

}
