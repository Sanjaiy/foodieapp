package store

import (
	"context"

	"github.com/Sanjaiy/foodieapp/internal/domain"
)

type ProductStore interface {
	ListProducts(ctx context.Context) ([]domain.Product, error)
	GetProduct(ctx context.Context, id string) (*domain.Product, error)
}

type CreateOrderInput struct {
	Items      []domain.OrderItem
	CouponCode string
	Products   []domain.Product
	Total      float64
	Discounts  float64
}

type OrderStore interface {
	ValidateProducts(ctx context.Context, productIDs []string) ([]domain.Product, error)
	CreateOrder(ctx context.Context, input CreateOrderInput) (*domain.Order, error)
}
