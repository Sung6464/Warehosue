package config

import (
	"fmt"
	"os"
	"strconv" // Import for string to int conversion
)

// AppConfig holds the application-wide configuration
type AppConfig struct {
	APIGatewayPort        int    // Port for the API Gateway
	GinMode               string // Gin mode (debug or release)
	CustomerServiceURL    string
	WarehouseServiceURL   string
	CommoditiesServiceURL string
	InventoryServiceURL   string
}

// Cfg is the global configuration instance
var Cfg AppConfig

// LoadConfig loads configuration from environment variables,
// falling back to localhost defaults for local execution.
func LoadConfig() error {
	// API Gateway Port
	portStr := os.Getenv("API_GATEWAY_PORT")
	if portStr == "" {
		portStr = "8080" // Default port for API Gateway
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid API_GATEWAY_PORT: %v", err)
	}
	Cfg.APIGatewayPort = port

	// Gin Mode
	Cfg.GinMode = os.Getenv("GIN_MODE")
	if Cfg.GinMode == "" {
		Cfg.GinMode = "debug" // Default Gin mode
	}

	// Service URLs - IMPORTANT: Use localhost for direct local execution!
	Cfg.CustomerServiceURL = os.Getenv("CUSTOMER_SERVICE_URL")
	if Cfg.CustomerServiceURL == "" {
		Cfg.CustomerServiceURL = "http://localhost:8087"
	}

	Cfg.WarehouseServiceURL = os.Getenv("WAREHOUSE_SERVICE_URL")
	if Cfg.WarehouseServiceURL == "" {
		Cfg.WarehouseServiceURL = "http://localhost:8085"
	}

	Cfg.CommoditiesServiceURL = os.Getenv("COMMODITIES_SERVICE_URL")
	if Cfg.CommoditiesServiceURL == "" {
		Cfg.CommoditiesServiceURL = "http://localhost:8086" // This was the problem!
	}

	Cfg.InventoryServiceURL = os.Getenv("INVENTORY_SERVICE_URL")
	if Cfg.InventoryServiceURL == "" {
		Cfg.InventoryServiceURL = "http://localhost:8088"
	}

	return nil
}
