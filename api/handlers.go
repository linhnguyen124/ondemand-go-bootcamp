package api

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/linhnguyen124/ondemand-go-bootcamp/model"
)

// CSVFilePath is the file path of the CSV data file.
var CSVFilePath string

type PokemonResponse struct {
	PokemonList []*model.Pokemon `json:"pokemonList"`
}

// ConcurrentReadFromCSV is the handler for the "/pokemon" endpoint.
// It reads items from the CSV file concurrently using a worker pool.
func ConcurrentReadFromCSV(w http.ResponseWriter, r *http.Request) {
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

	// Open the CSV file for reading
	CSVFilePath := "./resources/data.csv"
	file, err := os.Open(CSVFilePath)
	if err != nil {
		log.Println("Error opening CSV file:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Read records from CSV
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		log.Println("Error reading CSV:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create a wait group to synchronize the worker goroutines
	var wg sync.WaitGroup

	// Create a channel to receive the processed items
	results := make(chan *model.Pokemon)

	// Launch workers to process the items concurrently
	numWorkers := (items + itemsPerWorker - 1) / itemsPerWorker // Round up the number of workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, itemsPerWorker, pokemonType, records, results, &wg)
	}

	// Create a channel to collect the completed items
	completedItems := make(chan []*model.Pokemon)

	// Collect the processed items from the results channel
	go func() {
		pokemonList := make([]*model.Pokemon, 0, items)
		for pokemon := range results {
			pokemonList = append(pokemonList, pokemon)
			if len(pokemonList) >= items {
				break
			}
		}
		completedItems <- pokemonList
	}()

	// Wait for the completed items or worker limit
	select {
	case pokemonList := <-completedItems:
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

	case <-r.Context().Done():
		log.Println("Request canceled:", r.Context().Err())
	}
}

func worker(id, itemsPerWorker int, pokemonType string, records [][]string, results chan<- *model.Pokemon, wg *sync.WaitGroup) {
	defer wg.Done()

	// Process each record assigned to the worker
	for _, record := range records {
		if len(results) >= itemsPerWorker {
			return
		}

		pokemonID, err := strconv.Atoi(record[0])
		if err != nil {
			log.Println("Error parsing Pokemon ID:", err)
			continue
		}

		pokemonName := record[1]

		// Check if the Pokemon type matches the specified type (odd or even)
		if (pokemonType == "odd" && pokemonID%2 == 0) || (pokemonType == "even" && pokemonID%2 == 1) {
			continue
		}

		pokemon := &model.Pokemon{
			ID:   pokemonID,
			Name: pokemonName,
		}

		results <- pokemon
	}
}
