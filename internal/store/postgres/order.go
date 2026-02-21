package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/Sanjaiy/foodieapp/internal/db"
	"github.com/Sanjaiy/foodieapp/internal/domain"
)

// OrderStore is the PostgreSQL implementation of store.OrderStore.
type OrderStore struct {
	db *sql.DB
	q  *db.Queries
}

// NewOrderStore creates a new PostgreSQL-backed order store.
func NewOrderStore(dbConn *sql.DB) *OrderStore {
	return &OrderStore{
		db: dbConn,
		q:  db.New(dbConn),
	}
}

// ValidateProducts checks that all given product IDs exist and returns them.
func (s *OrderStore) ValidateProducts(ctx context.Context, productIDs []string) ([]domain.Product, error) {
	rows, err := s.q.GetProductsByIDs(ctx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("fetching products: %w", err)
	}

	// Check if all requested IDs were found
	if len(rows) != len(productIDs) {
		return nil, nil // signal that some products are missing
	}

	products := make([]domain.Product, len(rows))
	for i, row := range rows {
		products[i] = toProduct(row)
	}
	return products, nil
}

// CreateOrder persists a new order with its items in a transaction.
func (s *OrderStore) CreateOrder(ctx context.Context, items []domain.OrderItem, coupon string, products []domain.Product, total, discounts float64) (*domain.Order, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := s.q.WithTx(tx)

	// Create the order
	var couponCode sql.NullString
	if coupon != "" {
		couponCode = sql.NullString{String: coupon, Valid: true}
	}

	orderRow, err := qtx.CreateOrder(ctx, db.CreateOrderParams{
		CouponCode: couponCode,
		Total:      strconv.FormatFloat(total, 'f', 2, 64),
		Discounts:  strconv.FormatFloat(discounts, 'f', 2, 64),
	})
	if err != nil {
		return nil, fmt.Errorf("creating order: %w", err)
	}

	// Create order items
	for _, item := range items {
		err = qtx.CreateOrderItem(ctx, db.CreateOrderItemParams{
			OrderID:   orderRow.ID,
			ProductID: item.ProductID,
			Quantity:  int32(item.Quantity),
		})
		if err != nil {
			return nil, fmt.Errorf("creating order item: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("committing transaction: %w", err)
	}

	// Build the response
	order := &domain.Order{
		ID:        orderRow.ID.String(),
		Items:     items,
		Total:     total,
		Discounts: discounts,
		Products:  products,
	}

	return order, nil
}
