package models

type Transaction struct {
	PK_Id          string  `json:"id" gorm:"primaryKey"`
	Created_At     string  `json:"created_at"`
	FK_Customer_Id string  `json:"customer_id" gorm:"foreignKey"`
	Type           string  `json:"type"`
	Mode           string  `json:"mode"`
	Recipient      string  `json:"recipient"`
	Amount         float64 `json:"amount"`
}
