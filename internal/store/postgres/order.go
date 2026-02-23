package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/Sanjaiy/foodieapp/internal/db"
	"github.com/Sanjaiy/foodieapp/internal/domain"
	"github.com/Sanjaiy/foodieapp/internal/store"
)

type OrderStore struct {
	db *sql.DB
	q  *db.Queries
}

func NewOrderStore(dbConn *sql.DB) *OrderStore {
	return &OrderStore{
		db: dbConn,
		q:  db.New(dbConn),
	}
}

func (s *OrderStore) ValidateProducts(ctx context.Context, productIDs []string) ([]domain.Product, error) {
	rows, err := s.q.GetProductsByIDs(ctx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("fetching products: %w", err)
	}

	if len(rows) != len(productIDs) {
		return nil, nil
	}

	products := make([]domain.Product, len(rows))
	for i, row := range rows {
		products[i] = toProduct(row)
	}

	return products, nil
}

func (s *OrderStore) CreateOrder(ctx context.Context, input store.CreateOrderInput) (*domain.Order, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := s.q.WithTx(tx)

	var validCouponCode sql.NullString
	if input.CouponCode != "" {
		validCouponCode = sql.NullString{String: input.CouponCode, Valid: true}
	}

	orderRow, err := qtx.CreateOrder(ctx, db.CreateOrderParams{
		CouponCode: validCouponCode,
		Total:      strconv.FormatFloat(input.Total, 'f', 2, 64),
		Discounts:  strconv.FormatFloat(input.Discounts, 'f', 2, 64),
	})
	if err != nil {
		return nil, fmt.Errorf("creating order: %w", err)
	}

	for _, item := range input.Items {
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

	order := &domain.Order{
		ID:        orderRow.ID.String(),
		Items:     input.Items,
		Total:     input.Total,
		Discounts: input.Discounts,
		Products:  input.Products,
	}

	return order, nil
}
