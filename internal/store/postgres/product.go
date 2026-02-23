package postgres

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/Sanjaiy/foodieapp/internal/db"
	"github.com/Sanjaiy/foodieapp/internal/domain"
)

type ProductStore struct {
	q *db.Queries
}

func NewProductStore(dbConn *sql.DB) *ProductStore {
	return &ProductStore{q: db.New(dbConn)}
}

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
