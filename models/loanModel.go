package models

type Loan struct {
	PK_Acc_No      string  `json:"accNo"`
	FK_Customer_Id string  `json:"customer_id" gorm:"primaryKey"`
	Bank_IFSC      string  `json:"bankIFSC"`
	Bank_Branch_Id string  `json:"bankBranchId"`
	AccOpenDate    string  `json:"accOpenDate"`
	ReturnDate     string  `json:"returnDate"`
	Principal      float64 `json:"principal"`
	Roi            float64 `json:"roi"`
}
