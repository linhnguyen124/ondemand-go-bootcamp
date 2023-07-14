package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/linhnguyen124/ondemand-go-bootcamp/model"
	"github.com/linhnguyen124/ondemand-go-bootcamp/repository"
	"github.com/linhnguyen124/ondemand-go-bootcamp/service"

	"github.com/gorilla/mux"
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

type HTTPHandler struct {
	repo *repository.PokemonRepository
}

func NewHTTPHandler(repo *repository.PokemonRepository) *HTTPHandler {
	return &HTTPHandler{repo: repo}
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

func (h *PokemonHandler) GetPokemonByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	pokemon, err := h.pokemonService.GetPokemonByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if pokemon == nil {
		http.Error(w, "Pokemon not found", http.StatusNotFound)
		return
	}

	jsonResponse, err := json.Marshal(pokemon)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (h *PokemonHandler) ConsumeExternalAPI(w http.ResponseWriter, r *http.Request) {
	// Make a request to the external API
	response, err := http.Get("https://pokeapi.co/api/v2/pokemon")
	if err != nil {
		http.Error(w, "Failed to fetch data from external API", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	// Parse the response JSON
	var apiResponse struct {
		Results []struct {
			Name string `json:"name"`
		} `json:"results"`
	}
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	if err != nil {
		http.Error(w, "Failed to parse data from external API", http.StatusInternalServerError)
		return
	}

	// Convert the API response to the Pokemon struct
	pokemonList := make([]model.Pokemon, len(apiResponse.Results))
	for i, result := range apiResponse.Results {
		pokemonList[i] = model.Pokemon{
			ID:   i + 1,
			Name: result.Name,
		}
	}

	// Store the data in the CSV file
	err = h.pokemonService.WritePokemonCSV(pokemonList)
	if err != nil {
		http.Error(w, "Failed to store data in CSV file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data stored in CSV file successfully"))
}
