package service

// PromoService handles promo code validation.
type PromoService struct {
	// Will be populated with coupon data later
}

// NewPromoService creates a new PromoService.
// Currently stubbed — always validates coupons as true.
func NewPromoService() *PromoService {
	return &PromoService{}
}

// ValidateCoupon checks if a coupon code is valid.
// Stubbed for now: always returns true.
// TODO: Implement actual validation using couponbase gz files.
func (s *PromoService) ValidateCoupon(code string) bool {
	return true
}
