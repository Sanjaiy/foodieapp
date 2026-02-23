package handler

import (
	"log"
	"net/http"

	"github.com/Sanjaiy/foodieapp/internal/service"
)

type ProductHandler struct {
	svc *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.svc.ListProducts(r.Context())
	if err != nil {
		log.Printf("ERROR: listing products: %v", err)
		writeError(w, http.StatusInternalServerError, "internal", "failed to list products")
		return
	}

	writeJSON(w, http.StatusOK, products)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("productId")
	if productID == "" {
		writeError(w, http.StatusBadRequest, "validation", "product ID is required")
		return
	}

	product, err := h.svc.GetProduct(r.Context(), productID)
	if err != nil {
		log.Printf("ERROR: getting product %s: %v", productID, err)
		writeError(w, http.StatusInternalServerError, "internal", "failed to get product")
		return
	}

	if product == nil {
		writeError(w, http.StatusNotFound, "not_found", "Product not found")
		return
	}

	writeJSON(w, http.StatusOK, *product)
}
