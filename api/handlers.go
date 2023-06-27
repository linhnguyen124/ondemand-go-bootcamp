package api

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/linhnguyen124/ondemand-go-bootcamp/model"
)

func GetPokemonByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	pokemonList, err := readPokemonCSV("./resources/pokemon.csv")
	if err != nil {
		log.Println("Error reading CSV:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	for _, pokemon := range pokemonList {
		if pokemon.ID == id {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(pokemon.Name))
			return
		}
	}

	http.Error(w, "Pokemon not found", http.StatusNotFound)
}

func readPokemonCSV(filePath string) ([]*model.Pokemon, error) {
	//Read the CSV file and return the data as a slice of Pokemon structs

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

	pokemonList := []*model.Pokemon{}
	for _, record := range records {
		id, err := strconv.Atoi(record[0]) // Assuming ID is at index 0
		if err != nil {
			log.Println("Error parsing ID:", err)
			continue
		}
		pokemon := &model.Pokemon{
			ID:   id,
			Name: record[1], // Assuming Name is at index 1
		}
		pokemonList = append(pokemonList, pokemon)
	}

	return pokemonList, nil
}
