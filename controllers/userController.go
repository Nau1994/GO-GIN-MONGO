package controllers

import (
	"GoGinMongo/models"
	"GoGinMongo/services"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection

func init() {
	fmt.Println("userController init")
	userCollection = services.GetMongoClient().Database("app").Collection("users")
}

// CreateUser inserts a new user document into the "users" collection
func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	_, err := userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User inserted"})
}

// GetUser retrieves a user document based on the name provided in the URL parameters
func GetUser(c *gin.Context) {
	name := c.Param("name")
	filter := bson.D{{"name", name}}

	var user models.User
	err := userCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetAllUsers retrieves all user documents from the "users" collection
func GetAllUsers(c *gin.Context) {
	cursor, err := userCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}
	defer cursor.Close(context.TODO())

	var users []models.User
	if err = cursor.All(context.TODO(), &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// UpdateUser updates a user document based on the name provided in the URL parameters
func UpdateUser(c *gin.Context) {
	name := c.Param("name")
	filter := bson.D{{"name", name}}

	var updateData models.User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.D{{"$set", updateData}}

	_, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User updated"})
}

// DeleteUser deletes a user document based on the name provided in the URL parameters
func DeleteUser(c *gin.Context) {
	name := c.Param("name")
	filter := bson.D{{"name", name}}

	_, err := userCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
