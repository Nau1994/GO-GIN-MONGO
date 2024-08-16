package services

import (
	"GoGinMongo/models"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func init() {
	fmt.Println("mongoservice init")
	ConnectToMongo("mongodb://localhost:27017")
}

// ConnectToMongo establishes a connection to MongoDB
func ConnectToMongo(uri string) {
	var err error
	fmt.Println("ConnectToMongo")
	clientOptions := options.Client().ApplyURI(uri)
	// Connect to MongoDB
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// err = client.Connect(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")

	// Create collections and apply validators
	createCollectionWithValidators("users", models.UserValidator)
	createCollectionWithValidators("posts", models.PostValidator)
}

// createCollectionWithValidators creates a collection with schema validation
func createCollectionWithValidators(collectionName string, validator bson.D) {
	db := client.Database("app")
	cmd := bson.D{
		{"create", collectionName},
		{"validator", validator},
	}
	err := db.RunCommand(context.TODO(), cmd).Err()
	if err != nil && err.Error() != "(NamespaceExists) a collection 'app."+collectionName+"' already exists" {
		log.Fatal(err)
	}
}

// GetMongoClient returns the MongoDB client instance
func GetMongoClient() *mongo.Client {
	return client
}
