package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// ConnectToMongo establishes a connection to MongoDB
func ConnectToMongo() {
	var err error
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")
}

// CreateDocument inserts a new user document into the "users" collection
func CreateDocument(c *gin.Context) {
	collection := client.Database("app").Collection("users")

	var document bson.M
	if err := c.BindJSON(&document); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := collection.InsertOne(context.TODO(), document)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Document inserted"})
}

// InsertManyDocuments inserts multiple user documents into the "users" collection
// InsertManyDocuments inserts multiple documents into the "users" collection
func InsertManyDocuments(c *gin.Context) {
	collection := client.Database("app").Collection("users")

	var documents []bson.M
	if err := c.BindJSON(&documents); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if documents slice is empty
	if len(documents) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No documents provided"})
		return
	}

	// Attempt to insert the documents
	result, err := collection.InsertMany(context.TODO(), toInterfaceSlice(documents))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Documents inserted", "result": result})
}

// Helper function to convert []bson.M to []interface{}
func toInterfaceSlice(docs []bson.M) []interface{} {
	var result []interface{}
	for _, doc := range docs {
		result = append(result, doc)
	}
	return result
}

// GetDocument retrieves a user document based on the name provided in the URL parameters
func GetDocument(c *gin.Context) {
	collection := client.Database("app").Collection("users")

	name := c.Param("name")
	filter := bson.D{{"name", name}}

	var result bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetAllDocuments retrieves all user documents from the "users" collection
func GetAllDocuments(c *gin.Context) {
	collection := client.Database("app").Collection("users")

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.TODO())

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

// UpdateDocument updates a user document based on the name provided in the URL parameters
func UpdateDocument(c *gin.Context) {
	collection := client.Database("app").Collection("users")

	name := c.Param("name")
	filter := bson.D{{"name", name}}

	var updateData bson.M
	if err := c.BindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.D{{"$set", updateData}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Document updated"})
}

// DeleteDocument deletes a user document based on the name provided in the URL parameters
func DeleteDocument(c *gin.Context) {
	collection := client.Database("app").Collection("users")

	name := c.Param("name")
	filter := bson.D{{"name", name}}

	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Document deleted"})
}

// AggregateDocuments performs an aggregation on the "users" collection
func AggregateDocuments(c *gin.Context) {
	collection := client.Database("app").Collection("users")

	// Extract age from the request body
	var requestBody struct {
		Age int `json:"age"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Form the match stage based on the provided age
	matchStage := bson.D{}
	if requestBody.Age > 0 {
		matchStage = bson.D{{"age", bson.D{{"$gte", requestBody.Age}}}}
	}

	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{
		// Stage 1: Match documents to filter by age
		{{"$match", matchStage}},

		// Stage 2: Group documents by age and calculate the count
		{{"$group", bson.D{
			{"_id", "$age"},
			{"count", bson.D{{"$sum", 1}}},
		}}},

		// Stage 3: Sort groups by count in descending order
		{{"$sort", bson.D{{"count", -1}}}},

		// Stage 4: Project the fields to include in the output
		{{"$project", bson.D{
			{"_id", 0}, // Exclude the _id field from the output
			{"age", "$_id"},
			{"count", 1},
		}}},

		// Stage 5: Limit the number of documents in the result
		{{"$limit", 10}},
	}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

func main() {
	ConnectToMongo()
	r := gin.Default()

	r.POST("/create", CreateDocument)
	r.GET("/read/:name", GetDocument)
	r.PUT("/update/:name", UpdateDocument)
	r.DELETE("/delete/:name", DeleteDocument)
	r.POST("/aggregate", AggregateDocuments)

	// New endpoints
	r.GET("/documents", GetAllDocuments)
	r.POST("/documents", InsertManyDocuments)

	r.Run(":8080")
}
