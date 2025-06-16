package main

import (
	"context"
	"fmt"
	"log"

	"Warehouse-Services/config"
	"Warehouse-Services/controller"
	"Warehouse-Services/database" // Importing model here for database connection (e.g., if database/connect.go needs to know about models)
	"Warehouse-Services/repository"
	"Warehouse-Services/routes"
	"Warehouse-Services/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 2. Initialize MongoDB connection for THIS service
	mongoClient, err := database.ConnectDB(config.Cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting MongoDB client: %v", err)
		}
	}()

	// Get a handle to the specific MongoDB collection for warehouses.
	warehouseCollection := mongoClient.Database("wms_db").Collection("warehouses")

	// 3. Initialize Repository, Service, and Controller layers
	warehouseRepo := repository.NewMongoWarehouseRepository(warehouseCollection)
	warehouseService := service.NewWarehouseService(warehouseRepo)
	warehouseController := controller.NewWarehouseController(warehouseService)

	// 4. Set up Gin router
	router := gin.Default()

	// 5. Setup routes
	routes.SetupWarehouseRoutes(router, warehouseController)

	// 6. Start the HTTP server
	port := fmt.Sprintf(":%s", config.Cfg.WarehouseServicePort)
	fmt.Printf("Warehouse Service API listening on port %s...\n", port)
	log.Fatal(router.Run(port))
}
