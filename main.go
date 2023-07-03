package main

import (
	"log"
	"net/http"

	"github.com/linhnguyen124/ondemand-go-bootcamp/api"
)

func main() {
	http.HandleFunc("/pokemon", api.ConcurrentReadFromCSV)

	log.Println("Server started on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server startup failed:", err)
	}
}
