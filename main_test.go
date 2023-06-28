package main_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/linhnguyen124/ondemand-go-bootcamp/api"
	"github.com/linhnguyen124/ondemand-go-bootcamp/model"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetPokemonByID(t *testing.T) {
	// Create a new router
	router := mux.NewRouter()

	// Register the GetPokemonByID handler
	router.HandleFunc("/pokemon/{id}", api.GetPokemonByID).Methods("GET")

	// Create a new request with the test URL
	req, err := http.NewRequest("GET", "/pokemon/1", nil)
	assert.NoError(t, err, "Failed to create request")

	// Create a new response recorder
	rec := httptest.NewRecorder()

	// Serve the request and record the response
	router.ServeHTTP(rec, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rec.Code, "Expected status code %d, got %d", http.StatusOK, rec.Code)

	// Decode the response body into a Pokemon struct
	var pokemon model.Pokemon
	err = json.NewDecoder(rec.Body).Decode(&pokemon)
	assert.NoError(t, err, "Failed to decode response body")

	// Check the Pokemon ID and Name
	assert.Equal(t, 1, pokemon.ID, "Expected Pokemon ID to be 1")
	assert.Equal(t, "bulbasaur", pokemon.Name, "Expected Pokemon Name to be Bulbasaur")
}

func TestGetPokemonByID_InvalidID(t *testing.T) {
	// Create a new router
	router := mux.NewRouter()

	// Register the GetPokemonByID handler
	router.HandleFunc("/pokemon/{id}", api.GetPokemonByID).Methods("GET")

	// Create a new request with an invalid ID (non-numeric)
	req, err := http.NewRequest("GET", "/pokemon/abc", nil)
	assert.NoError(t, err, "Failed to create request")

	// Create a new response recorder
	rec := httptest.NewRecorder()

	// Serve the request and record the response
	router.ServeHTTP(rec, req)

	// Check the response status code
	assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status code %d, got %d", http.StatusBadRequest, rec.Code)
}
