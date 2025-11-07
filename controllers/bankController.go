package controllers

import (
	"fmt"
	"golang-banking-management-system/database"
	"golang-banking-management-system/models"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	//"gorm.io/gorm"
)

func saveNumber(n int) {
	os.WriteFile("static.txt", []byte(fmt.Sprintf("%d", n)), 0644)
}

func loadNumber() int {
	data, err := os.ReadFile("bankId.txt")
	if err != nil {
		return 0
	}

	num, _ := strconv.Atoi(strings.TrimSpace(string(data)))
	return num
}

func BankIdGenerator(name string) string {
	var newBankSerial int
	var newBankNo string

	loadLastNo := loadNumber()

	if loadLastNo == 0 {
		newBankSerial = 1
	} else {
		newBankSerial = loadLastNo + 1
	}

	saveNumber(newBankSerial)

	prefix := name
	if len(name) >= 4 {
		prefix = name[:4]
	}

	newBankNo = fmt.Sprintf("%s%04d", strings.ToUpper(prefix), newBankSerial)

	return newBankNo
}

func CreateBank(c *gin.Context) {
	log.Println("create bank called")
	var bank models.Bank
	if err := c.ShouldBindJSON(&bank); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uniqBankId := BankIdGenerator(bank.Name)
	bank.PK_Bank_Id = uniqBankId
	log.Println("create bank called, bank id: ", bank.PK_Bank_Id)
	err := validateStruct(bank)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		if err := database.DB.Create(&bank).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, bank)
	log.Println("status sent")

}

func DeleteBank(c *gin.Context) {
	bankId := c.Param("BankId")

	//db handling through GORM
	database.DB.Delete(&models.Bank{}, bankId)

	c.JSON(http.StatusOK, gin.H{"message": "Bank deleted"})
}