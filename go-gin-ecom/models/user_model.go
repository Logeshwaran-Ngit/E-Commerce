package models

type Users struct {
	ID       uint           `gorm:"primaryKey" json:"user_Id"`
	Name     string         `json:"name"`
	Email    string         `gorm:"uniqueIndex" json:"email"`
	Password string         `json:"password"`
	Role     string         `json:"role"`
	Address  string         `json:"address"`
	Phone_No string         `json:"phone_no"`
	AddCart  []Add_cart     `gorm:"foreignKey:User_Id"`
	Orders   []ProductOrder `gorm:"foreignKey:UserID"`
}
