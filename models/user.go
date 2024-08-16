package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID    primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name  string               `bson:"name" json:"name" binding:"required"`
	Age   int                  `bson:"age" json:"age" binding:"required"`
	Posts []primitive.ObjectID `bson:"posts" json:"posts"`
}

var UserValidator = bson.D{
	{"$jsonSchema", bson.D{
		{"bsonType", "object"},
		{"required", bson.A{"name", "age", "posts"}},
		{"properties", bson.D{
			{"name", bson.D{
				{"bsonType", "string"},
				{"description", "must be a string and is required"},
			}},
			{"age", bson.D{
				{"bsonType", "int"},
				{"description", "must be an integer"},
			}},
			{"posts", bson.D{
				{"bsonType", "array"},
				{"items", bson.D{
					{"bsonType", "objectId"},
					{"description", "must be a valid ObjectId"},
				}},
				{"description", "must be an array of post references"},
			}},
		}},
	}},
}
