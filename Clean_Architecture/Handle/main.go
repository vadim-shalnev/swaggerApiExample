package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type TokenString struct {
	T string `json:"token"`
}
type NewUser struct {
	Username string `json:"user_name"`
	Password string `json:"password"`
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/api/login", Login)
	r.Post("/api/register", Register)
	r.Route("/api/address", func(r chi.Router) {
		r.Use(AuthMiddleware)
		r.Post("/search", Handle)
		r.Post("/geocode", Handle)
	})
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // Укажите путь к файлу swagger.json
	))
	http.ListenAndServe(":8080", r)
}

func Login(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	valid := PrivaseCheker(tokenString)
	w.Write([]byte(valid))

}

func PrivaseCheker(Usertoken string) string {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost:8070/api/login", nil)
	if err != nil {
		log.Fatal("Ошибка при логине", err)
	}

	req.Header.Set("Authorization", "Bearer "+Usertoken)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Ошибка логина к сервису", err)
	}
	defer resp.Body.Close()

	bodyJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Ошибка чтения ответа сервиса ", err)
	}
	return string(bodyJSON)

}

func Register(w http.ResponseWriter, r *http.Request) {
	bodyJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Не удалось прочитать запрос", http.StatusBadRequest)
	}
	var regData NewUser
	err = json.Unmarshal(bodyJSON, &regData)
	if err != nil {
		http.Error(w, "Не удалось дессериализировать JSON", http.StatusBadRequest)
	}

	tokenString := TokenReqGenerate(bodyJSON)

	var tokenStr TokenString
	tokenStr.T = tokenString

	tokenJSON, err := json.Marshal(tokenStr)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(tokenJSON)
}
func TokenReqGenerate(User []byte) string {
	req, err := http.NewRequest("POST", "http://localhost:8070/api/register", bytes.NewReader(User))
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Неверный ответ от сервиса регистрации", err)
	}

	var tokenstr TokenString

	err = json.Unmarshal(bodyJSON, &tokenstr)
	if err != nil {
		log.Fatal("Анмарш токена сервиса реги", err)
	}

	tokenToUser := tokenstr.T
	return tokenToUser

}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		fmt.Println("tests")
		client := &http.Client{}

		req, err := http.NewRequest("GET", "http://localhost:8070/api/login", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		req.Header.Set("Authorization", "Bearer "+Usertoken)

		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()

		bodyJ, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(string(bodyJ))
		if resp.StatusCode != http.StatusOK {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Handle(w http.ResponseWriter, r *http.Request) {
	bodyJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Ошибка чтения запроса пользователя", err)
	}
	Usertoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	client := &http.Client{}
	url := "http://localhost:8090"
	url += r.URL.Path
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		log.Fatal("Ошибка в ответе сервиса поиска", err)
	}
	req.Header.Set("Authorization", "Bearer "+Usertoken)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyJSON, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(bodyJSON)

}
