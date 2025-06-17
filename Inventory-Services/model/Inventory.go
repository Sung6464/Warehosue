package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Inventory represents an inventory record in the database.
type Inventory struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ProductID   primitive.ObjectID `bson:"product_id" json:"productId"`
	Quantity    int                `bson:"quantity" json:"quantity"`
	Location    string             `bson:"location" json:"location"`
	LastUpdated time.Time          `bson:"last_updated" json:"lastUpdated"`
}
