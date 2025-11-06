package routes

import (
	"golang-banking-management-system/controllers"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {
	admin := r.Group("/admin")
	{
		admin.POST("/bank/create", controllers.CreateBank)
		admin.POST("/branches/create", controllers.CreateBankBranch)
		admin.DELETE("/bank/:BankId", controllers.DeleteBank)
		admin.DELETE("/branches/:BranchId", controllers.DeleteBankBranch)
	}
}
