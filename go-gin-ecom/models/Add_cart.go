package models

type Add_cart struct {
	Card_Id       uint    `gorm:"primaryKey" json:"p_id"`
	User_Id       uint    `json:"user_id"`
	Product_Id    uint    `json:"product_id"`
	Product_Stock uint    `json:"stock"`
	User          Users   `gorm:"foreignKey:User_Id" json:"user,omitempty"`
	Product       Product `gorm:"foreignKey:Product_Id" json:"product,omitempty"`
}
