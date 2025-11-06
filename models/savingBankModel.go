package models

type SavingBankAcc struct {
	PK_Customer_Id    string  `json:"pk_customer_id" gorm:"primaryKey"`
	Bank_IFSC         string  `json:"bankIFSC"`
	FK_Bank_Branch_Id string  `json:"bankBranchId"`
	AccOpenDate       string  `json:"accOpenDate"`
	Acc_No            string  `json:"accNo"`
	Acc_Balance       float64 `json:"acc_balance"`
}
