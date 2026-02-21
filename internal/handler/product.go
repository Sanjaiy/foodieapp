package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Sanjaiy/foodieapp/internal/dto"
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

	writeJSON(w, http.StatusOK, dto.FromDomainProducts(products))
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

	writeJSON(w, http.StatusOK, dto.FromDomainProduct(*product))
}

// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("ERROR: encoding response: %v", err)
	}
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(dto.ErrorResponse{
		Code:    code,
		Message: message,
	})
}
