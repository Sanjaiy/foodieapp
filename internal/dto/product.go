package dto

import "github.com/Sanjaiy/foodieapp/internal/domain"

type Product struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	Price    float64       `json:"price"`
	Category string        `json:"category"`
	Image    *ProductImage `json:"image,omitempty"`
}

type ProductImage struct {
	Thumbnail string `json:"thumbnail"`
	Mobile    string `json:"mobile"`
	Tablet    string `json:"tablet"`
	Desktop   string `json:"desktop"`
}

func FromDomainProduct(p domain.Product) Product {
	var img *ProductImage
	if p.Image != nil {
		img = &ProductImage{
			Thumbnail: p.Image.Thumbnail,
			Mobile:    p.Image.Mobile,
			Tablet:    p.Image.Tablet,
			Desktop:   p.Image.Desktop,
		}
	}
	return Product{
		ID:       p.ID,
		Name:     p.Name,
		Price:    p.Price,
		Category: p.Category,
		Image:    img,
	}
}

func FromDomainProducts(products []domain.Product) []Product {
	if products == nil {
		return nil
	}
	res := make([]Product, len(products))
	for i, p := range products {
		res[i] = FromDomainProduct(p)
	}
	return res
}
