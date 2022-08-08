package model

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Vat   int     `json:"vat"`
}

type Customer struct {
	gorm.Model
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Order struct {
	gorm.Model
	CustomerId  uint
	IsCompleted bool
	Total       float64
	Customer    Customer
	Items       []OrderItem
}

type OrderItem struct {
	gorm.Model
	OrderId   uint
	ProductId uint
	Product   Product
	Amount    float64 `json:"amount"`
}

type Cart struct {
	gorm.Model
	CustomerId  uint
	Customer    Customer
	Items       []CartItem
	TotalAmount float64 `gorm:"-:all"`
}

type CartItem struct {
	gorm.Model
	CartId    uint
	ProductId uint
	Product   Product
	Amount    float64 `json:"amount"`
}

type AppSettings struct {
	gorm.Model
	DiscountThreshold float64
}
