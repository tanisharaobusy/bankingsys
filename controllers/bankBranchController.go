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
	uniqBankId := BankBranchIdGenerator(bankBranch.Name, bankBranch.Branch, bankBranch.FK_Bank_Id)
	bankBranch.PK_Bank_Branch_Id = uniqBankId
	uniqBankIfsc := BankBranchIFSC(bankBranch.Name)
	bankBranch.Bank_IFSC = uniqBankIfsc
	if !BankExists(database.DB, bankBranch.FK_Bank_Id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bank Id does not exist"})
		return
	}

	log.Println("create bank branch called, bank id: ", bankBranch.FK_Bank_Id)
	//db handling through GORM

	err := validateStruct(bankBranch)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		if err := database.DB.Create(&bankBranch).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"bank_branch_id":   bankBranch.PK_Bank_Branch_Id,
		"bank_branch_ifsc": bankBranch.Bank_IFSC,
	})
	log.Println("status sent")

}

func DeleteBankBranch(c *gin.Context) {
	branchId := c.Param("BranchId")

	//db handling through GORM
	database.DB.Delete(&models.BankBranch{}, branchId)

	c.JSON(http.StatusOK, gin.H{"message": "Branch deleted"})
}

func DisplayBranches(c *gin.Context) {
	var branches []models.BankBranch
	bankId := c.Param("BankId")

	if err := database.DB.Where("fk_bank_id = ?", bankId).Find(&branches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(branches) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No branches found for this bank"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"branches": branches})
}

func UpdateBranch(c *gin.Context) {
	id := c.Param("BranchId")
	fmt.Println("Updating branch:", id)

	var branch models.BankBranch
	if err := database.DB.First(&branch, "pk_bank_branch_id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
		return
	}

	var updatedBranch models.BankBranch
	if err := c.ShouldBindJSON(&updatedBranch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if updatedBranch.Bank_IFSC != "" {
		branch.Bank_IFSC = updatedBranch.Bank_IFSC
	}
	if updatedBranch.Phone != 0 {
		branch.Phone = updatedBranch.Phone
	}

	if updatedBranch.Email != "" {
		branch.Email = updatedBranch.Email
	}
	if updatedBranch.Address != "" {
		branch.Address = updatedBranch.Address
	}

	if err := database.DB.Save(&branch).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update branch: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Branch updated successfully",
		"branch":  branch,
	})
}
