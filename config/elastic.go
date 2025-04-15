package config

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

var ES *elasticsearch.Client

func InitElastic() {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://192.168.100.22:9200"},
		Username:  "elastic",
		Password:  "password",
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("elastic init error: %s", err)
	}

	_, err = client.Info()
	if err != nil {
		log.Fatalf("elasticsearch connection failed: %s", err)
	}

	log.Println("elasticsearch connection established")
	ES = client
}
