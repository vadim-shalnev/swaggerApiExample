package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type NewUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
type NewUserResponse struct {
	Email string      `json:"email"`
	Role  string      `json:"role"`
	Token TokenString `json:"token"`
}
type TokenString struct {
	Token string `json:"authController"`
}
type RequestQuery struct {
	Query string `json:"query"`
}

func main() {
	// регистрируемся в сервисе и проверяем 'endpoints User
	client := &http.Client{}
	newUser := NewUserRequest{Email: "user204", Password: "124", Role: "User"}
	userJSON, err := json.Marshal(newUser)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := client.Post("http://localhost:8080/api/register", "application/json", bytes.NewReader(userJSON))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	bodyJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var reginfo NewUserResponse
	err = json.Unmarshal(bodyJSON, &reginfo)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("reginfo", reginfo)

	req, err := http.NewRequest("POST", "http://localhost:8080/api/login", bytes.NewReader(userJSON))
	if err != nil {
		fmt.Println(err)
		return
	}

	token := reginfo.Token.Token
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("login body:", string(bodyByte))

	var loginInfo NewUserResponse
	err = json.Unmarshal(bodyByte, &loginInfo)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("logininfo", loginInfo)

	req, err = http.NewRequest("GET", "http://localhost:8080/api/user/get/13", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	bodyByte, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("get User body:", string(bodyByte))
	err = json.Unmarshal(bodyByte, &loginInfo)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("getuserinfo", loginInfo)
	// Выполняем запросы к геосервису
	var query RequestQuery
	query.Query = "Москва лениня 13"
	queryJSON, err := json.Marshal(query)
	if err != nil {
		fmt.Println(err)
		return
	}
	// создаем тймаут для запроса к геосервису
	time.Sleep(time.Second * 2)
	req, err = http.NewRequest("POST", "http://localhost:8080/api/address/search", bytes.NewReader(queryJSON))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	bodyByte, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("search body:", string(bodyByte))
	var dataSearch RequestQuery
	err = json.Unmarshal(bodyByte, &dataSearch)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("searchinfo", dataSearch)
	req, err = http.NewRequest("POST", "http://localhost:8080/api/address/geocode", bytes.NewReader(queryJSON))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	bodyByte, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("geocode body:", string(bodyByte))
	var geocodeaddres RequestQuery
	err = json.Unmarshal(bodyByte, &geocodeaddres)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("addresinfo", geocodeaddres)

	// Выполняем запрос на удаление к бд
	req, err = http.NewRequest("DELETE", "http://localhost:8080/api/user/del/15", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	bodyByte, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("del is:", string(bodyByte))

}
