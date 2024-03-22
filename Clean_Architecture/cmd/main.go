package main

import (
	"database/sql"
	"github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Controller"
	repository "github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Repository"
	"github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Route"
	"github.com/vadim-shalnev/swaggerApiExample/Clean_Architecture/Service"
	"log"
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
	Route.New_router(contr)

}
