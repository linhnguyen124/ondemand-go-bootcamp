package api

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/linhnguyen124/ondemand-go-bootcamp/model"

	"github.com/gorilla/mux"
)

// GetPokemonByID handles an HTTP request to retrieve the name of a Pokemon based on the given ID.
// It expects the ID to be provided as a path parameter in the request URL.
// If the ID is valid and corresponds to a Pokemon in the CSV data, it returns the Pokemon's name with a status code of 200.
// If the ID is invalid or no Pokemon with the given ID is found, it returns an appropriate error response with the corresponding status code.
func GetPokemonByID(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the path parameter in the request URL
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert the ID string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Read the Pokemon data from the CSV file
	pokemonList, err := readPokemonCSV("./resources/pokemon.csv")
	if err != nil {
		log.Println("Error reading CSV:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Find the Pokemon with the given ID and return its name
	for _, pokemon := range pokemonList {
		if pokemon.ID == id {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(pokemon.Name))
			return
		}
	}

	// No Pokemon found with the given ID
	http.Error(w, "Pokemon not found", http.StatusNotFound)
}

func readPokemonCSV(filePath string) ([]model.Pokemon, error) {
	// Read the CSV file and return the data as a slice of Pokemon structs

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

	pokemonList := make([]model.Pokemon, 0, len(records))
	for _, record := range records {
		id, err := strconv.Atoi(record[0]) // Assuming ID is at index 0
		if err != nil {
			log.Println("Error parsing ID:", err)
			continue
		}
		pokemon := model.Pokemon{
			ID:   id,
			Name: record[1], // Assuming Name is at index 1
		}
		pokemonList = append(pokemonList, pokemon)
	}

	return pokemonList, nil
}
