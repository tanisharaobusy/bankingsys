package models

// import "gorm.io/gorm"
// ---------- MODEL ----------
type User struct {
	PK_Customer_Id    string `json:"customer_id" gorm:"primaryKey"`
	Name              string `json:"name"`
	Email             string `json:"email" `
	Phone             uint   `json:"phoneNo"`
	Address           string `json:"address"`
	Occupation        string `json:"occupation"`
	Bank_IFSC         string `json:"bankIFSC"`
	FK_Bank_Branch_Id string `json:"bankBranchId"`
	PAN               string `json:"pan"`
	AccOpenDate       string `json:"accOpenDate"`
	Acc_No            string `json:"accNo"`
}
