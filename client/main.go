package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	// регистрируемся в сервисе и проверяем 'endpoints user
	client := &http.Client{}
	newUser := NewUserRequest{Email: "user1234", Password: "123", Role: "user"}
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
	var datainfo NewUserResponse
	err = json.Unmarshal(bodyJSON, &datainfo)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(datainfo)

	req, err := http.NewRequest("POST", "http://localhost:8080/api/login", bytes.NewReader(userJSON))
	if err != nil {
		fmt.Println(err)
		return
	}

	token := datainfo.Token.Token
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var loginInfo NewUserResponse
	err = json.Unmarshal(bodyByte, &loginInfo)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bodyByte)

	req, err = http.NewRequest("GET", "http://localhost:8080/user/get/{1}", nil)
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
	err = json.Unmarshal(bodyByte, &loginInfo)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bodyByte)
	// Выполняем запросы к геосервису
	var query RequestQuery
	query.Query = "Москва лениня 13"
	queryJSON, err := json.Marshal(query)
	if err != nil {
		fmt.Println(err)
		return
	}
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
	var dataSearch RequestQuery
	err = json.Unmarshal(bodyByte, &dataSearch)
	if err != nil {
		fmt.Println(err)
		return
	}
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

	err = json.Unmarshal(bodyByte, &dataSearch)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bodyByte)

	// Выполняем запрос на удаление к бд
	req, err = http.NewRequest("DELETE", "http://localhost:8080/user/del/{1}", nil)
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
	fmt.Println(resp.Body)

}
