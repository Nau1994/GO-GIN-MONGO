package routes

import (
	"GoGinMongo/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine) {
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/", controllers.CreateUser)
		userRoutes.GET("/:name", controllers.GetUser)
		userRoutes.GET("/", controllers.GetAllUsers)
		userRoutes.PUT("/:name", controllers.UpdateUser)
		userRoutes.DELETE("/:name", controllers.DeleteUser)
	}
}
