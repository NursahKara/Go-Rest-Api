package orderController

import (
	"encoding/json"
	"final_project/database"
	entities "final_project/entities"
	"math"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cartId := mux.Vars(r)["cartId"]

	var order entities.Order
	var cart entities.Cart
	var customerOrders []entities.Order
	//var cartItems []entities.CartItem

	database.Instance.Preload("Items.Product").Find(&cart, cartId)
	// id, err := strconv.Atoi(cartId)
	// if err != nil {
	if cart.ID == 0 {
		json.NewEncoder(w).Encode("Cart Not Found!")
		return
	}
	//}

	database.Instance.Preload("Items.Product").Where("customer_id = ?", cart.CustomerId).Find(&customerOrders)
	//database.Instance.Preload("Items.Product").Where("cart_id = ?", cartId).Find(&cartItems)

	var appSettings entities.AppSettings
	database.Instance.First(&appSettings)

	if appSettings.ID == 0 {
		json.NewEncoder(w).Encode("Set Discount Thresold First!")
		return
	}

	discountThreshold := appSettings.DiscountThreshold

	var fourthOrderDiscount float64
	var moreThanThreeProductsDiscount float64
	var exceededThresholdInTheLastMonthDiscount float64

	var totalAmount float64

	if len(cart.Items) > 0 {
		for _, item := range cart.Items {
			totalAmount += item.Product.Price * float64(100+item.Product.Vat) / 100 * item.Amount
		}
	}

	if len(customerOrders) > 0 {
		var totalAmountOfOrders float64
		for _, order := range customerOrders {

			diff := order.CreatedAt.Sub(time.Now())

			if diff.Hours() < float64(720) {
				var orderTotal float64
				for _, item := range order.Items {
					orderTotal += item.Product.Price * float64(100+item.Product.Vat) / 100 * item.Amount
				}
				totalAmountOfOrders += orderTotal
			}
		}

		if totalAmountOfOrders > discountThreshold {
			exceededThresholdInTheLastMonthDiscount = totalAmount * 10 / 100
		}
	}

	var orderItems []entities.OrderItem
	if len(cart.Items) > 0 {
		for _, item := range cart.Items {
			lineTotal := item.Product.Price * float64(100+item.Product.Vat) / 100 * item.Amount

			if rune(item.Amount) > 3 {
				moreThanThreeProductsDiscount += (item.Product.Price * float64(100+item.Product.Vat) / 100 * (item.Amount - 3)) * 8 / 100
			}

			if (len(customerOrders)+1)%4 == 0 && totalAmount > discountThreshold {
				if item.Product.Vat == 18 {
					fourthOrderDiscount += lineTotal * 15 / 100
				} else if item.Product.Vat == 8 {
					fourthOrderDiscount += lineTotal * 10 / 100
				}
			}

			orderItems = append(orderItems, entities.OrderItem{ProductId: item.ProductId, Product: item.Product, Amount: item.Amount})
		}
	}

	order.Items = orderItems
	order.IsCompleted = true
	order.CustomerId = cart.CustomerId
	order.Customer = cart.Customer
	order.Total = totalAmount - math.Max(exceededThresholdInTheLastMonthDiscount, math.Max(fourthOrderDiscount, moreThanThreeProductsDiscount))
	database.Instance.Create(&order)

	json.NewEncoder(w).Encode(order)

}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	var orders []entities.Order
	database.Instance.Preload("Items.Product").Joins("Customer").Find(&orders)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

func SetDiscountThreshold(w http.ResponseWriter, r *http.Request) {
	var discountThreshold float64
	json.NewDecoder(r.Body).Decode(&discountThreshold)

	var appSettings entities.AppSettings

	appSettings.DiscountThreshold = discountThreshold

	database.Instance.FirstOrCreate(&appSettings)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appSettings)
}
