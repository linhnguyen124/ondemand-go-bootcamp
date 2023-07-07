package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/linhnguyen124/ondemand-go-bootcamp/model"
	"github.com/linhnguyen124/ondemand-go-bootcamp/service"
)

type PokemonResponse struct {
	PokemonList []*model.Pokemon `json:"pokemonList"`
}

type PokemonHandler struct {
	pokemonService *service.PokemonService
}

func NewPokemonHandler(pokemonService *service.PokemonService) *PokemonHandler {
	return &PokemonHandler{
		pokemonService: pokemonService,
	}
}

// ConcurrentReadFromCSV is the handler for the "/pokemon" endpoint.
// It reads items from the CSV file concurrently using a worker pool.
func (h *PokemonHandler) ConcurrentReadFromCSV(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	pokemonType := r.URL.Query().Get("type")
	itemsParam := r.URL.Query().Get("items")
	itemsPerWorkerParam := r.URL.Query().Get("items_per_worker")

	// Validate and convert query parameters
	items, err := strconv.Atoi(itemsParam)
	if err != nil || items <= 0 {
		http.Error(w, "Invalid items parameter", http.StatusBadRequest)
		return
	}

	itemsPerWorker, err := strconv.Atoi(itemsPerWorkerParam)
	if err != nil || itemsPerWorker <= 0 {
		http.Error(w, "Invalid items_per_worker parameter", http.StatusBadRequest)
		return
	}

	pokemonList, err := h.pokemonService.GetPokemonListConcurrently(pokemonType, items, itemsPerWorker)
	if err != nil {
		log.Println("Error retrieving Pokemon list:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create a response object
	response := PokemonResponse{
		PokemonList: pokemonList,
	}

	// Encode the response as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding JSON response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
