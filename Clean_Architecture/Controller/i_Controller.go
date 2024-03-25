package Controller

/*
type Controller struct {
	Service service.UserService
	Auth    Auth
}
type Auth interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	AuthMiddleware(next http.Handler) http.Handler
	HandleSearch(w http.ResponseWriter, r *http.Request)
	HandleGeo(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
}

func NewController(service service.UserService) *Controller {
	return &Controller{Service: service}
}

func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	UserInfo, err := c.Service.Register(r.Context(), r.Body)
	if err != nil {
		handleError(w, err)
		return
	}
	sendJSONResponse(w, UserInfo)
}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	ctx := context.WithValue(r.Context(), "jwt_token", Usertoken)
	UserInfo, err := c.Service.Login(ctx, r.Body)
	if err != nil {
		handleError(w, err)
		return
	}
	sendJSONResponse(w, UserInfo)
}
func (c *Controller) GetUser(w http.ResponseWriter, r *http.Request) {
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

func (c *Controller) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		_, _, token := c.Service.UserInfoChecker(r.Context(), "", "", Usertoken)
		if !token {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (c *Controller) HandleSearch(w http.ResponseWriter, r *http.Request) {
	Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	ctx := context.WithValue(r.Context(), "jwt_token", Usertoken)
	respSearch, err := c.Service.Search(ctx, r.Body)
	if err != nil {
		handleError(w, err)
		return
	}
	sendJSONResponse(w, respSearch)
}

func (c *Controller) HandleGeo(w http.ResponseWriter, r *http.Request) {
	Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	ctx := context.WithValue(r.Context(), "jwt_token", Usertoken)
	respGeo, err := c.Service.Address(ctx, r.Body)
	if err != nil {
		handleError(w, err)
		return
	}
	sendJSONResponse(w, respGeo)
}

func handleError(w http.ResponseWriter, err error) {
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

func sendJSONResponse(w http.ResponseWriter, data interface{}) {
	respJSON, err := json.Marshal(data)
	if err != nil {
		handleError(w, err)
		return
	}
	w.Write(respJSON)
}

*/