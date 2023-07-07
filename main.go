package main

import (
	"log"
	"net/http"

	"github.com/linhnguyen124/ondemand-go-bootcamp/api"
	"github.com/linhnguyen124/ondemand-go-bootcamp/repository"
	"github.com/linhnguyen124/ondemand-go-bootcamp/service"
)

func main() {
	// Create repositories
	pokemonRepo := repository.NewPokemonRepository()

	// Create services
	pokemonService := service.NewPokemonService(pokemonRepo)

	// Create handlers
	pokemonHandler := api.NewPokemonHandler(pokemonService)

	// Register handlers
	http.HandleFunc("/pokemon", pokemonHandler.ConcurrentReadFromCSV)

	// Start the server
	log.Println("Server started on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server startup failed:", err)
	}
}
