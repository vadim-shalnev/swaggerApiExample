package Controller

/*
func Register(w http.ResponseWriter, r *http.Request) {

	UserInfo, err := Service.Register(r.Body)
	if err != nil {
		if err.Error() == "не удалось прочитать запрос" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		if err.Error() == "не удалось дессериализировать JSON" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		if err.Error() == "ошибка генерации токена" {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if err.Error() == "не удалось добавить нового пользователя в БД" {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	tokenJSON, err := json.Marshal(UserInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(tokenJSON)
}

func Login(w http.ResponseWriter, r *http.Request) {
	Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	UserInfo, err := Service.Login(Usertoken, r.Body)
	if err != nil {
		if err.Error() == "не удалось прочитать запрос" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		if err.Error() == "не удалось дессериализировать JSON" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		if err.Error() == "не верный логин" {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		if err.Error() == "не верный пароль" {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		if err.Error() == "вы успешно вышли из сервиса" {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
	}

	tokenJSON, err := json.Marshal(UserInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(tokenJSON)
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		_, _, token := Service.UserInfo_Checker("", "", Usertoken)
		if !token {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func HandleSearch(w http.ResponseWriter, r *http.Request) {
	respSearch, err := Service.Search(r.Body)

	respSearchJSON, err := json.Marshal(respSearch)
	if err != nil {
		if err.Error() == "не удалось прочитать запрос" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		if err.Error() == "не удалось дессериализировать JSON" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		if err.Error() == "ошибка в работе dadata" {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if err.Error() == "ошибка запроса Select" {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	w.Write(respSearchJSON)
}
func HandleGeo(w http.ResponseWriter, r *http.Request) {
	respGeo, err := Service.Address(r.Body)

	respGeoJSON, err := json.Marshal(respGeo)
	if err != nil {
		if err.Error() == "не удалось прочитать запрос" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		if err.Error() == "не удалось дессериализировать JSON" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		if err.Error() == "ошибка в работе dadata" {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if err.Error() == "ошибка запроса Select" {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	w.Write(respGeoJSON)
}

*/
