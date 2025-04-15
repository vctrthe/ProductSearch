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
					"fields":    []string{"product_name", "drug_generic", "company"},
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
					"fields":    []string{"product_name", "drug_generic", "company"},
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
	)

	if err != nil {
		log.Fatalf("search error: %s", err)
	}

	defer res.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)

	var hits []map[string]interface{}
	for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {
		doc := hit.(map[string]interface{})
		source := doc["_source"].(map[string]interface{})
		source["score"] = doc["_score"]
		hits = append(hits, source)
	}

	c.JSON(http.StatusOK, gin.H{"results": hits})
}
