package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds the application configuration for this microservice.
type Config struct { // Ensure struct is named Config
	Port         int    `json:"port"`
	GinMode      string `json:"gin_mode"` // This field must exist
	MongoDBURI   string `json:"mongodb_uri"`
	DatabaseName string `json:"database_name"`
}

// Cfg is the global configuration instance.
var Cfg *Config

// LoadConfig loads configuration from environment variables or defaults.
func LoadConfig() error {
	Cfg = &Config{
		Port:         8086,                        // Default port for commodity service
		GinMode:      "debug",                     // Default value
		MongoDBURI:   "mongodb://localhost:27017", // For individual testing outside Docker
		DatabaseName: "wms_commodities_db",
	}

	if portStr := os.Getenv("PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			Cfg.Port = port
		}
	}
	if ginMode := os.Getenv("GIN_MODE"); ginMode != "" {
		Cfg.GinMode = ginMode
	}
	if mongoURI := os.Getenv("MONGODB_URI"); mongoURI != "" {
		Cfg.MongoDBURI = mongoURI
	}
	if dbName := os.Getenv("DATABASE_NAME"); dbName != "" {
		Cfg.DatabaseName = dbName
	}

	fmt.Printf("Commodity Service Configuration: Port=%d, GinMode=%s, MongoDBURI=%s, DatabaseName=%s\n",
		Cfg.Port, Cfg.GinMode, Cfg.MongoDBURI, Cfg.DatabaseName)

	return nil
}
