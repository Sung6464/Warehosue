package model

import "time" // Required for InventoryItem's LastUpdated field

// Warehouse represents a warehouse entity in the WMS.
type Warehouse struct {
	ID          string `json:"id" bson:"_id"`
	Name        string `json:"name" bson:"name"`
	Address     string `json:"location" bson:"address"` // Mapped from 'location' in JSON
	Storage     int    `json:"storage" bson:"storage"`
	CommodityID string `json:"commodity_id" bson:"commodity_id,omitempty"` // Link to Commodity (used for examples)
	CustomerID  string `json:"customer_id" bson:"customer_id,omitempty"`   // New: Link to the booking Customer (for 1-to-many relationship)
}

// Commodity represents the data received from the external Commodities service.
// This is a local representation of the *contract* of the Commodity service's response.
type Commodity struct {
	ID     string `json:"id" bson:"_id"`
	Name   string `json:"name" bson:"name"`
	Amount int    `json:"amount" bson:"amount"`
}

// Customer represents the data received from the external Customers service.
// This is a local representation of the *contract* of the Customer service's response.
type Customer struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
	// Note: WarehouseIDs are managed by the Customer Service.
	// We only need ID and Name here for validation purposes.
}

// InventoryItem represents data received from the Inventory Service.
// This is a local representation of the *contract* of the Inventory service's response.
type InventoryItem struct {
	ID          string    `json:"id" bson:"_id"`
	WarehouseID string    `json:"warehouse_id" bson:"warehouse_id"`
	CommodityID string    `json:"commodity_id" bson:"commodity_id"`
	CustomerID  string    `json:"customer_id" bson:"customer_id,omitempty"` // Optional customer for specific stock
	Quantity    int       `json:"quantity" bson:"quantity"`
	LastUpdated time.Time `json:"last_updated" bson:"last_updated"`
}
