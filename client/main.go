package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RequestAddressSearch struct {
	Query string `json:"query"`
}

func main() {
	var SearchResp RequestAddressSearch
	//var addr Address
	addr := "комсомольский проспект 18"
	SearchResp.Query = addr
	queryJSON, err := json.Marshal(SearchResp)
	if err != nil {
		fmt.Println(err)
	}
	respSwag, err := http.Get("http://localhost:8080/swagger/index.html")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(respSwag.Body)
	resp, err := http.Post("http://localhost:8080/api/address/geocode", "application/json", bytes.NewReader(queryJSON))
	if err != nil {
		fmt.Println(err)
	}
	bodyJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var datainfo RequestAddressSearch

	err = json.Unmarshal(bodyJSON, &datainfo)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(datainfo)

}
