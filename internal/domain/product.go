package domain

type Product struct {
	ID       string
	Name     string
	Price    float64
	Category string
	Image    *ProductImage
}

type ProductImage struct {
	Thumbnail string
	Mobile    string
	Tablet    string
	Desktop   string
}
