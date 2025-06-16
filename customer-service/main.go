package main

import (
	"context"
	"fmt"
	"log"

	"Customer-Services/config"     // Fixed import path
	"Customer-Services/controller" // Fixed import path
	"Customer-Services/database"   // Fixed import path
	"Customer-Services/repository" // Fixed import path
	"Customer-Services/routes"     // Fixed import path
	"Customer-Services/service"    // Fixed import path

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

	customerCollection := mongoClient.Database("wms_db").Collection("customers")

	customerRepo := repository.NewMongoCustomerRepository(customerCollection)
	customerService := service.NewCustomerService(customerRepo)
	customerController := controller.NewCustomerController(customerService)

	router := gin.Default()

	routes.SetupCustomerRoutes(router, customerController)

	port := fmt.Sprintf(":%s", config.Cfg.CustomerServicePort)
	fmt.Printf("Customer Service API listening on port %s...\n", port)
	log.Fatal(router.Run(port))
}
