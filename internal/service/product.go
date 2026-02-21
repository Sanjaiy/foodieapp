package service

import (
	"context"

	"github.com/Sanjaiy/foodieapp/internal/domain"
	"github.com/Sanjaiy/foodieapp/internal/store"
)

type ProductService struct {
	store store.ProductStore
}

func NewProductService(s store.ProductStore) *ProductService {
	return &ProductService{
		store: s,
	}
}

func (s *ProductService) ListProducts(ctx context.Context) ([]domain.Product, error) {
	return s.store.ListProducts(ctx)
}

func (s *ProductService) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	return s.store.GetProduct(ctx, id)
}
