package routes

import (
	"golang-banking-management-system/controllers"

	"github.com/gin-gonic/gin"
)

func LoanRoutes(r *gin.Engine) {
	loan := r.Group("/loans")
	{
		loan.POST("/create", controllers.CreateLoan)
		loan.DELETE("/close/:LoanAccNo", controllers.DeleteLoan)
		loan.GET("/details/:LoanAccNo", controllers.LoanDetails)
		loan.GET("/transH/:LoanAccNo", controllers.LoanHistory)
	}
}
