package util

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/vctrthe/ProductSearch/model"
)

func LoadAndIndexData(ES *elasticsearch.Client) {
	file, err := os.Open("data/data_product.csv")
	if err != nil {
		log.Fatalf("csv read error: %s", err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("error reading CSV: %s", err)
	}

	// First, delete the existing index to avoid duplicates
	_, err = ES.Indices.Delete([]string{"products"})
	if err != nil {
		log.Printf("warning: could not delete existing index: %s", err)
	}

	// Create the index with proper mappings
	_, err = ES.Indices.Create("products")
	if err != nil {
		log.Printf("warning: could not create index: %s", err)
	}

	// Index the data
	for _, row := range records[1:] {
		if len(row) < 4 {
			log.Printf("warning: skipping invalid row: %v", row)
			continue
		}

		product := model.Product{
			ID:          strings.TrimSpace(row[0]),
			ProductName: strings.TrimSpace(row[1]),
			DrugGeneric: strings.TrimSpace(row[2]),
			Company:     strings.TrimSpace(row[3]),
		}

		// Skip if ID is empty
		if product.ID == "" {
			log.Printf("warning: skipping row with empty ID: %v", row)
			continue
		}

		data, err := json.Marshal(product)
		if err != nil {
			log.Printf("error marshaling product: %s", err)
			continue
		}

		// Use product ID as document ID to ensure uniqueness
		_, err = ES.Index(
			"products",
			strings.NewReader(string(data)),
			ES.Index.WithDocumentID(product.ID),
		)
		if err != nil {
			log.Printf("error indexing product %s: %s", product.ID, err)
			continue
		}
	}

	// Refresh the index to make all operations visible
	_, err = ES.Indices.Refresh(ES.Indices.Refresh.WithIndex("products"))
	if err != nil {
		log.Printf("warning: could not refresh index: %s", err)
	}

	fmt.Println("data loaded to Elasticsearch")
}
