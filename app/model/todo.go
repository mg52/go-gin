package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type ToDo struct {
	Id     primitive.ObjectID `bson:"_id" json:"id"`
	UserId string             `bson:"userId" json:"userId"`
	Name   string             `bson:"name" json:"name"`
}
