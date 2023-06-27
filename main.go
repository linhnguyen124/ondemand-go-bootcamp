package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/linhnguyen124/ondemand-go-bootcamp/api"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/pokemon/{id}", api.GetPokemonByID).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", r))
}
