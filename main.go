package main

import (
	"github.com/gin-gonic/gin"
	//"gorm.io/gorm"
	"golang-banking-management-system/database"
	"golang-banking-management-system/routes"
	"log"
)

func main() {
	// gin ka instance bn rha h
	r := gin.Default()

	// database initialise ho rha h
	log.Println("db initialised")
	database.InitDB()
	//routes import ho rhe h

	routes.UserRoutes(r)
	routes.AdminRoutes(r)
	routes.LoanRoutes(r)
	routes.TransRoutes(r)
	//server kis port pe run hoga
	r.Run(":8080")
}
