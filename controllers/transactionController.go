package controllers

import (
	"fmt"
	"log"

	//"golang-banking-management-system/controllers"

	"golang-banking-management-system/database"
	"golang-banking-management-system/models"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	//"gorm.io/gorm"
)

func getTime() time.Time {

	now := time.Now()
	return now

}

func loadTransNo() int {
	data, err := os.ReadFile("transSerial.txt")
	if err != nil {
		return 0
	}

	num, _ := strconv.Atoi(strings.TrimSpace(string(data)))
	return num
}

func saveTransNo(n int) {
	os.WriteFile("transSerial.txt", []byte(fmt.Sprintf("%d", n)), 0644)
}

func TransNoGenerator(custId, typeTrans string) string {
	//code to generate a unique account no

	//bankid(first 3 digits) + branchid(last 3 digits) + 6 digit actual serial no, which can't be duplicate or repeated ever

	//safe side: check db for same accno - if not found, then only assign a global variable with the acc no so that it's value can be stored in the database - user table mainly

	var newTransSerialNo int
	var newTransNo string

	loadLastNo := loadTransNo()

	if loadLastNo == 0 {
		newTransSerialNo = 1
	} else {
		newTransSerialNo = loadLastNo + 1
	}

	saveTransNo(newTransSerialNo)

	prefix1 := custId
	if len(custId) >= 4 {
		prefix1 = custId[:4]
	}

	prefix2 := typeTrans
	if len(typeTrans) >= 1 {
		prefix2 = typeTrans[:1]
	}

	newTransNo = fmt.Sprintf("%s%s%06d", strings.ToUpper(prefix1), strings.ToUpper(prefix2), newTransSerialNo)

	return newTransNo
}

func Credit(c *gin.Context) {
	var trans models.Transaction

	transaction, exists := c.Get("transaction")
	if !exists {
		log.Println("error in 82")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction not found"})
		return
	}

	if t, ok := transaction.(models.Transaction); ok {
		fmt.Println("Transaction:", t)
		trans = t
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid transaction type"})
		return
	}

	var user models.User
	var loan models.Loan
	loanNo, exists := c.Get("Loan")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Loan Account not found"})
		return
	}

	if l, ok := loanNo.(models.Loan); ok {
		fmt.Println("Loan:", l)
		loan.PK_Acc_No = l.PK_Acc_No
	} else {
		fmt.Println("line 107")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid transaction type"})
		return
	}

	// if err := c.ShouldBindJSON(&trans); err != nil {
	// 	log.Println("error in 112")
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	if err := database.DB.Where("pk_customer_id = ?", trans.FK_Customer_Id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	trans.Type = "Credit"
	trans.PK_Id = TransNoGenerator(trans.FK_Customer_Id, trans.Type)
	trans.Created_At = getTime().String()

	if err := database.DB.Create(&trans).Error; err != nil {
		log.Println("error in 127")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var acc models.SavingBankAcc
	if err := database.DB.Where("pk_customer_id = ?", trans.FK_Customer_Id).First(&acc).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Saving account not found"})
		return
	}

	acc.Acc_Balance += trans.Amount
	success := 0
	err := database.DB.Save(&acc).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update balance"})
		return
	} else {
		if loan.PK_Acc_No != "" && trans.Tran == "Loan" {
			success = SaveLoanTrans(loan.PK_Acc_No, trans.PK_Id)
			c.Set("tranId", trans.PK_Id)
		}

	}
	if success == 0 && trans.Tran == "Loan" {
		c.JSON(http.StatusOK, gin.H{
			"transaction":  trans,
			"new_balance":  acc.Acc_Balance,
			"loan_account": loan.PK_Acc_No,
		})
	} else if success == 0 && trans.Tran == "Normal" {
		c.JSON(http.StatusOK, gin.H{
			"transaction": trans,
			"new_balance": acc.Acc_Balance,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Couldn't save transaction",
		})
	}

}

func History(c *gin.Context) {
	var trans []models.Transaction
	CustId := c.Param("CustomerId")

	if err := database.DB.Raw("SELECT * FROM transactions WHERE fk_customer_id = ?", CustId).Scan(&trans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(trans) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No transactions found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": trans})
}

func Debit(c *gin.Context) {
	var trans models.Transaction
	var user models.User

	loanAccNo := c.Param("LoanAccNo")
	if err := c.ShouldBindJSON(&trans); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Where("pk_customer_id = ?", trans.FK_Customer_Id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	trans.Type = "Debit"
	trans.PK_Id = TransNoGenerator(trans.FK_Customer_Id, trans.Type)
	trans.Created_At = getTime().String()

	var acc models.SavingBankAcc
	if err := database.DB.Where("pk_customer_id = ?", trans.FK_Customer_Id).First(&acc).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Saving account not found"})
		return
	}

	if acc.Acc_Balance < trans.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	var loan models.Loan
	if err := database.DB.Where("pk_acc_no= ?", loanAccNo).First(&loan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan account not found"})
		return
	}
	acc.Acc_Balance -= trans.Amount
	success := 0
	loan.Principal -= trans.Amount
	err1 := database.DB.Create(&trans).Error
	err2 := database.DB.Save(&loan).Error
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update loan pricipal"})
		return
	}
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	} else {
		if loanAccNo != "" && trans.Tran == "Loan" {
			success = SaveLoanTrans(loanAccNo, trans.PK_Id)
			log.Println(success)
		}

	}

	if err := database.DB.Save(&acc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update balance"})
		return
	}

	if trans.Tran == "Loan" {
		c.JSON(http.StatusOK, gin.H{
			"transaction":  trans,
			"new_balance":  acc.Acc_Balance,
			"loan_account": loanAccNo,
		})
	} else if trans.Tran == "Normal" {
		c.JSON(http.StatusOK, gin.H{
			"transaction": trans,
			"new_balance": acc.Acc_Balance,
		})
	} else if success == -1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Couldn't save transaction",
		})
	}
}
