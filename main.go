package main

import (
	cartController "final_project/controller/cartController"
	customerController "final_project/controller/customerController"
	orderController "final_project/controller/orderController"
	productController "final_project/controller/productController"
	"final_project/database"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {

	// Load Configurations from config.json using Viper
	//LoadAppConfig()

	// Initialize Database
	//database.Connect(AppConfig.ConnectionString)
	database.Connect("host=localhost port=5432 user=postgres password=aqweds123 dbname=PropertyFinder sslmode=disable")
	database.Migrate()

	// Initialize the router
	router := mux.NewRouter().StrictSlash(true)

	// Register Routes
	RegisterProductRoutes(router)

	// Start the server
	log.Println(fmt.Sprintf("Starting Server on port %s", "4000"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", "4000"), router))
}

func RegisterProductRoutes(router *mux.Router) {
	//PRODUCT
	router.HandleFunc("/api/products", productController.GetProducts).Methods("GET") // List Products
	router.HandleFunc("/api/products/{id}", productController.GetProductById).Methods("GET")
	router.HandleFunc("/api/products/{id}", productController.UpdateProduct).Methods("PUT")
	router.HandleFunc("/api/products/{id}", productController.DeleteProduct).Methods("DELETE")
	router.HandleFunc("/api/products", productController.CreateProduct).Methods("POST")

	//CUSTOMER
	router.HandleFunc("/api/customers", customerController.GetCustomers).Methods("GET")
	router.HandleFunc("/api/customers", customerController.CreateCustomer).Methods("POST")

	//ORDER
	router.HandleFunc("/api/orders/{cartId}", orderController.CreateOrder).Methods("POST") // Complete Order
	router.HandleFunc("/api/order/set-discount-threshold", orderController.SetDiscountThreshold).Methods("POST")
	router.HandleFunc("/api/orders", orderController.GetOrders).Methods("GET")

	//CART
	router.HandleFunc("/api/carts", cartController.GetCarts).Methods("GET")
	router.HandleFunc("/api/cartItems", cartController.GetCartItems).Methods("GET")
	router.HandleFunc("/api/carts/{id}", cartController.GetCartById).Methods("GET") // Show Cart
	router.HandleFunc("/api/carts", cartController.CreateCart).Methods("POST")      // Add To Cart
	router.HandleFunc("/api/carts/items", cartController.CreateCartItem).Methods("POST")
	router.HandleFunc("/api/cartItem/{id}", cartController.DeleteCartItemById).Methods("DELETE") // Delete Cart Item

}
