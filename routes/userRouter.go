package routes

import (
	"golang-banking-management-system/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	user := r.Group("/users")
	{
		user.POST("/", controllers.CreateUser)
		user.DELETE("/delete/:CustomerId", controllers.DeleteUser)
		user.GET("/:CustomerId", controllers.UserDetails)
		user.GET("/loan/:CustomerId", controllers.DisplayLoans)
	}
}
