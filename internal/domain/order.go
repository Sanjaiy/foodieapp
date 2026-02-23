package domain

type Order struct {
	ID        string      `json:"id"`
	Items     []OrderItem `json:"items"`
	Total     float64     `json:"total"`
	Discounts float64     `json:"discounts"`
	Products  []Product   `json:"products"`
}

type OrderItem struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}
