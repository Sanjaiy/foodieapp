package service

import (
	"context"
	"fmt"
	"log"
	"math"

	"github.com/Sanjaiy/foodieapp/internal/domain"
	"github.com/Sanjaiy/foodieapp/internal/store"
)

const discountPercent = 10.0

type OrderService struct {
	store store.OrderStore
	promo *PromoService
}

func NewOrderService(s store.OrderStore, promo *PromoService) *OrderService {
	return &OrderService{
		store: s,
		promo: promo,
	}
}

func (s *OrderService) PlaceOrder(ctx context.Context, items []domain.OrderItem, couponCode string) (*domain.Order, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("at least one item is required")
	}

	for _, item := range items {
		if item.ProductID == "" {
			return nil, fmt.Errorf("productId is required for each item")
		}
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("quantity must be greater than 0")
		}
	}

	productIDSet := make(map[string]struct{})
	for _, item := range items {
		productIDSet[item.ProductID] = struct{}{}
	}
	productIDs := make([]string, 0, len(productIDSet))
	for id := range productIDSet {
		productIDs = append(productIDs, id)
	}

	products, err := s.store.ValidateProducts(ctx, productIDs)
	if err != nil {
		log.Printf("ERROR: validating products: %v", err)
		return nil, fmt.Errorf("failed to validate products")
	}
	if products == nil {
		return nil, fmt.Errorf("invalid product specified")
	}

	priceMap := make(map[string]float64, len(products))
	for _, p := range products {
		priceMap[p.ID] = p.Price
	}

	var total float64
	for _, item := range items {
		total += priceMap[item.ProductID] * float64(item.Quantity)
	}
	total = math.Round(total*100) / 100

	var discounts float64
	if couponCode != "" && s.promo.ValidateCoupon(couponCode) {
		discounts = math.Round(total*discountPercent) / 100
		total = math.Round((total-discounts)*100) / 100
	}

	order, err := s.store.CreateOrder(ctx, items, couponCode, products, total, discounts)
	if err != nil {
		log.Printf("ERROR: creating order: %v", err)
		return nil, fmt.Errorf("failed to create order")
	}

	return order, nil
}
