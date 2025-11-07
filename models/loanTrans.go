package models

type LoanTrans struct {
	FK_Trans_Id string `json:"tranId" gorm:"foreignKey"`
	PK_Loan_Id  string `json:"loanId" gorm:"primaryKey"`
}
