package main

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Route"
)

const ApiKey string = "22d3fa86b8743e497b32195cbc690abc06b42436"
const SecretKey string = "adf07bdd63b240ae60087efd2e72269b9c65cc91"

var tokenAuth *jwtauth.JWTAuth

func main() {

	Route.New_route()

}
