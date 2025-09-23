package main

import (
	"log"
	"net/http"
	"practice-go-crud/database"
	"practice-go-crud/routes"
)

func main() {
	database.ConnectDB()

	routes.SetupRoutes()

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
