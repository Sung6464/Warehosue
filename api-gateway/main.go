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
	// This ensures our explicit routing and proxy.Director have full control over path normalization.
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
		// --- CUSTOMER SERVICE ROUTES ---
		// Explicitly define routes for the root collection path (e.g., /api/customers)
		// This handles GET/POST/PUT/DELETE/OPTIONS requests directly to /api/customers
		apiGroup.GET("/customers", gatewayController.ProxyToCustomerService)
		apiGroup.POST("/customers", gatewayController.ProxyToCustomerService)
		apiGroup.PUT("/customers", gatewayController.ProxyToCustomerService)     // Not typically used for root, but covers all methods
		apiGroup.DELETE("/customers", gatewayController.ProxyToCustomerService)  // Not typically used for root, but covers all methods
		apiGroup.OPTIONS("/customers", gatewayController.ProxyToCustomerService) // For CORS preflight

		// Use Any for paths with a wildcard (e.g., /api/customers/123)
		apiGroup.Any("/customers/*proxyPath", gatewayController.ProxyToCustomerService)

		// --- WAREHOUSE SERVICE ROUTES ---
		apiGroup.GET("/warehouses", gatewayController.ProxyToWarehouseService)
		apiGroup.POST("/warehouses", gatewayController.ProxyToWarehouseService)
		apiGroup.PUT("/warehouses", gatewayController.ProxyToWarehouseService)
		apiGroup.DELETE("/warehouses", gatewayController.ProxyToWarehouseService)
		apiGroup.OPTIONS("/warehouses", gatewayController.ProxyToWarehouseService)

		apiGroup.Any("/warehouses/*proxyPath", gatewayController.ProxyToWarehouseService)

		// --- COMMODITY SERVICE ROUTES ---
		apiGroup.GET("/commodities", gatewayController.ProxyToCommoditiesService)
		apiGroup.POST("/commodities", gatewayController.ProxyToCommoditiesService)
		apiGroup.PUT("/commodities", gatewayController.ProxyToCommoditiesService)
		apiGroup.DELETE("/commodities", gatewayController.ProxyToCommoditiesService)
		apiGroup.OPTIONS("/commodities", gatewayController.ProxyToCommoditiesService)

		apiGroup.Any("/commodities/*proxyPath", gatewayController.ProxyToCommoditiesService)

		// --- INVENTORY SERVICE ROUTES ---
		apiGroup.GET("/inventory", gatewayController.ProxyToInventoryService)
		apiGroup.POST("/inventory", gatewayController.ProxyToInventoryService)
		apiGroup.PUT("/inventory", gatewayController.ProxyToInventoryService)
		apiGroup.DELETE("/inventory", gatewayController.ProxyToInventoryService)
		apiGroup.OPTIONS("/inventory", gatewayController.ProxyToInventoryService)

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
