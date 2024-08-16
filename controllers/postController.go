package controllers

import (
	"GoGinMongo/models"
	"GoGinMongo/services"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var postsCollection *mongo.Collection
var usersCollection *mongo.Collection

func init() {
	fmt.Println("postController init")
	postsCollection = services.GetMongoClient().Database("app").Collection("posts")
	usersCollection = services.GetMongoClient().Database("app").Collection("users")
}

// CreatePost creates a new post and updates the user's posts array
func CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.BindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate UserID
	if _, err := primitive.ObjectIDFromHex(post.UserID.Hex()); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UserID"})
		return
	}

	// Insert the post into the posts collection
	result, err := postsCollection.InsertOne(context.TODO(), post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the inserted post ID
	postId, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve post ID"})
		return
	}

	// Update the user's posts array
	filter := bson.M{"_id": post.UserID}
	update := bson.M{"$push": bson.M{"posts": postId}}

	_, err = usersCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post created", "postId": postId})
}

// GetPost retrieves a post document by post ID
func GetPost(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	filter := bson.D{{"_id", id}}

	var post models.Post
	err := postsCollection.FindOne(context.TODO(), filter).Decode(&post)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	c.JSON(http.StatusOK, post)
}

// GetAllPosts retrieves all post documents from the "posts" collection
func GetAllPosts(c *gin.Context) {
	cursor, err := postsCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve posts"})
		return
	}
	defer cursor.Close(context.TODO())

	var posts []models.Post
	if err = cursor.All(context.TODO(), &posts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process posts"})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// GetPostsByUserID retrieves all posts created by a specific user
func GetPostsByUserID(c *gin.Context) {

	userId, _ := primitive.ObjectIDFromHex(c.Param("userId"))

	filter := bson.D{{"userId", userId}}

	cursor, err := postsCollection.Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve posts"})
		return
	}
	defer cursor.Close(context.TODO())

	var posts []models.Post
	if err = cursor.All(context.TODO(), &posts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process posts"})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// UpdatePost updates a post document based on the post ID
func UpdatePost(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	filter := bson.D{{"_id", id}}

	var updateData bson.M
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	update := bson.D{{"$set", updateData}}

	_, err := postsCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post updated"})
}

// DeletePost deletes a post and removes the postId from the user's posts array
func DeletePost(c *gin.Context) {

	postId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Find the post to get the UserID
	var post models.Post
	err = postsCollection.FindOne(context.TODO(), bson.M{"_id": postId}).Decode(&post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete the post
	_, err = postsCollection.DeleteOne(context.TODO(), bson.M{"_id": postId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update the user's posts array by pulling the postId
	filter := bson.M{"_id": post.UserID}
	update := bson.M{"$pull": bson.M{"posts": postId}}

	_, err = usersCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}
