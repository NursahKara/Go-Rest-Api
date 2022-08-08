package cartController

import (
	"encoding/json"
	"final_project/database"
	entities "final_project/entities"
	"net/http"

	"github.com/gorilla/mux"
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

func GetCartItems(w http.ResponseWriter, r *http.Request) {
	var cartItems []entities.CartItem
	database.Instance.Joins("Product").Find(&cartItems)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cartItems)
}

func GetCartById(w http.ResponseWriter, r *http.Request) {

	cartId := mux.Vars(r)["id"]
	if checkIfCartExists(cartId) == false {
		json.NewEncoder(w).Encode("Cart Not Found!")
		return
	}
	var cart entities.Cart
	database.Instance.Preload("Items.Product").Joins("Customer").First(&cart, cartId)

	if len(cart.Items) > 0 {
		for _, item := range cart.Items {
			cart.TotalAmount += item.Product.Price * float64(100+item.Product.Vat) / 100 * item.Amount
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

func CreateCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cart entities.Cart
	json.NewDecoder(r.Body).Decode(&cart)
	database.Instance.Create(&cart)
	json.NewEncoder(w).Encode(cart)
}

func CreateCartItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cartItem entities.CartItem
	json.NewDecoder(r.Body).Decode(&cartItem)
	database.Instance.Create(&cartItem)
	json.NewEncoder(w).Encode(cartItem)
}

func DeleteCartItemById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cartItemId := mux.Vars(r)["id"]

	if checkIfCartItemExists(cartItemId) == false {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Cart Not Found!")
		return
	}
	database.Instance.Delete(&entities.CartItem{}, cartItemId)
	json.NewEncoder(w).Encode("Cart Deleted Successfully!")
}

func checkIfCartExists(cartId string) bool {
	var cart entities.Cart
	database.Instance.First(&cart, cartId)
	if cart.ID == 0 {
		return false
	}
	return true
}

func checkIfCartItemExists(cartItemId string) bool {
	var item entities.CartItem
	database.Instance.First(&item, cartItemId)
	if item.ID == 0 {
		return false
	}
	return true
}
