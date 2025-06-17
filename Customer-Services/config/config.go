package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds the application configuration for this microservice.
type Config struct {
	Port         int    `json:"port"`
	GinMode      string `json:"gin_mode"`
	MongoDBURI   string `json:"mongodb_uri"`
	DatabaseName string `json:"database_name"`
}

// Cfg is the global configuration instance.
var Cfg *Config

// LoadConfig loads configuration from environment variables or defaults.
func LoadConfig() error {
	Cfg = &Config{
		Port:         8087, // Default port for customer service
		GinMode:      "debug",
		MongoDBURI:   "mongodb://mongodb-wms:27017", // Default for Docker Compose local
		DatabaseName: "wms_customer_db",
	}

	// Override with environment variables if set (Render will set these)
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

	fmt.Printf("Customer Service Configuration: Port=%d, GinMode=%s, MongoDBURI=%s, DatabaseName=%s\n",
		Cfg.Port, Cfg.GinMode, Cfg.MongoDBURI, Cfg.DatabaseName)

	return nil
}
