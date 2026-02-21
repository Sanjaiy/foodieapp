package dto

import "github.com/Sanjaiy/foodieapp/internal/domain"

type OrderRequest struct {
	CouponCode string      `json:"couponCode,omitempty"`
	Items      []OrderItem `json:"items"`
}

type OrderItem struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type OrderResponse struct {
	ID        string      `json:"id"`
	Items     []OrderItem `json:"items"`
	Total     float64     `json:"total"`
	Discounts float64     `json:"discounts"`
	Products  []Product   `json:"products"`
}

func FromDomainOrder(o *domain.Order) *OrderResponse {
	if o == nil {
		return nil
	}
	items := make([]OrderItem, len(o.Items))
	for i, item := range o.Items {
		items[i] = OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}
	return &OrderResponse{
		ID:        o.ID,
		Items:     items,
		Total:     o.Total,
		Discounts: o.Discounts,
		Products:  FromDomainProducts(o.Products),
	}
}

func (req OrderRequest) ToDomainItems() []domain.OrderItem {
	if req.Items == nil {
		return nil
	}
	items := make([]domain.OrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}
	return items
}
