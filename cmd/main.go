package main

import (
	"fmt"
	"log"
	"net/http"

	"notes-api/config"
	"notes-api/routes"
)

func main() {

	db := config.ConnectDB()

	r := routes.SetupRoutes(db)

	fmt.Println("Server started on :8080")

	log.Fatal(http.ListenAndServe(":8080", r))
}
