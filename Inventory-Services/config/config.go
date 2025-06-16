package config

import (
	"fmt"
	"os"
)

// AppConfig holds the application-wide configuration for the Inventory Service.
type AppConfig struct {
	MongoURI              string
	InventoryServicePort  string
	WarehouseServiceURL   string // FIXED: Added this field
	CommoditiesServiceURL string // FIXED: Added this field
	CustomerServiceURL    string // FIXED: Added this field
}

// Global variable to hold the loaded configuration.
var Cfg AppConfig

// LoadConfig loads configuration from environment variables.
func LoadConfig() error {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017" // Default for local development
		fmt.Println("MONGO_URI environment variable not set, using default: mongodb://localhost:27017")
	}
	Cfg.MongoURI = mongoURI

	// This service's port (read from environment variable or use default)
	Cfg.InventoryServicePort = os.Getenv("INVENTORY_SERVICE_PORT")
	if Cfg.InventoryServicePort == "" {
		Cfg.InventoryServicePort = "8088" // Default port for Inventory Service
	}

	// URLs for other services that this service needs to communicate with for validation.
	Cfg.WarehouseServiceURL = os.Getenv("WAREHOUSE_SERVICE_URL")
	if Cfg.WarehouseServiceURL == "" {
		Cfg.WarehouseServiceURL = "http://localhost:8085" // Default for Warehouse Service
	}
	Cfg.CommoditiesServiceURL = os.Getenv("COMMODITIES_SERVICE_URL")
	if Cfg.CommoditiesServiceURL == "" {
		Cfg.CommoditiesServiceURL = "http://localhost:8086" // Default for Commodities Service
	}
	Cfg.CustomerServiceURL = os.Getenv("CUSTOMER_SERVICE_URL")
	if Cfg.CustomerServiceURL == "" {
		Cfg.CustomerServiceURL = "http://localhost:8087" // Default for Customer Service
	}

	return nil
}
