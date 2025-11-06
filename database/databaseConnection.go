package database

import (
	"fmt"
	"golang-banking-management-system/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=localhost user=postgres password=busy123 dbname=banksysdb port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Database connected")

	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Loan{})
	DB.AutoMigrate(&models.Transaction{})
	DB.AutoMigrate(&models.SavingBankAcc{})
	DB.AutoMigrate(&models.Bank{})
	DB.AutoMigrate(&models.BankBranch{})

	fmt.Println("Database migrated")
}
