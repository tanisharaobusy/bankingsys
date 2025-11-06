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

	"github.com/gin-gonic/gin"
	//"gorm.io/gorm"
)

func saveAccSerialNo(n int) {
	os.WriteFile("AccSerialNo.txt", []byte(fmt.Sprintf("%d", n)), 0644)
}

func loadAccSerialNo() int {
	data, err := os.ReadFile("AccSerialNo.txt")
	if err != nil {
		return 0
	}

	num, _ := strconv.Atoi(strings.TrimSpace(string(data)))
	return num
}

func saveCustomerId(n int) {
	os.WriteFile("customerId.txt", []byte(fmt.Sprintf("%d", n)), 0644)
}

func loadCustomerId() int {
	data, err := os.ReadFile("customerId.txt")
	if err != nil {
		return 0
	}

	num, _ := strconv.Atoi(strings.TrimSpace(string(data)))
	return num
}

func AccNoGenerator(bankIfsc, branchId string) string {
	//code to generate a unique account no

	//bankid(first 3 digits) + branchid(last 3 digits) + 6 digit actual serial no, which can't be duplicate or repeated ever

	//safe side: check db for same accno - if not found, then only assign a global variable with the acc no so that it's value can be stored in the database - user table mainly

	var newAccSerialNo int
	var newAccNo string

	loadLastNo := loadAccSerialNo()

	if loadLastNo == 0 {
		newAccSerialNo = 1
	} else {
		newAccSerialNo = loadLastNo + 1
	}

	saveAccSerialNo(newAccSerialNo)

	prefix1 := bankIfsc
	if len(bankIfsc) >= 4 {
		prefix1 = bankIfsc[:4]
	}

	prefix2 := branchId
	if len(branchId) >= 2 {
		prefix2 = branchId[len(branchId)-4:]
	}

	newAccNo = fmt.Sprintf("%s%s%06d", strings.ToUpper(prefix1), strings.ToUpper(prefix2), newAccSerialNo)

	return newAccNo
}

func CustomerIdGenerator(bankBranch string) string {
	//code to generate a unique customerid
	//bankId unique for all/ if deleted then also can't be repeated

	//safe side: check db for same customerId - if not found, then only assign a global variable with the acc no so that it's value can be stored in the database - user table mainly

	var newCustSerialNo int
	var newCustomerIdNo string

	loadLastNo := loadCustomerId()

	if loadLastNo == 0 {
		newCustSerialNo = 1
	} else {
		newCustSerialNo = loadLastNo + 1
	}

	saveCustomerId(newCustSerialNo)

	prefix1 := bankBranch
	if len(bankBranch) >= 4 {
		prefix1 = bankBranch[:4]
	}

	newCustomerIdNo = fmt.Sprintf("%s%06d", strings.ToUpper(prefix1), newCustSerialNo)

	return newCustomerIdNo

}

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uniqAccNo := AccNoGenerator(user.Bank_IFSC, user.FK_Bank_Branch_Id)
	user.Acc_No = uniqAccNo
	uniqCustId := CustomerIdGenerator(user.FK_Bank_Branch_Id)
	user.PK_Customer_Id = uniqCustId
	//db handling through GORM
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	savingAcc, err := CreateSavingsAcc(user.Bank_IFSC, user.FK_Bank_Branch_Id, user.AccOpenDate, user.Acc_No, user.PK_Customer_Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, gin.H{
		"user":           user,
		"savingsAccount": savingAcc,
	})

}

func DeleteUser(c *gin.Context) {
	id := c.Param("CustomerId")

	//db handling through GORM
	database.DB.Delete(&models.User{}, id)

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func UserDetails(c *gin.Context) {
	var user models.User
	CustId := c.Param("CustomerId")

	if err := database.DB.Raw(`
	SELECT *
	FROM users 
	WHERE pk_customer_id = ?`, CustId).Scan(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"User Details ": user})
}

func DisplayLoans(c *gin.Context) {
	var loan []models.Loan
	CustomerId := c.Param("CustomerId")

	if err := database.DB.Raw("SELECT * FROM loans WHERE fk_customer_id = ?", CustomerId).Scan(&loan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(loan) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No loans found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"loans details": loan})
}
