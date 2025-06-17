package main

import (
	"api-gateway/config"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// GatewayController handles proxying requests.
type GatewayController struct{}

// NewGatewayController creates a new GatewayController instance.
func NewGatewayController() *GatewayController {
	return &GatewayController{}
}

// proxyRequest handles the actual proxying logic.
// targetServiceURL is the base URL of the downstream microservice (e.g., "http://customer-service:8087").
// downstreamServiceBasePath is the base path that the target microservice expects (e.g., "/customers").
// This function will append the `proxyPath` (the part captured by Gin's wildcard) to this base path,
// ensuring correct slash handling.
func (gc *GatewayController) proxyRequest(c *gin.Context, targetServiceURL, downstreamServiceBasePath string) {
	remote, err := url.Parse(targetServiceURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse target URL"})
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.Host = remote.Host // Important for target service to receive correct Host header

		// Get the specific part of the path captured by Gin's wildcard.
		// If incoming is /api/customers, c.Param("proxyPath") will be "".
		// If incoming is /api/customers/123, c.Param("proxyPath") will be "/123".
		proxyPathSegment := c.Param("proxyPath")

		// Normalize the downstream base path: ensure no trailing slash.
		// e.g., "/customers" (not "/customers/")
		normalizedDownstreamBasePath := strings.TrimSuffix(downstreamServiceBasePath, "/")

		// Normalize the proxy path segment: ensure no leading slash.
		// e.g., "123" (not "/123") if it's an ID
		normalizedProxyPathSegment := strings.TrimPrefix(proxyPathSegment, "/")

		// Construct the final path for the downstream service.
		// If there's a non-empty segment from the proxyPath (e.g., an ID), append it with a slash.
		// Otherwise, just use the normalized base path.
		if normalizedProxyPathSegment != "" {
			req.URL.Path = fmt.Sprintf("%s/%s", normalizedDownstreamBasePath, normalizedProxyPathSegment)
		} else {
			req.URL.Path = normalizedDownstreamBasePath
		}

		req.URL.RawQuery = c.Request.URL.RawQuery // Preserve query parameters

		// Log the rewritten path for debugging
		log.Printf("Proxying request: Original Client Path: %s, Rewritten Target Path: '%s', Target Host: %s", c.Request.URL.Path, req.URL.Path, req.URL.Host)
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Proxy error for %s %s: %v", r.Method, r.URL.Path, err)
		if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "no such host") {
			c.JSON(http.StatusBadGateway, gin.H{"message": "Service unavailable or not reachable", "error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal proxy error", "error": err.Error()})
		}
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}

// Proxy functions for each service.
// Each passes the *exact base path* that its corresponding downstream microservice expects,
// e.g., "/customers" for the customer service.
func (gc *GatewayController) ProxyToCustomerService(c *gin.Context) {
	gc.proxyRequest(c, config.Cfg.CustomerServiceURL, "/customers")
}
func (gc *GatewayController) ProxyToWarehouseService(c *gin.Context) {
	gc.proxyRequest(c, config.Cfg.WarehouseServiceURL, "/warehouses")
}
func (gc *GatewayController) ProxyToCommoditiesService(c *gin.Context) {
	gc.proxyRequest(c, config.Cfg.CommoditiesServiceURL, "/commodities")
}
func (gc *GatewayController) ProxyToInventoryService(c *gin.Context) {
	gc.proxyRequest(c, config.Cfg.InventoryServiceURL, "/inventory")
}

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	gin.SetMode(config.Cfg.GinMode)
	router := gin.Default()

	// --- NO CORS MIDDLEWARE INCLUDED ---
	// This version is designed for backend-only testing (e.g., via Postman) where CORS is not relevant.

	gatewayController := NewGatewayController()

	// Health check endpoint for the API Gateway itself
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "API Gateway is healthy!"})
	})

	// Group API routes under "/api" prefix.
	// The wildcard *proxyPath will capture everything after /api/<service_name>.
	// E.g., for /api/customers/*proxyPath:
	// - If incoming is /api/customers, proxyPath will be ""
	// - If incoming is /api/customers/123, proxyPath will be "/123"
	apiGroup := router.Group("/api")
	{
		apiGroup.Any("/customers/*proxyPath", gatewayController.ProxyToCustomerService)
		apiGroup.Any("/warehouses/*proxyPath", gatewayController.ProxyToWarehouseService)
		apiGroup.Any("/commodities/*proxyPath", gatewayController.ProxyToCommoditiesService)
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
		log.Printf("--- WMS API Gateway: Ultra-Precise Routing (No CORS) ---") // Updated unique log
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
