package models

type Product struct {
	ID          uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Description string  `json:"description"`
	Model       string  `json:"model"`
	Price       float64 `json:"price"`
	Stock       uint    `json:"stock"`
}
