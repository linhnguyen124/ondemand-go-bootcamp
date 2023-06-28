package api

// import (
// 	"encoding/csv"
// 	"encoding/json"
// 	"io"
// 	"log"
// 	"net/http"
// 	"os"
// 	"strconv"
// 	"testing"

// 	"github.com/gorilla/mux"
// )

// type Pokemon struct {
// 	ID   int    `json:"id"`
// 	Name string `json:"name"`
// }

// type ExternalAPIClient struct {
// 	BaseURL string
// }

// func (c *ExternalAPIClient) GetPokemonByID(id int) (*Pokemon, error) {
// 	url := c.BaseURL + "/pokemon/" + strconv.Itoa(id)
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, nil
// 	}

// 	pokemon := &Pokemon{}
// 	err = json.NewDecoder(resp.Body).Decode(pokemon)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return pokemon, nil
// }

// func SavePokemonToCSV(filePath string, pokemon *Pokemon) error {
// 	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	writer := csv.NewWriter(file)
// 	defer writer.Flush()

// 	err = writer.Write([]string{strconv.Itoa(pokemon.ID), pokemon.Name})
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

// 	client := &ExternalAPIClient{BaseURL: "https://api.example.com"}

// 	pokemon, err := client.GetPokemonByID(id)
// 	if err != nil {
// 		log.Println("Error retrieving Pokemon:", err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}

// 	if pokemon == nil {
// 		http.Error(w, "Pokemon not found", http.StatusNotFound)
// 		return
// 	}

// 	err = SavePokemonToCSV("./resources/pokemon.csv", pokemon)
// 	if err != nil {
// 		log.Println("Error saving Pokemon to CSV:", err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}

// 	jsonData, err := json.Marshal(pokemon)
// 	if err != nil {
// 		log.Println("Error marshaling Pokemon to JSON:", err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(jsonData)
// }

// func TestGetPokemonByID(t *testing.T) {
// 	// Initialize the router
// 	r := mux.NewRouter()
// 	r.HandleFunc("/pokemon/{id}", GetPokemonByID).Methods("GET")

// 	// Create a test request
// 	req, err := http.NewRequest("GET", "/pokemon/1", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Create a response recorder to record the response
// 	rr := httptest.NewRecorder()

// 	// Serve the request
// 	r.ServeHTTP(rr, req)

// 	// Check the status code
// 	if rr.Code != http.StatusOK {
// 		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
// 	}

// 	// Check the response body
// 	expected := `{"id":1,"name":"bulbasaur"}`
// 	if rr.Body.String() != expected {
// 		t.Errorf("expected response body %q, got %q", expected, rr.Body.String())
// 	}
// }
