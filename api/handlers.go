package api

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/linhnguyen124/ondemand-go-bootcamp/model"

	"github.com/gorilla/mux"
)

// PokemonAPI represents the external client for interacting with API.
type PokemonAPI struct {
	BaseURL string
}

// NewPokemonAPI creates a new instance of PokemonAPI with the specified base URL.
func NewPokemonAPI(baseURL string) *PokemonAPI {
	return &PokemonAPI{BaseURL: baseURL}
}

// GetPokemonByID is the handler for the "/pokemon/{id}" endpoint.
// It retrieves the Pokemon data from the CSV file if it exists,
// otherwise, it calls the external API to fetch the data, stores it in the CSV file,
// and returns the result as JSON in the HTTP response.
func GetPokemonByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Check if the Pokemon data exists in the CSV file
	pokemon, err := GetPokemonFromCSV(id, "./resources/data.csv")
	if err == nil {
		// If the data exists in the CSV file, return it as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pokemon)
		return
	}

	// Create the PokemonAPI client with the base URL of the external API
	api := NewPokemonAPI("https://pokeapi.co/api/v2")

	// Fetch the Pokemon data from the external API
	pokemon, err = api.GetPokemonByID(id)
	if err != nil {
		log.Println("Error retrieving Pokemon:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Store the Pokemon data in the CSV file
	err = StorePokemonInCSV(pokemon, "./resources/data.csv")
	if err != nil {
		log.Println("Error storing Pokemon data in CSV:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Display the result as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pokemon)
}

// GetPokemonFromCSV retrieves the Pokemon data from the CSV file based on the provided ID.
func GetPokemonFromCSV(id int, filePath string) (*model.Pokemon, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		recordID, err := strconv.Atoi(record[0])
		if err != nil {
			log.Println("Error parsing ID:", err)
			continue
		}

		if recordID == id {
			pokemon := &model.Pokemon{
				ID:   recordID,
				Name: record[1],
			}
			return pokemon, nil
		}
	}

	return nil, fmt.Errorf("pokemon not found in CSV")
}

// GetPokemonByID fetches the Pokemon data from the external API based on the provided ID.
func (c *PokemonAPI) GetPokemonByID(id int) (*model.Pokemon, error) {
	url := fmt.Sprintf("%s/pokemon/%d", c.BaseURL, id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve Pokemon: %s", resp.Status)
	}

	pokemon := &model.Pokemon{}
	err = json.NewDecoder(resp.Body).Decode(pokemon)
	if err != nil {
		return nil, err
	}

	return pokemon, nil
}

// StorePokemonInCSV stores the Pokemon data in the CSV file.
func StorePokemonInCSV(pokemon *model.Pokemon, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{strconv.Itoa(pokemon.ID), pokemon.Name}
	err = writer.Write(record)
	if err != nil {
		return err
	}

	return nil
}
