package model

// Customer represents a customer entity, now with Name and a list of associated WarehouseIDs.
type Customer struct {
	ID           string   `json:"id" bson:"_id"`
	Name         string   `json:"name" bson:"name"`
	WarehouseIDs []string `json:"warehouse_ids" bson:"warehouse_ids,omitempty"` // List of warehouse IDs a customer is associated with
}

// Warehouse represents the data received from the external Warehouse service.
// This struct is needed for deserializing responses when Customer Service calls Warehouse Service.
type Warehouse struct {
	ID      string `json:"id" bson:"_id"`
	Name    string `json:"name" bson:"name"`
	Address string `json:"location" bson:"address"`
	Storage int    `json:"storage" bson:"storage"`
	// Note: We don't need CustomerID here as this struct is for receiving from Warehouse Service
}
