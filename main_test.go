package main_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/linhnguyen124/ondemand-go-bootcamp/api"
	"github.com/linhnguyen124/ondemand-go-bootcamp/model"
	"github.com/stretchr/testify/assert"
)

type PokemonResponse struct {
	PokemonList []*model.Pokemon `json:"pokemonList"`
}

func TestConcurrentReadFromCSV(t *testing.T) {
	// Create a test request with query parameters
	req, err := http.NewRequest("GET", "/pokemon?type=odd&items=10&items_per_worker=5", nil)
	assert.NoError(t, err)

	// Create a response recorder to capture the API response
	recorder := httptest.NewRecorder()

	// Call the handler function
	api.ConcurrentReadFromCSV(recorder, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Check the response content type
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	// Parse the response body
	var response api.PokemonResponse
	err = json.NewDecoder(recorder.Body).Decode(&response)
	assert.NoError(t, err)

	// Check the number of items in the response
	assert.Equal(t, 10, len(response.PokemonList))

	// Check if the IDs of the Pokemon are odd
	for _, pokemon := range response.PokemonList {
		assert.True(t, pokemon.ID%2 == 1)
	}
}
