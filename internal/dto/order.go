package dto

import "github.com/Sanjaiy/foodieapp/internal/domain"

type OrderRequest struct {
	CouponCode string             `json:"couponCode,omitempty"`
	Items      []domain.OrderItem `json:"items"`
}

type OrderResponse struct {
	ID        string             `json:"id"`
	Items     []domain.OrderItem `json:"items"`
	Total     float64            `json:"total"`
	Discounts float64            `json:"discounts"`
	Products  []domain.Product   `json:"products"`
}

func FromDomainOrder(o *domain.Order) *OrderResponse {
	if o == nil {
		return nil
	}
	return &OrderResponse{
		ID:        o.ID,
		Items:     o.Items,
		Total:     o.Total,
		Discounts: o.Discounts,
		Products:  o.Products,
	}
}

func (req OrderRequest) ToDomainItems() []domain.OrderItem {
	return req.Items
}
