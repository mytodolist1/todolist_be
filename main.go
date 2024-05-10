package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/mytodolist1/todolist_be/router"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := router.Router()

	fmt.Println("Server is running on port 8081")
	fmt.Println("Local : http://localhost:8081")

	err = http.ListenAndServe(":8081", router)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
