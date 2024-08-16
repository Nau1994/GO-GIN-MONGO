package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title   string             `bson:"title" json:"title" binding:"required"`
	Message string             `bson:"message" json:"message" binding:"required"`
	UserID  primitive.ObjectID `bson:"userId" json:"userId" binding:"required"`
}

var PostValidator = bson.D{
	{"$jsonSchema", bson.D{
		{"bsonType", "object"},
		{"required", bson.A{"title", "message", "userId"}},
		{"properties", bson.D{
			{"title", bson.D{
				{"bsonType", "string"},
				{"description", "must be a string and is required"},
			}},
			{"message", bson.D{
				{"bsonType", "string"},
				{"description", "must be a string and is required"},
			}},
			{"userId", bson.D{
				{"bsonType", "objectId"},
				{"description", "must be a valid ObjectId and is required"},
			}},
		}},
	}},
}
