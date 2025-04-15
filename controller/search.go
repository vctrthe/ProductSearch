package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vctrthe/ProductSearch/config"
	"github.com/vctrthe/ProductSearch/model"
)

func SearchProduct(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing query param ?q="})
		return
	}

	var buf bytes.Buffer
	words := strings.Fields(keyword)

	var query map[string]interface{}

	if len(words) == 1 {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"multi_match": map[string]interface{}{
					"query":     keyword,
					"fields":    []string{"product_name^5", "drug_generic^2", "company^1"},
					"fuzziness": "AUTO",
				},
			},
		}
	} else {
		mustQueries := []map[string]interface{}{}
		for _, word := range words {
			mustQueries = append(mustQueries, map[string]interface{}{
				"multi_match": map[string]interface{}{
					"query":     word,
					"fields":    []string{"product_name^5", "drug_generic^2", "company^1"},
					"fuzziness": "AUTO",
				},
			})
		}

		query = map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": mustQueries,
				},
			},
		}
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("error encoding query: %s", err)
	}

	res, err := config.ES.Search(
		config.ES.Search.WithContext(context.Background()),
		config.ES.Search.WithIndex("products"),
		config.ES.Search.WithBody(&buf),
		config.ES.Search.WithTrackTotalHits(true),
		config.ES.Search.WithSort("_score:asc"),
		config.ES.Search.WithPretty(),
	)

	if err != nil {
		log.Fatalf("search error: %s", err)
	}

	defer res.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)

	var hits []model.SearchResult
	for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {
		doc := hit.(map[string]interface{})
		source := doc["_source"].(map[string]interface{})
		score := doc["_score"].(float64)

		hits = append(hits, model.SearchResult{
			ID:          source["id"].(string),
			ProductName: source["product_name"].(string),
			DrugGeneric: source["drug_generic"].(string),
			Company:     source["company"].(string),
			Score:       score,
		})
	}

	c.JSON(http.StatusOK, model.SearchResponse{Results: hits})
}
