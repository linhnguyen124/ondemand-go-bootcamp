package service

import (
	"errors"
	"log"
	"sync"

	"github.com/linhnguyen124/ondemand-go-bootcamp/model"
	"github.com/linhnguyen124/ondemand-go-bootcamp/repository"
)

// PokemonService handles business logic for Pokemon
type PokemonService struct {
	pokemonRepo *repository.PokemonRepository
}

// NewPokemonService creates a new instance of PokemonService
func NewPokemonService(pokemonRepo *repository.PokemonRepository) *PokemonService {
	return &PokemonService{
		pokemonRepo: pokemonRepo,
	}
}

// GetPokemonListConcurrently retrieves a list of Pokemon concurrently based on the specified type and item count
func (s *PokemonService) GetPokemonListConcurrently(pokemonType string, items, itemsPerWorker int) ([]*model.Pokemon, error) {
	pokemonList, err := s.pokemonRepo.ReadAll()
	if err != nil {
		return nil, err
	}

	// Create a wait group to synchronize the worker goroutines
	var wg sync.WaitGroup

	// Create a channel to receive the processed items
	results := make(chan *model.Pokemon)

	// Launch workers to process the items concurrently
	numWorkers := (items + itemsPerWorker - 1) / itemsPerWorker // Round up the number of workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go s.worker(i, itemsPerWorker, pokemonType, pokemonList, results, &wg)
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
	pokemonList, ok := <-completedItems
	if !ok {
		return nil, errors.New("failed to retrieve completed items")
	}
	return pokemonList, nil
}

func (s *PokemonService) worker(id, itemsPerWorker int, pokemonType string, pokemonList []*model.Pokemon, results chan<- *model.Pokemon, wg *sync.WaitGroup) {
	defer wg.Done()

	// Process each record assigned to the worker
	for _, pokemon := range pokemonList {
		// Add log statement to indicate worker start
		log.Printf("Worker %d: Start processing", id)
		if len(results) >= itemsPerWorker {
			return
		}

		pokemonID := pokemon.ID

		// Check if the Pokemon type matches the specified type (odd or even)
		if (pokemonType == "odd" && pokemonID%2 == 0) || (pokemonType == "even" && pokemonID%2 == 1) {
			continue
		}

		results <- pokemon
		// Add log statement to indicate worker finish
		log.Printf("Worker %d: Finish processing", id)
	}
}
