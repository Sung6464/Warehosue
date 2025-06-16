package controller

import (
	"api-gateway/config"
	"log" // Import log for debugging
	"net/http"
	"net/http/httputil"
	"net/url" // Import the strings package

	"github.com/gin-gonic/gin"
)

type GatewayController struct {
	// You might have dependencies here, e.g., for authentication or logging
}

// NewGatewayController creates a new GatewayController
func NewGatewayController() *GatewayController {
	return &GatewayController{}
}

// ProxyToCustomerService proxies requests to the Customer Service
func (gc *GatewayController) ProxyToCustomerService(c *gin.Context) {
	// The target service expects paths like /customers, /customers/:id, etc.
	gc.proxyRequest(c, config.Cfg.CustomerServiceURL, "/customers")
}

// ProxyToWarehouseService proxies requests to the Warehouse Service
func (gc *GatewayController) ProxyToWarehouseService(c *gin.Context) {
	// The target service expects paths like /warehouses, /warehouses/:id, etc.
	gc.proxyRequest(c, config.Cfg.WarehouseServiceURL, "/warehouses")
}

// ProxyToCommoditiesService proxies requests to the Commodities Service
func (gc *GatewayController) ProxyToCommoditiesService(c *gin.Context) {
	// The target service expects paths like /commodities, /commodities/:id, etc.
	gc.proxyRequest(c, config.Cfg.CommoditiesServiceURL, "/commodities")
}

// ProxyToInventoryService proxies requests to the Inventory Service
func (gc *GatewayController) ProxyToInventoryService(c *gin.Context) {
	// The target service expects paths like /inventory, /inventory/:id, etc.
	gc.proxyRequest(c, config.Cfg.InventoryServiceURL, "/inventory")
}

// proxyRequest handles the actual proxying logic
// targetBasePath is the base path that the *target service* expects (e.g., "/customers", "/commodities")
func (gc *GatewayController) proxyRequest(c *gin.Context, targetURL string, targetBasePath string) {
	remote, err := url.Parse(targetURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse target URL"})
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)

	// Custom director to rewrite the request URL
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.Host = remote.Host // Important for target service to receive correct Host header

		// Get the captured path from the Gin route (*proxyPath)
		// This will be "" for /api/commodities, "/" for /api/commodities/, "/123" for /api/commodities/123
		proxyPath := c.Param("proxyPath")

		// Construct the new path for the target service
		// We want: targetBasePath + proxyPath (e.g., "/commodities" + "/123" = "/commodities/123")
		// Special handling for when proxyPath is just "" or "/"
		if proxyPath == "" || proxyPath == "/" {
			req.URL.Path = targetBasePath // If no sub-path (e.g., /api/commodities or /api/commodities/), just use the base path
		} else {
			req.URL.Path = targetBasePath + proxyPath // Otherwise, append the captured sub-path
		}

		// Add logging for debugging the rewritten path
		log.Printf("Proxying request: Original Client Path: %s, Captured ProxyPath: '%s', Rewritten Target Path: '%s', Target Host: %s", c.Request.URL.Path, proxyPath, req.URL.Path, req.URL.Host)
	}

	// ServeHTTP will modify the request's URL and then forward it.
	proxy.ServeHTTP(c.Writer, c.Request)
}
