package main

import (
	"log"
	"os"

	"api_gateway/configs"
	"api_gateway/db"
	"api_gateway/endpoint"
)

// @title api gateway
// @version 1.0
// @description swagger for api gateway.

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Format: "Bearer <token>"
func main() {
	configs.InitLogging()

	db, err := db.NewDatabase()

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	endpoint.NewRouter(db).Router.Run(os.Getenv("URL_SERVER"))
}
