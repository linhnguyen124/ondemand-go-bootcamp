package main

import (
	"log"
	"net/http"

	"github.com/linhnguyen124/ondemand-go-bootcamp/api"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/pokemon/{id}", api.GetPokemonByID).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
