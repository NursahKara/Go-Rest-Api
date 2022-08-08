package test

import (
	"encoding/json"
	"final_project/database"
	entities "final_project/entities"
	"net/http"
	"net/http/httptest"
	"testing"
)

func GetCarts(w http.ResponseWriter, r *http.Request) {

	var carts []entities.Cart
	database.Instance.Preload("Items.Product").Joins("Customer").Find(&carts)

	if len(carts) > 0 {
		for i, cart := range carts {
			if len(cart.Items) > 0 {
				for _, item := range cart.Items {
					cart.TotalAmount += item.Product.Price * float64(100+item.Product.Vat) / 100 * item.Amount
				}
				carts[i] = cart
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(carts)
}

func TestGetCarts(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/carts", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(GetCarts)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `[{"id":1,"created_at":"2022-08-07 19:27:53.910073+03","updated_at":"2022-08-07 19:27:53.910073+03","deleted_at":"","customer_id":"1"},
				  {"id":2,"created_at":"2022-08-07 19:41:27.408444+03","updated_at":"2022-08-07 19:41:27.408444+03","deleted_at":"","customer_id":"1"}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
