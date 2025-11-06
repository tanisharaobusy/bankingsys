package routes

import (
	"golang-banking-management-system/controllers"

	"github.com/gin-gonic/gin"
)

func LoanRoutes(r *gin.Engine) {
	loan := r.Group("/loan")
	{
		loan.POST("/create", controllers.CreateLoan)
		loan.DELETE("/bank/:Cust_Id", controllers.DeleteLoan)
		loan.GET("/details/:loanAccNo", controllers.LoanDetails)

	}
}
