package user

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type UserControllerImpl struct {
	Service   service.user
	Responder responder.Responder
}
type UserController interface {
	GetUser(w http.ResponseWriter, r *http.Request)
}

func NewUserController(service service.user) UserController {
	return &UserControllerImpl{Service: service}
}

func (c *UserControllerImpl) GetUser(w http.ResponseWriter, r *http.Request) {
	Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	ctx := context.WithValue(r.Context(), "jwt_token", Usertoken)
	vars := mux.Vars(r)
	userID := vars["id"]
	UserInfo, err := c.Service.GetUser(ctx, userID)
	if err != nil {
		handleError(w, err)
		return
	}
	sendJSONResponse(w, UserInfo)
}
