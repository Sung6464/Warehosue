package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Commodity represents a commodity in the database.
type Commodity struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name   string             `bson:"name" json:"name"`
	Amount int                `bson:"amount" json:"amount"`
}
