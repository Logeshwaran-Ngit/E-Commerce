package models

import "gorm.io/gorm"

type ProductOrder struct {
	gorm.Model
	OrderID        string         `gorm:"uniqueIndex" json:"orderid"`
	Name           string         `json:"name"`
	Email          string         `json:"email"`
	PhoneNo        string         `json:"phoneno"`
	Address        string         `json:"address"`
	UserID         uint           `json:"user_id"`
	CurrentAmount  float64        `json:"currentamount"`
	DeliveryCharge float64        `json:"deliverycharge"`
	TotalAmount    float64        `json:"totalamount"`
	Status         string         `json:"status"`
	OrderItems     []OrderProduct `gorm:"foreignKey:OrderID;references:OrderID" json:"order_items"`
}

type OrderProduct struct {
	ID         uint    `gorm:"primaryKey"`
	OrderID    string  `json:"order_id"`
	ProductID  uint    `json:"product_id"`
	Quantity   uint    `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
}
