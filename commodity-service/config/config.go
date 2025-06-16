package config

import (
	"fmt"
	"os"
)

type AppConfig struct {
	MongoURI               string
	CommoditiesServicePort string
}

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
	Cfg.CommoditiesServicePort = os.Getenv("COMMODITIES_SERVICE_PORT")
	if Cfg.CommoditiesServicePort == "" {
		Cfg.CommoditiesServicePort = "8086" // Default port for Commodity Service
	}

	return nil
}
