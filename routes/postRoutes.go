package routes

import (
	"GoGinMongo/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterPostRoutes(router *gin.Engine) {
	postRoutes := router.Group("/posts")
	{
		postRoutes.POST("/", controllers.CreatePost)
		postRoutes.GET("/:id", controllers.GetPost)
		postRoutes.GET("/", controllers.GetAllPosts)
		postRoutes.GET("/user/:userId", controllers.GetPostsByUserID)
		postRoutes.PUT("/:id", controllers.UpdatePost)
		postRoutes.DELETE("/:id", controllers.DeletePost)
	}
}
