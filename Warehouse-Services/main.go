package main

import (
	"Warehouse-Services/config"
	"Warehouse-Services/database"
	"Warehouse-Services/routes"
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

	client, err := database.ConnectDB() // ConnectDB should take no arguments here
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatalf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	gin.SetMode(config.Cfg.GinMode)
	router := gin.Default()

	// --- CORS Configuration for Warehouse Service ---
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Register warehouse-specific routes
	routes.WarehouseRoutes(router) // Correct function name

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Cfg.Port), // Use config.Cfg.Port
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("--- Warehouse Service (Warehouse-Services): Listening on port :%d with CORS ---", config.Cfg.Port) // Use config.Cfg.Port
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on port %d: %v\n", config.Cfg.Port, err) // Use config.Cfg.Port
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1) // Corrected: make it a channel
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan // Corrected: receive from channel

	log.Println("Warehouse Service (Warehouse-Services) shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Warehouse Service (Warehouse-Services) graceful shutdown failed: %v\n", err)
	}
	log.Println("Warehouse Service (Warehouse-Services) stopped.")
}
