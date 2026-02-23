package service

import "github.com/Sanjaiy/foodieapp/internal/helpers"

type PromoService struct {
	lookup *helpers.CouponLookup
}

func NewPromoService(lookup *helpers.CouponLookup) *PromoService {
	return &PromoService{lookup: lookup}
}

func (s *PromoService) ValidateCoupon(code string) bool {
	return s.lookup.IsValid(code)
}
