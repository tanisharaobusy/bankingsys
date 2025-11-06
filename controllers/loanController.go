package controllers

import (
	"fmt"
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

func getDate() string {

	now := time.Now()
	date := now.Format("2006-01-02")
	return date

}

func loadLoanSerialNo() int {
	data, err := os.ReadFile("LoanSerialNo.txt")
	if err != nil {
		return 0
	}

	num, _ := strconv.Atoi(strings.TrimSpace(string(data)))
	return num
}

func saveLoanNo(n int) {
	os.WriteFile("LoanSerialNotxt", []byte(fmt.Sprintf("%d", n)), 0644)
}

func LoanNoGenerator(custId string) string {
	//code to generate a unique account no

	//bankid(first 3 digits) + branchid(last 3 digits) + 6 digit actual serial no, which can't be duplicate or repeated ever

	//safe side: check db for same accno - if not found, then only assign a global variable with the acc no so that it's value can be stored in the database - user table mainly

	var newLoanSerialNo int
	var newLoanNo string

	loadLastNo := loadLoanSerialNo()

	if loadLastNo == 0 {
		newLoanSerialNo = 1
	} else {
		newLoanSerialNo = loadLastNo + 1
	}

	saveLoanNo(newLoanSerialNo)

	prefix1 := custId
	if len(custId) >= 4 {
		prefix1 = custId[:4]
	}

	newLoanNo = fmt.Sprintf("%s%s%05d", strings.ToUpper(prefix1), "L", newLoanSerialNo)

	return newLoanNo
}

func CreateLoan(c *gin.Context) {
	var loan models.Loan
	var user models.User
	if err := c.ShouldBindJSON(&loan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Where("pk_customer_id = ?", loan.FK_Customer_Id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	uniqLoanNo := LoanNoGenerator(loan.FK_Customer_Id)
	loan.PK_Acc_No = uniqLoanNo
	loan.AccOpenDate = getDate()
	loan.Roi = 12

	//db handling through GORM
	if err := database.DB.Create(&loan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Loan account": loan})

}

func DeleteLoan(c *gin.Context) {
	LoanAccNo := c.Param("loanAccNo")
	var loan models.Loan

	// Find the loan
	if err := database.DB.Where("pk_acc_no = ?", LoanAccNo).First(&loan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
		return
	}

	// Check if principal is zero
	if loan.Principal != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Loan Principal not 0, can't close the loan."})
		return
	}

	// Delete the loan
	if err := database.DB.Where("pk_acc_no = ?", LoanAccNo).Delete(&models.Loan{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete loan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Loan deleted"})
}

func LoanDetails(c *gin.Context) {
	var loan models.Loan
	LoanAccNo := c.Param("loanAccNo")

	if err := database.DB.Raw("SELECT * FROM loans WHERE pk_acc_no = ?", LoanAccNo).Scan(&loan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Loan detail for account: " + LoanAccNo: loan,
	})

}
