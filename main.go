package main

import (
	"log"
	"net/http"

	"github.com/linhnguyen124/ondemand-go-bootcamp/api"
	"github.com/linhnguyen124/ondemand-go-bootcamp/repository"
	"github.com/linhnguyen124/ondemand-go-bootcamp/service"

	"github.com/gorilla/mux"
)

func main() {
	// Create repositories
	pokemonRepo := repository.NewPokemonRepository("./resources/data.csv")

	// Create services
	pokemonService := service.NewPokemonService(pokemonRepo)

	// Create handlers
	pokemonHandler := api.NewPokemonHandler(pokemonService)

	// Create a new router
	router := mux.NewRouter()

	// Register the endpoints with the router
	router.HandleFunc("/pokemon/{id}", pokemonHandler.GetPokemonByID).Methods("GET")
	router.HandleFunc("/pokemon", pokemonHandler.ConcurrentReadFromCSV).Methods("POST")

	// Start the server with the router
	log.Println("Server started on port 8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Server startup failed:", err)
	}
}
