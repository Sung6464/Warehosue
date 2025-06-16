package model // The package for files in the 'model' folder

import "time" // Required for LastUpdated field

// InventoryItem represents a specific commodity stored in a specific warehouse, potentially for a customer.
type InventoryItem struct {
	ID          string    `json:"id" bson:"_id"`
	WarehouseID string    `json:"warehouse_id" bson:"warehouse_id"`         // ID of the warehouse where this commodity is stored
	CommodityID string    `json:"commodity_id" bson:"commodity_id"`         // ID of the commodity itself
	CustomerID  string    `json:"customer_id" bson:"customer_id,omitempty"` // Optional: If this stock belongs to a specific customer
	Quantity    int       `json:"quantity" bson:"quantity"`                 // Current quantity of this commodity at this location
	LastUpdated time.Time `json:"last_updated" bson:"last_updated"`         // Timestamp of the last update
}

// Commodity represents the data received from the external Commodities service.
type Commodity struct {
	ID     string `json:"id" bson:"_id"`
	Name   string `json:"name" bson:"name"`
	Amount int    `json:"amount" bson:"amount"`
}

// Warehouse represents the data received from the external Warehouse service.
type Warehouse struct {
	ID      string `json:"id" bson:"_id"`
	Name    string `json:"name" bson:"name"`
	Address string `json:"location" bson:"address"`
	Storage int    `json:"storage" bson:"storage"`
}

// Customer represents the data received from the external Customers service.
type Customer struct {
	ID           string   `json:"id" bson:"_id"`
	Name         string   `json:"name" bson:"name"`
	WarehouseIDs []string `json:"warehouse_ids" bson:"warehouse_ids,omitempty"`
}
