package models

type Product struct {
	Product_Id          uint    `gorm:"primaryKey" json:"p_id"`
	Product_Name        string  `json:"name"`
	Product_Email       string  `json:"email"`
	Product_Discrpition string  `json:"discrpition"`
	Product_Model       string  `json:"model"`
	Product_Prize       float32 `json:"prize"`
	Product_Stock       uint    `json:"stock"`
}
