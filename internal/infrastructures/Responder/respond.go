package Responder

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type Respond struct {
	log *zap.Logger
}

type Responder interface {
	SendJSONResponse(w http.ResponseWriter, data interface{})
	HandleError(w http.ResponseWriter, err error)
}

func NewResponder(log *zap.Logger) *Respond {
	return &Respond{log: log}
}

func (r *Respond) SendJSONResponse(w http.ResponseWriter, data interface{}) {
	respJSON, err := json.Marshal(data)
	if err != nil {
		r.HandleError(w, err)
		return
	}
	w.Write(respJSON)
}

func (r *Respond) HandleError(w http.ResponseWriter, err error) {
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
