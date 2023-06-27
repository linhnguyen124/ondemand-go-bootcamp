package api

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type Pokemon struct {
	ID   int
	Name string
}

func GetPokemonByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	file, err := os.Open("../resources/pokemon.csv")
	if err != nil {
		log.Println("Error opening CSV file:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		log.Println("Error reading CSV:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	pokemonList := []*Pokemon{}
	for _, record := range records {
		id, err := strconv.Atoi(record[0])
		if err != nil {
			log.Println("Error parsing ID:", err)
			continue
		}
		pokemon := &Pokemon{
			ID:   id,
			Name: record[1],
		}
		pokemonList = append(pokemonList, pokemon)
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
