package model

type Product struct {
	ID          string `json:"id"`
	ProductName string `json:"product_name"`
	DrugGeneric string `json:"drug_generic"`
	Company     string `json:"company"`
}
