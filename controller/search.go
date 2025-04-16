package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vctrthe/ProductSearch/config"
	"github.com/vctrthe/ProductSearch/model"
)

// @Summary Search for products by keyword or full-text
// @Description Search for products using Elasticsearch. Results are sorted by score in ascending order. (Lowest to highest score)
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Success 200 {object} model.SearchResponse
// @Failure 400 {object} object "Bad Request"
// @Router /search [get]
func SearchProduct(c *gin.Context) {
	// Get the search query from URL parameter
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing query param ?q="})
		return
	}

	// Create a timeout context to prevent long-running searches
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	// Prepare the search query
	var searchQuery map[string]interface{}
	words := strings.Fields(query)

	if len(words) == 1 {
		// Single word search
		searchQuery = map[string]interface{}{
			"query": map[string]interface{}{
				"multi_match": map[string]interface{}{
					"query":     query,
					"fields":    []string{"product_name^5", "drug_generic^3", "company^2"},
					"fuzziness": "AUTO",
				},
			},
			"size": 10, // Return only top 10 results
			"sort": []map[string]interface{}{
				{"_score": "asc"}, // Sort by score in ascending order
			},
		}
	} else {
		// Multiple words search
		mustQueries := []map[string]interface{}{}
		for _, word := range words {
			mustQueries = append(mustQueries, map[string]interface{}{
				"multi_match": map[string]interface{}{
					"query":     word,
					"fields":    []string{"product_name^5", "drug_generic^3", "company^2"},
					"fuzziness": "AUTO",
				},
			})
		}

		searchQuery = map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": mustQueries,
				},
			},
			"size": 10, // Return only top 10 results
			"sort": []map[string]interface{}{
				{"_score": "asc"}, // Sort by score in ascending order
			},
		}
	}

	// Convert query to JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		log.Printf("error encoding query: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Execute the search
	res, err := config.ES.Search(
		config.ES.Search.WithContext(ctx),
		config.ES.Search.WithIndex("products"),
		config.ES.Search.WithBody(&buf),
		config.ES.Search.WithTrackTotalHits(true),
	)

	if err != nil {
		log.Printf("search error: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	defer res.Body.Close()

	// Parse the response
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Printf("error decoding response: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Extract hits and total count
	hits := result["hits"].(map[string]interface{})
	_ = hits["total"].(map[string]interface{})["value"].(float64)
	hitsArray := hits["hits"].([]interface{})

	// Convert hits to our response format
	var searchResults []model.SearchResult
	for _, hit := range hitsArray {
		doc := hit.(map[string]interface{})
		source := doc["_source"].(map[string]interface{})
		score := doc["_score"].(float64)

		searchResults = append(searchResults, model.SearchResult{
			ID:          source["id"].(string),
			ProductName: source["product_name"].(string),
			DrugGeneric: source["drug_generic"].(string),
			Company:     source["company"].(string),
			Score:       score,
		})
	}

	// Return the response
	c.JSON(http.StatusOK, model.SearchResponse{
		Results: searchResults,
	})
}
