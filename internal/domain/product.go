package domain

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
