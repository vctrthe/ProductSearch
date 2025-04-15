package model

type SearchResult struct {
	ID          string  `json:"id"`
	ProductName string  `json:"product_name"`
	DrugGeneric string  `json:"drug_generic"`
	Company     string  `json:"company"`
	Score       float64 `json:"score"`
}

type SearchResponse struct {
	Results []SearchResult `json:"results"`
}
