package main

import (
	"database/sql"
	"github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Controller"
	repository "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Repository"
	"github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Router"
	"github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Service"
	"log"
	"net/http"
)

func main() {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres dbname=mydatabase sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	repo := repository.NewRepositoryImpl(db)
	serv := Service.NewUserServiceImpl(repo)
	contr := Controller.NewController(serv)
	router := Router.New_router(contr)

	http.ListenAndServe(":8080", router)

}
