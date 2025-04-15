package util

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/google/uuid"
	"github.com/vctrthe/ProductSearch/model"
)

func LoadAndIndexData(ES *elasticsearch.Client) {
	file, err := os.Open("data/data_product.csv")
	if err != nil {
		log.Fatalf("csv read error: %s", err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	for _, row := range records[1:] {
		product := model.Product{
			ID:          row[0],
			ProductName: row[1],
			DrugGeneric: row[2],
			Company:     row[3],
		}

		data, _ := json.Marshal(product)
		ES.Index("products", strings.NewReader(string(data)), ES.Index.WithDocumentID(uuid.New().String()))
	}

	fmt.Println("data loaded to Elasticsearch")
}
