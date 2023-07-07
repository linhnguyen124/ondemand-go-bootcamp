package repository

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"

	"github.com/linhnguyen124/ondemand-go-bootcamp/model"
)

// PokemonRepository handles data operations for Pokemon
type PokemonRepository struct {
	CSVFilePath string // File path of the CSV data file
}

// NewPokemonRepository creates a new instance of PokemonRepository
func NewPokemonRepository() *PokemonRepository {
	return &PokemonRepository{
		CSVFilePath: "./resources/data.csv",
	}
}

// ReadAll reads all Pokemon records from the CSV file
func (r *PokemonRepository) ReadAll() ([]*model.Pokemon, error) {
	file, err := os.Open(r.CSVFilePath)
	if err != nil {
		log.Println("Error opening CSV file:", err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		log.Println("Error reading CSV:", err)
		return nil, err
	}

	pokemonList := make([]*model.Pokemon, 0, len(records))
	for _, record := range records {
		pokemonID, err := strconv.Atoi(record[0])
		if err != nil {
			log.Println("Error parsing Pokemon ID:", err)
			continue
		}

		pokemonName := record[1]

		pokemon := &model.Pokemon{
			ID:   pokemonID,
			Name: pokemonName,
		}

		pokemonList = append(pokemonList, pokemon)
	}

	return pokemonList, nil
}
