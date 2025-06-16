package main

import (
	"context"
	"fmt"
	"log"

	"commodity-service/config"     // Fixed import path
	"commodity-service/controller" // Fixed import path
	"commodity-service/database"   // Fixed import path
	"commodity-service/repository" // Fixed import path
	"commodity-service/routes"     // Fixed import path
	"commodity-service/service"    // Fixed import path

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

	commodityCollection := mongoClient.Database("wms_db").Collection("commodities")

	commodityRepo := repository.NewMongoCommodityRepository(commodityCollection)
	commodityService := service.NewCommodityService(commodityRepo)
	commodityController := controller.NewCommodityController(commodityService)

	router := gin.Default()

	routes.SetupCommodityRoutes(router, commodityController)

	port := fmt.Sprintf(":%s", config.Cfg.CommoditiesServicePort)
	fmt.Printf("Commodity Service API listening on port %s...\n", port)
	log.Fatal(router.Run(port))
}
