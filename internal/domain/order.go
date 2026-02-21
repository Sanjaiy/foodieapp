package domain

type Order struct {
	ID        string
	Items     []OrderItem
	Total     float64
	Discounts float64
	Products  []Product
}

type OrderItem struct {
	ProductID string
	Quantity  int
}
