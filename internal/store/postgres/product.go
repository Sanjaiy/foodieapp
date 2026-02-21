package postgres

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/Sanjaiy/foodieapp/internal/db"
	"github.com/Sanjaiy/foodieapp/internal/domain"
)

// ProductStore is the PostgreSQL implementation of store.ProductStore.
type ProductStore struct {
	q *db.Queries
}

// NewProductStore creates a new PostgreSQL-backed product store.
func NewProductStore(dbConn *sql.DB) *ProductStore {
	return &ProductStore{q: db.New(dbConn)}
}

// ListProducts returns all available products.
func (s *ProductStore) ListProducts(ctx context.Context) ([]domain.Product, error) {
	rows, err := s.q.ListProducts(ctx)
	if err != nil {
		return nil, err
	}

	products := make([]domain.Product, len(rows))
	for i, row := range rows {
		products[i] = toProduct(row)
	}
	return products, nil
}

// GetProduct returns a single product by ID. Returns nil if not found.
func (s *ProductStore) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	row, err := s.q.GetProduct(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	p := toProduct(row)
	return &p, nil
}

// toProduct converts a sqlcdb.Product to a domain.Product.
func toProduct(row db.Product) domain.Product {
	price, _ := strconv.ParseFloat(row.Price, 64)
	return domain.Product{
		ID:       row.ID,
		Name:     row.Name,
		Price:    price,
		Category: row.Category,
		Image: &domain.ProductImage{
			Thumbnail: row.ImgThumb,
			Mobile:    row.ImgMobile,
			Tablet:    row.ImgTablet,
			Desktop:   row.ImgDesktop,
		},
	}
}
