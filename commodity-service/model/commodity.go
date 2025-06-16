package model

// Commodity represents a product or item stored in a warehouse.
type Commodity struct {
	ID     string `json:"id" bson:"_id"`
	Name   string `json:"name" bson:"name"`
	Amount int    `json:"amount" bson:"amount"` // Represents the quantity or amount of the commodity
}
