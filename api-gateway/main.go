package main

import (
	"api-gateway/config"
	"api-gateway/controller"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	gin.SetMode(config.Cfg.GinMode)
	router := gin.Default()

	// CRITICAL: Disable Gin's automatic trailing slash redirects and fixed path redirects on API Gateway.
	// This makes our proxy.Director fully responsible for path normalization and prevents undesired 301/307s.
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	// --- Robust CORS Configuration for API Gateway ---
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"}, // Allow React dev server origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	gatewayController := controller.NewGatewayController()

	// Health check endpoint for the API Gateway itself
	router.GET("/health", controller.HealthCheck)

	// Group API routes under "/api" prefix.
	apiGroup := router.Group("/api")
	{
		// Inside apiGroup
		apiGroup.POST("/customers", gatewayController.ProxyToCustomerService)           // <-- ADD THIS LINE for explicit POST to root collection
		apiGroup.Any("/customers/*proxyPath", gatewayController.ProxyToCustomerService) // Keep this for other methods and subpaths
		// You will need to do this for POST for ALL services that receive POST to their root collection
		apiGroup.POST("/warehouses", gatewayController.ProxyToWarehouseService)
		apiGroup.Any("/warehouses/*proxyPath", gatewayController.ProxyToWarehouseService)

		apiGroup.POST("/commodities", gatewayController.ProxyToCommoditiesService)
		apiGroup.Any("/commodities/*proxyPath", gatewayController.ProxyToCommoditiesService)

		apiGroup.POST("/inventory", gatewayController.ProxyToInventoryService)
		apiGroup.Any("/inventory/*proxyPath", gatewayController.ProxyToInventoryService)
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Cfg.APIGatewayPort),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("--- WMS API Gateway: CORS Configured for Frontend Operations ---")
		log.Printf("API Gateway listening on port :%d...", config.Cfg.APIGatewayPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on port %d: %v\n", config.Cfg.APIGatewayPort, err)
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down API Gateway gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("API Gateway graceful shutdown failed: %v\n", err)
	}
	log.Println("API Gateway stopped.")
}
