package api

// import (
// 	"encoding/csv"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"strconv"

// 	"github.com/gorilla/mux"
// )

// type Pokemon struct {
// 	ID   int    `json:"id"`
// 	Name string `json:"name"`
// }

// type PokemonAPI struct {
// 	BaseURL string
// }

// func NewPokemonAPI(baseURL string) *PokemonAPI {
// 	return &PokemonAPI{BaseURL: baseURL}
// }

// func (c *PokemonAPI) GetPokemonByID(id int) (*Pokemon, error) {
// 	url := fmt.Sprintf("%s/pokemon/%d", c.BaseURL, id)
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("failed to retrieve Pokemon: %s", resp.Status)
// 	}

// 	pokemon := &Pokemon{}
// 	err = json.NewDecoder(resp.Body).Decode(pokemon)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return pokemon, nil
// }

// func StorePokemonInCSV(pokemon *Pokemon, filePath string) error {
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

// func GetPokemonByID(w http.ResponseWriter, r *http.Request) {
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

// 	err = StorePokemonInCSV(pokemon, "./resources/data.csv")
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
// 	r.HandleFunc("/pokemon/{id}", GetPokemonByID).Methods("GET")
// 	log.Fatal(http.ListenAndServe(":8080", r))
// }
