package customerController

import (
	"encoding/json"
	"final_project/database"
	entities "final_project/entities"
	"net/http"
)

func GetCustomers(w http.ResponseWriter, r *http.Request) {

	var customers []entities.Customer
	database.Instance.Find(&customers)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customers)
}

func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var customer entities.Customer
	json.NewDecoder(r.Body).Decode(&customer)
	database.Instance.Create(&customer)
	json.NewEncoder(w).Encode(customer)
}
