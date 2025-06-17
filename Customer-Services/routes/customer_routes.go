package routes

import (
	"Customer-Services/controller"
	"Customer-Services/service"
	"fmt"      // Import fmt for string formatting
	"net/http" // Import http for redirects

	"github.com/gin-gonic/gin"
)

// CustomerRoutes sets up the API routes for customer operations.
func CustomerRoutes(router *gin.Engine) {
	customerController := controller.NewCustomerController(service.NewCustomerService())

	// Primary routes: define WITHOUT a trailing slash for collection endpoints
	customerGroup := router.Group("/customers")
	{
		// Explicitly handle all HTTP methods for the base /customers path (no trailing slash)
		customerGroup.POST("", customerController.CreateCustomer) // Matches /customers
		customerGroup.GET("", customerController.GetAllCustomers) // Matches /customers

		// Routes for specific IDs
		customerGroup.GET("/:id", customerController.GetCustomerByID) // Matches /customers/:id
		customerGroup.PUT("/:id", customerController.UpdateCustomer)
		customerGroup.DELETE("/:id", customerController.DeleteCustomer)
	}

	// Add explicit 301 redirects for paths that might come in WITH trailing slashes.
	// This ensures consistency and guides clients (and proxies) to the non-trailing-slash URL.
	router.GET("/customers/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/customers")
	})
	router.POST("/customers/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/customers")
	})
	router.PUT("/customers/:id/", func(c *gin.Context) { // Handle /customers/:id/
		id := c.Param("id")
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/customers/%s", id))
	})
	router.DELETE("/customers/:id/", func(c *gin.Context) { // Handle /customers/:id/
		id := c.Param("id")
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/customers/%s", id))
	})
	router.GET("/customers/:id/", func(c *gin.Context) { // Explicitly redirect GET /customers/:id/
		id := c.Param("id")
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/customers/%s", id))
	})
	// Add OPTIONS redirect for CORS preflight if a trailing slash is sent
	router.OPTIONS("/customers/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/customers")
	})
	router.OPTIONS("/customers/:id/", func(c *gin.Context) {
		id := c.Param("id")
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/customers/%s", id))
	})
}
