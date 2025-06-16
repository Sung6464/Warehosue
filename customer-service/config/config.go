package config

import (
	"fmt"
	"os"
)

// AppConfig holds the application-wide configuration for the Customer Service.
type AppConfig struct {
	MongoURI            string
	CustomerServicePort string
	WarehouseServiceURL string // URL for the Warehouse Service (for validation)
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
	Cfg.CustomerServicePort = os.Getenv("CUSTOMER_SERVICE_PORT")
	if Cfg.CustomerServicePort == "" {
		Cfg.CustomerServicePort = "8087" // Default port for Customer Service
	}

	// URL for other services that this service needs to communicate with.
	Cfg.WarehouseServiceURL = os.Getenv("WAREHOUSE_SERVICE_URL") // Load environment variable for Warehouse Service URL
	if Cfg.WarehouseServiceURL == "" {
		Cfg.WarehouseServiceURL = "http://localhost:8085" // Default for Warehouse Service
	}

	return nil
}
