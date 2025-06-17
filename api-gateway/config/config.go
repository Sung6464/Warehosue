package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// Config holds the application configuration.
type Config struct {
	APIGatewayPort        int    `json:"api_gateway_port"`
	GinMode               string `json:"gin_mode"`
	CustomerServiceURL    string `json:"customer_service_url"`
	WarehouseServiceURL   string `json:"warehouse_service_url"`
	CommoditiesServiceURL string `json:"commodities_service_url"`
	InventoryServiceURL   string `json:"inventory_service_url"`
}

// Cfg is the global configuration instance.
var Cfg *Config

// LoadConfig loads configuration from environment variables or defaults.
func LoadConfig() error {
	Cfg = &Config{
		APIGatewayPort: 8080, // Default for local
		GinMode:        "debug",
		// Default to Docker Compose service names for inter-service communication
		// These will be overridden by environment variables if set (e.g., from docker-compose.yml or Render)
		CustomerServiceURL:    "http://customer-service:8087",
		WarehouseServiceURL:   "http://warehouse-service:8085",
		CommoditiesServiceURL: "http://commodity-service:8086",
		InventoryServiceURL:   "http://inventory-service:8088",
	}

	// Override with environment variables if set
	if portStr := os.Getenv("PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			Cfg.APIGatewayPort = port
		}
	}
	if ginMode := os.Getenv("GIN_MODE"); ginMode != "" {
		Cfg.GinMode = ginMode
	}
	if customerURL := os.Getenv("CUSTOMER_SERVICE_URL"); customerURL != "" {
		Cfg.CustomerServiceURL = customerURL
	}
	if warehouseURL := os.Getenv("WAREHOUSE_SERVICE_URL"); warehouseURL != "" {
		Cfg.WarehouseServiceURL = warehouseURL
	}
	if commoditiesURL := os.Getenv("COMMODITIES_SERVICE_URL"); commoditiesURL != "" {
		Cfg.CommoditiesServiceURL = commoditiesURL
	}
	if inventoryURL := os.Getenv("INVENTORY_SERVICE_URL"); inventoryURL != "" {
		Cfg.InventoryServiceURL = inventoryURL
	}

	configJSON, _ := json.MarshalIndent(Cfg, "", "  ")
	fmt.Printf("API Gateway Configuration:\n%s\n", string(configJSON))

	return nil
}
