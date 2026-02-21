package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Sanjaiy/foodieapp/internal/dto"
	"github.com/Sanjaiy/foodieapp/internal/service"
)

type OrderHandler struct {
	svc *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{
		svc: svc,
	}
}

func (h *OrderHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var req dto.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "validation", "invalid JSON body")
		return
	}

	order, err := h.svc.PlaceOrder(r.Context(), req.ToDomainItems(), req.CouponCode)
	if err != nil {
		if err.Error() == "at least one item is required" ||
			err.Error() == "productId is required for each item" ||
			err.Error() == "quantity must be greater than 0" ||
			err.Error() == "invalid product specified" {
			writeError(w, http.StatusUnprocessableEntity, "validation", err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "internal", "failed to place order")
		return
	}

	writeJSON(w, http.StatusOK, dto.FromDomainOrder(order))
}
