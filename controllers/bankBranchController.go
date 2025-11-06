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

func saveBankBranchNumber(n int) {
	os.WriteFile("bankBranchId.txt", []byte(fmt.Sprintf("%d", n)), 0644)
}

func loadBankBranchNumber() int {
	data, err := os.ReadFile("bankBranchId.txt")
	if err != nil {
		return 0
	}

	num, _ := strconv.Atoi(strings.TrimSpace(string(data)))
	return num
}

func saveBranchIFSCNumber(n int) {
	os.WriteFile("branchIFSC.txt", []byte(fmt.Sprintf("%d", n)), 0644)
}

func loadBranchIFSCNumber() int {
	data, err := os.ReadFile("branchIFSC.txt")
	if err != nil {
		return 0
	}

	num, _ := strconv.Atoi(strings.TrimSpace(string(data)))
	return num
}

func BankBranchIdGenerator(name, branch, bankId string) string {
	var newBankBranchSerial int
	var newBankBranchNo string

	loadLastNo := loadBankBranchNumber()

	if loadLastNo == 0 {
		newBankBranchSerial = 1
	} else {
		newBankBranchSerial = loadLastNo + 1
	}

	saveBranchIFSCNumber(newBankBranchSerial)

	prefix1 := name
	if len(name) >= 2 {
		prefix1 = name[:2]
	}

	prefix2 := branch
	if len(branch) >= 2 {
		prefix2 = branch[:2]
	}

	newBankBranchNo = fmt.Sprintf("%s%s%s%04d", strings.ToUpper(prefix1), strings.ToUpper(prefix2), bankId[6:], newBankBranchSerial)

	return newBankBranchNo
}

func BankBranchIFSC(name string) string {
	var newBranchIFSCSerial int
	var newBranchIFSCNo string

	loadLastNo := loadBranchIFSCNumber()

	if loadLastNo == 0 {
		newBranchIFSCSerial = 1
	} else {
		newBranchIFSCSerial = loadLastNo + 1
	}

	saveBankBranchNumber(newBranchIFSCSerial)

	prefix1 := name
	if len(name) >= 4 {
		prefix1 = name[:4]
	}

	newBranchIFSCNo = fmt.Sprintf("%s%s%04d", strings.ToUpper(prefix1), "000", newBranchIFSCSerial)

	return newBranchIFSCNo
}

func CreateBankBranch(c *gin.Context) {
	log.Println("create bank branch called")
	var bankBranch models.BankBranch
	if err := c.ShouldBindJSON(&bankBranch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uniqBankId := BankBranchIdGenerator(bankBranch.Name, bankBranch.Branch, bankBranch.Bank_Id)
	bankBranch.Bank_Branch_Id = uniqBankId
	uniqBankIfsc := BankBranchIFSC(bankBranch.Name)
	bankBranch.Bank_IFSC = uniqBankIfsc
	log.Println("create bank branch called, bank id: ", bankBranch.Bank_Id)
	//db handling through GORM
	if err := database.DB.Create(&bankBranch).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bankBranch)
	log.Println("status sent")

}

func DeleteBankBranch(c *gin.Context) {
	branchId := c.Param("BranchId")

	//db handling through GORM
	database.DB.Delete(&models.BankBranch{}, branchId)

	c.JSON(http.StatusOK, gin.H{"message": "Branch deleted"})
}

func DisplayBranches(c *gin.Context) {
	var branch []models.BankBranch
	bankId := c.Param("BankId")

	if err := database.DB.Raw("SELECT * FROM bank_branches WHERE bank_id = ?", bankId).Scan(&branch).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(branch) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No branches found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Branches": branch})
}
