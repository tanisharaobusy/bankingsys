package routes

import (
	"golang-banking-management-system/controllers"

	"github.com/gin-gonic/gin"
)

func TransRoutes(r *gin.Engine) {
	trans := r.Group("/transaction")
	{
		trans.POST("/credit", controllers.Credit)
		trans.POST("/debit", controllers.Debit)
		trans.GET("/history/:CustomerId", controllers.History)
	}
}
