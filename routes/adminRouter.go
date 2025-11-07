package routes

import (
	"golang-banking-management-system/controllers"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {
	admin := r.Group("/admin")
	{
		admin.POST("/bank/create", controllers.CreateBank)
		admin.POST("/branch/create", controllers.CreateBankBranch)
		admin.DELETE("/bank/:BankId", controllers.DeleteBank)
		admin.DELETE("/branch/:BranchId", controllers.DeleteBankBranch)
		admin.GET("/branch/:BankId", controllers.DisplayBranches)
		admin.PUT("/update/:BranchId", controllers.UpdateBranch)
	}
}
