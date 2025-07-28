package models

import "gorm.io/gorm"

type ProductOrder struct {
	gorm.Model
	OrderID        string         `json:"orderid"` // Logical order ID like "ORD-12345"
	Name           string         `json:"name"`
	Email          string         `json:"email"`
	PhoneNo        string         `json:"phoneno"`
	Address        string         `json:"address"`
	UserID         uint           `json:"user_id"`
	CurrentAmount  float64        `json:"currentamount"`
	DeliveryCharge float64        `json:"deliverycharge"`
	TotalAmount    float64        `json:"totalamount"`
	Status         string         `json:"status"`
	OrderItems     []OrderProduct `gorm:"foreignKey:ProductOrderID" json:"order_items"` // Correct relation
}

type OrderProduct struct {
	ID             uint    `gorm:"primaryKey"`
	ProductOrderID uint    `json:"product_order_id"`
	ProductID      uint    `json:"product_id"`
	Quantity       uint    `json:"quantity"`
	UnitPrice      float64 `json:"unit_price"`
	TotalPrice     float64 `json:"total_price"`
}
