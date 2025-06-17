package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Warehouse represents a warehouse in the database.
type Warehouse struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string             `bson:"name" json:"name"`
	Location string             `bson:"location" json:"location"`
	Storage  int                `bson:"storage" json:"storage"` // Capacity in some unit
}
