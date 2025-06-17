package controller

import (
	"api-gateway/config"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings" // Ensure strings package is imported

	"github.com/gin-gonic/gin"
)

// GatewayController handles proxying requests to various microservices.
type GatewayController struct {
	CustomerServiceURL    *url.URL
	WarehouseServiceURL   *url.URL
	CommoditiesServiceURL *url.URL
	InventoryServiceURL   *url.URL
}

// NewGatewayController creates a new instance of GatewayController.
func NewGatewayController() *GatewayController {
	customerURL, err := url.Parse(config.Cfg.CustomerServiceURL)
	if err != nil {
		log.Fatalf("Invalid Customer Service URL: %v", err)
	}
	warehouseURL, err := url.Parse(config.Cfg.WarehouseServiceURL)
	if err != nil {
		log.Fatalf("Invalid Warehouse Service URL: %v", err)
	}
	commoditiesURL, err := url.Parse(config.Cfg.CommoditiesServiceURL)
	if err != nil {
		log.Fatalf("Invalid Commodities Service URL: %v", err)
	}
	inventoryURL, err := url.Parse(config.Cfg.InventoryServiceURL)
	if err != nil {
		log.Fatalf("Invalid Inventory Service URL: %v", err)
	}

	return &GatewayController{
		CustomerServiceURL:    customerURL,
		WarehouseServiceURL:   warehouseURL,
		CommoditiesServiceURL: commoditiesURL,
		InventoryServiceURL:   inventoryURL,
	}
}

// ProxyToService creates a generic reverse proxy for a given target URL.
// `apiPathPrefix` is the path on the API Gateway (e.g., "/api/customers")
// `downstreamRootPath` is the root path on the target service (e.g., "/customers")
func (gc *GatewayController) ProxyToService(targetURL *url.URL, apiPathPrefix string, downstreamRootPath string) gin.HandlerFunc {
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	proxy.Director = func(req *http.Request) {
		originalPath := req.URL.Path

		// Extract the portion of the path that comes after the API Gateway's prefix.
		// Gin's `*proxyPath` captures this, but accessing it directly from req.URL.Path
		// is robust.
		// Example: if originalPath is "/api/customers/123" and apiPathPrefix is "/api/customers"
		// then proxyPathSegment will be "/123".
		// If originalPath is "/api/customers" or "/api/customers/", proxyPathSegment will be "" or "/".
		proxyPathSegment := strings.TrimPrefix(originalPath, apiPathPrefix)

		// Normalize the proxyPathSegment: ensure it doesn't have leading/trailing slashes for concatenation,
		// unless it's just the root path.
		var finalDownstreamPath string
		if proxyPathSegment == "" || proxyPathSegment == "/" {
			// If client requested /api/customers or /api/customers/,
			// send just /customers to the downstream service.
			finalDownstreamPath = downstreamRootPath
		} else {
			// If client requested /api/customers/123, send /customers/123.
			// Trim potential leading slash from proxyPathSegment for correct concatenation.
			trimmedProxyPathSegment := strings.TrimPrefix(proxyPathSegment, "/")
			// Ensure downstreamRootPath does not have a trailing slash for concatenation
			cleanedDownstreamRootPath := strings.TrimSuffix(downstreamRootPath, "/")
			finalDownstreamPath = fmt.Sprintf("%s/%s", cleanedDownstreamRootPath, trimmedProxyPathSegment)
		}

		req.URL.Path = finalDownstreamPath
		req.URL.RawQuery = req.URL.RawQuery // Preserve original query parameters

		// Important: Set the Host header to the target service's host (e.g., "customer-service:8087")
		// This is crucial for Docker's internal DNS resolution.
		req.Host = targetURL.Host
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host

		log.Printf("Proxying request: Original Client Path: %s, Rewritten Target Path: '%s', Target Host: %s, Query: %s\n",
			originalPath, req.URL.Path, req.URL.Host, req.URL.RawQuery)
	}

	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		log.Printf("Proxy Error: %v for request %s %s\n", err, req.Method, req.URL.Path)
		rw.WriteHeader(http.StatusBadGateway)
		_, _ = rw.Write([]byte(fmt.Sprintf("Bad Gateway: Could not reach upstream service (%s) or proxy error occurred: %v", targetURL.String(), err)))
	}

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// ProxyToCustomerService proxies requests to the Customer Service.
func (gc *GatewayController) ProxyToCustomerService(c *gin.Context) {
	// API Gateway route is /api/customers/*proxyPath
	// Downstream route is /customers
	gc.ProxyToService(gc.CustomerServiceURL, "/api/customers", "/customers")(c)
}

// ProxyToWarehouseService proxies requests to the Warehouse Service.
func (gc *GatewayController) ProxyToWarehouseService(c *gin.Context) {
	gc.ProxyToService(gc.WarehouseServiceURL, "/api/warehouses", "/warehouses")(c)
}

// ProxyToCommoditiesService proxies requests to the Commodity Service.
func (gc *GatewayController) ProxyToCommoditiesService(c *gin.Context) {
	gc.ProxyToService(gc.CommoditiesServiceURL, "/api/commodities", "/commodities")(c)
}

// ProxyToInventoryService proxies requests to the Inventory Service.
func (gc *GatewayController) ProxyToInventoryService(c *gin.Context) {
	gc.ProxyToService(gc.InventoryServiceURL, "/api/inventory", "/inventory")(c)
}

// HealthCheck provides a simple health check endpoint.
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "API Gateway is healthy"})
}
