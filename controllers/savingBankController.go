package controllers

import (
	"golang-banking-management-system/database"
	"golang-banking-management-system/models"
	"log"
)

func CreateSavingsAcc(bankIFSC, bankBranchId, accOpenDate, accNo, custId string) (*models.SavingBankAcc, error) {
	var account models.SavingBankAcc

	account.Bank_IFSC = bankIFSC
	account.FK_Bank_Branch_Id = bankBranchId
	account.AccOpenDate = accOpenDate

	account.Acc_No = accNo
	account.PK_Customer_Id = custId

	account.Acc_Balance = 0

	// Save to DB

	err := validateStruct(account)
	if err != nil {
		return nil, err

	} else {
		if err := database.DB.Create(&account).Error; err != nil {
			return nil, err
		}
	}

	log.Println("Savings account created successfully!")
	return &account, nil
}
