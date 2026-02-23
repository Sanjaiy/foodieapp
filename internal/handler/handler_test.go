package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/Sanjaiy/foodieapp/internal/domain"
	"github.com/Sanjaiy/foodieapp/internal/dto"
)

var (
	baseURL = getEnv("TEST_BASE_URL", "http://localhost:8080")
	apiKey  = getEnv("TEST_API_KEY", "apitest")
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func TestListProducts(t *testing.T) {
	resp, err := http.Get(baseURL + "/api/product")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var products []domain.Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if len(products) == 0 {
		t.Fatal("expected at least 1 product, got 0")
	}

	for _, p := range products {
		if p.ID == "" || p.Name == "" || p.Price <= 0 {
			t.Errorf("invalid product: %+v", p)
		}
	}
}

func TestGetProductFound(t *testing.T) {
	resp, err := http.Get(baseURL + "/api/product/1")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var product domain.Product
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if product.ID != "1" {
		t.Errorf("expected id '1', got %q", product.ID)
	}
	if product.Name != "Waffle with Berries" {
		t.Errorf("expected 'Waffle with Berries', got %q", product.Name)
	}
}

func TestGetProductNotFound(t *testing.T) {
	resp, err := http.Get(baseURL + "/api/product/9999")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func postOrder(t *testing.T, body dto.OrderRequest) *http.Response {
	t.Helper()
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/api/order", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api_key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	return resp
}

func TestPlaceOrderWithoutCoupon(t *testing.T) {
	resp := postOrder(t, dto.OrderRequest{
		Items: []domain.OrderItem{
			{ProductID: "1", Quantity: 2},
			{ProductID: "3", Quantity: 1},
		},
	})
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var order dto.OrderResponse
	json.NewDecoder(resp.Body).Decode(&order)

	if order.ID == "" {
		t.Error("expected order ID, got empty")
	}
	// 6.50*2 + 8.00*1 = 21.00
	if order.Total != 21.0 {
		t.Errorf("expected total 21.00, got %.2f", order.Total)
	}
	if order.Discounts != 0 {
		t.Errorf("expected 0 discounts, got %.2f", order.Discounts)
	}
}

func TestPlaceOrderWithValidCoupon(t *testing.T) {
	resp := postOrder(t, dto.OrderRequest{
		Items: []domain.OrderItem{
			{ProductID: "1", Quantity: 2},
		},
		CouponCode: "OVER9000",
	})
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var order dto.OrderResponse
	json.NewDecoder(resp.Body).Decode(&order)

	// 6.50*2 = 13.00, 10% discount = 1.30, total = 11.70
	if order.Discounts != 1.3 {
		t.Errorf("expected discount 1.30, got %.2f", order.Discounts)
	}
	if order.Total != 11.7 {
		t.Errorf("expected total 11.70, got %.2f", order.Total)
	}
}

func TestPlaceOrderWithInvalidCoupon(t *testing.T) {
	resp := postOrder(t, dto.OrderRequest{
		Items: []domain.OrderItem{
			{ProductID: "1", Quantity: 2},
		},
		CouponCode: "FAKECODE1",
	})
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var order dto.OrderResponse
	json.NewDecoder(resp.Body).Decode(&order)

	// Invalid coupon → no discount
	if order.Total != 13.0 {
		t.Errorf("expected total 13.00, got %.2f", order.Total)
	}
	if order.Discounts != 0 {
		t.Errorf("expected 0 discounts, got %.2f", order.Discounts)
	}
}
