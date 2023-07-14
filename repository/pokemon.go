package repository

import (
	"encoding/csv"
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/linhnguyen124/ondemand-go-bootcamp/model"
)

var ErrCSVFileNotFound = errors.New("CSV file not found")

// PokemonRepository handles data operations for Pokemon
type PokemonRepository struct {
	CSVFilePath string // File path of the CSV data file
}

// NewPokemonRepository creates a new instance of PokemonRepository
func NewPokemonRepository(filePath string) *PokemonRepository {
	return &PokemonRepository{
		CSVFilePath: filePath,
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

func (r *PokemonRepository) GetPokemonByID(id int) (*model.Pokemon, error) {
	file, err := os.Open(r.CSVFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File not found error
			return nil, ErrCSVFileNotFound
		}
		// Other file-related error
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		// Error reading CSV data
		return nil, err
	}

	for _, record := range records {
		recordID, err := strconv.Atoi(record[0])
		if err != nil {
			continue
		}
		if recordID == id {
			pokemon := &model.Pokemon{
				ID:   recordID,
				Name: record[1],
			}
			return pokemon, nil
		}
	}

	return nil, nil // Pokemon not found
}
