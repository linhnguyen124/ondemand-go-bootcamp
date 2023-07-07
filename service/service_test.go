package service

import (
	"testing"

	"github.com/linhnguyen124/ondemand-go-bootcamp/model"
	"github.com/linhnguyen124/ondemand-go-bootcamp/repository"
	"github.com/stretchr/testify/assert"
)

func TestGetPokemonListConcurrently(t *testing.T) {
	// Create a new PokemonRepository with the test data file
	pokemonRepo := repository.NewPokemonRepository()
	pokemonRepo.CSVFilePath = "../testdata/test_data.csv"

	// Create a new PokemonService using the PokemonRepository
	pokemonService := NewPokemonService(pokemonRepo)

	// Test case parameters
	pokemonType := "odd"
	items := 3
	itemsPerWorker := 2

	// Call the GetPokemonListConcurrently method
	pokemonList, err := pokemonService.GetPokemonListConcurrently(pokemonType, items, itemsPerWorker)

	// Check for any errors
	assert.NoError(t, err)

	// Check the length of the returned Pokemon list
	assert.Equal(t, items, len(pokemonList))

	// Expected Pokemon list
	expectedPokemonList := []*model.Pokemon{
		{ID: 1, Name: "bulbasaur"},
		{ID: 3, Name: "venusaur"},
		{ID: 5, Name: "charmeleon"},
	}

	// Check each Pokemon in the list
	for i := 0; i < len(pokemonList); i++ {
		assert.Equal(t, expectedPokemonList[i].ID, pokemonList[i].ID)
		assert.Equal(t, expectedPokemonList[i].Name, pokemonList[i].Name)
	}
}
