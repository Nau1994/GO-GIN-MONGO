package main

import (
	"GoGinMongo/routes"
	_ "GoGinMongo/services"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("main")
	// Initialize MongoDB connection
	// services.ConnectToMongo("mongodb://localhost:27017")

	// Initialize Gin router
	r := gin.Default()

	// Register routes
	routes.RegisterUserRoutes(r)
	routes.RegisterPostRoutes(r)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server failed to start")
	}
}

//docker run -d --rm --name mongo-container -p 27017:27017 -e MONGO_INITDB_DATABASE=app mongo
