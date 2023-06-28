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

// func NewPokemonAPI(baseURL string) *model.PokemonAPI {
// 	return &model.PokemonAPI{BaseURL: baseURL}
// }

// func (c *model.PokemonAPI) GetPokemonByID(id int) (*model.Pokemon, error) {
// 	url := fmt.Sprintf("%s/pokemon/%d", c.BaseURL, id)
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("failed to retrieve Pokemon: %s", resp.Status)
// 	}

// 	pokemon := &model.Pokemon{}
// 	err = json.NewDecoder(resp.Body).Decode(pokemon)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return pokemon, nil
// }

// func StorePokemonInCSV(pokemon *model.Pokemon, filePath string) error {
// 	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	writer := csv.NewWriter(file)
// 	defer writer.Flush()

// 	record := []string{strconv.Itoa(pokemon.ID), pokemon.Name}
// 	err = writer.Write(record)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func GetPokemonByIDByExternal(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	idStr := vars["id"]

// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
// 		return
// 	}

// 	api := NewPokemonAPI("https://pokeapi.co/api/v2")
// 	pokemon, err := api.GetPokemonByID(id)
// 	if err != nil {
// 		log.Println("Error retrieving Pokemon:", err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}

// 	err = StorePokemonInCSV(pokemon, "./resources/pokemon.csv")
// 	if err != nil {
// 		log.Println("Error storing Pokemon data in CSV:", err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}

// 	// Display the result as JSON
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(pokemon)
// }

// func main() {
// 	r := mux.NewRouter()
// 	r.HandleFunc("/pokemon/{id}", GetPokemonByIDByExternal).Methods("GET")
// 	log.Fatal(http.ListenAndServe(":8080", r))
// }
