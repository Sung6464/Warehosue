package main

import (
	"context"
	"fmt"
	"log"

	"Inventory-Services/config"     // FIXED: Matches go.mod module name
	"Inventory-Services/controller" // FIXED: Matches go.mod module name
	"Inventory-Services/database"   // FIXED: Matches go.mod module name
	"Inventory-Services/repository" // FIXED: Matches go.mod module name
	"Inventory-Services/routes"     // FIXED: Matches go.mod module name
	"Inventory-Services/service"    // FIXED: Matches go.mod module name

	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	mongoClient, err := database.ConnectDB(config.Cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting MongoDB client: %v", err)
		}
	}()

	inventoryCollection := mongoClient.Database("wms_db").Collection("inventory")

	inventoryRepo := repository.NewMongoInventoryRepository(inventoryCollection)
	inventoryService := service.NewInventoryService(inventoryRepo)
	inventoryController := controller.NewInventoryController(inventoryService)

	router := gin.Default()

	routes.SetupInventoryRoutes(router, inventoryController)

	port := fmt.Sprintf(":%s", config.Cfg.InventoryServicePort)
	fmt.Printf("Inventory Service API listening on port %s...\n", port)
	log.Fatal(router.Run(port))
}
