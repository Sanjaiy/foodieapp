package store

import (
	"context"

	"github.com/Sanjaiy/foodieapp/internal/domain"
)

// ProductStore defines the interface for product data access.
// Implementations can use PostgreSQL, in-memory, or any other backend.
type ProductStore interface {
	ListProducts(ctx context.Context) ([]domain.Product, error)
	GetProduct(ctx context.Context, id string) (*domain.Product, error)
}

// OrderStore defines the interface for order data access.
type OrderStore interface {
	// ValidateProducts checks that all given product IDs exist and returns them.
	ValidateProducts(ctx context.Context, productIDs []string) ([]domain.Product, error)
	// CreateOrder persists a new order and returns it with its generated ID.
	CreateOrder(ctx context.Context, items []domain.OrderItem, couponCode string, products []domain.Product, total, discounts float64) (*domain.Order, error)
}
