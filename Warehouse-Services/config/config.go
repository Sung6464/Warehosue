package config

import (
	"fmt"
	"os"
)

// AppConfig holds the application-wide configuration.
type AppConfig struct {
	MongoURI              string
	WarehouseServicePort  string
	CommoditiesServiceURL string
	CustomerServiceURL    string
	InventoryServiceURL   string
}

// Global variable to hold the loaded configuration.
var Cfg AppConfig

// LoadConfig loads configuration from environment variables.
func LoadConfig() error {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017" // Default for local dev
		fmt.Println("MONGO_URI environment variable not set, using default: mongodb://localhost:27017")
	}
	Cfg.MongoURI = mongoURI

	// This service's port
	Cfg.WarehouseServicePort = os.Getenv("WAREHOUSE_SERVICE_PORT")
	if Cfg.WarehouseServicePort == "" {
		Cfg.WarehouseServicePort = "8085" // Default port for Warehouse Service
	}

	// Other services' URLs for inter-service communication
	Cfg.CommoditiesServiceURL = os.Getenv("COMMODITIES_SERVICE_URL")
	if Cfg.CommoditiesServiceURL == "" {
		Cfg.CommoditiesServiceURL = "http://localhost:8086" // Default
	}
	Cfg.CustomerServiceURL = os.Getenv("CUSTOMER_SERVICE_URL")
	if Cfg.CustomerServiceURL == "" {
		Cfg.CustomerServiceURL = "http://localhost:8087" // Default
	}
	Cfg.InventoryServiceURL = os.Getenv("INVENTORY_SERVICE_URL")
	if Cfg.InventoryServiceURL == "" {
		Cfg.InventoryServiceURL = "http://localhost:8088" // Default for Inventory Service
	}

	return nil
}
