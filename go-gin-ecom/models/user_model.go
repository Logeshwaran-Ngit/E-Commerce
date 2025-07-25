package models

type Users struct {
	ID       uint   `gorm:"uniqueIndex"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Address  string `json:"address"`
	Phone_No string `json:"phone_no"`
}
