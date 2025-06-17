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
		customerGroup.POST("", customerController.CreateCustomer)     // Matches /customers
		customerGroup.GET("", customerController.GetAllCustomers)     // Matches /customers
		customerGroup.GET("/:id", customerController.GetCustomerByID) // Matches /customers/:id
		customerGroup.PUT("/:id", customerController.UpdateCustomer)
		customerGroup.DELETE("/:id", customerController.DeleteCustomer)
	}

	// Add explicit 301 redirects for paths that might come in WITH trailing slashes.
	// This ensures consistency and guides clients to the non-trailing-slash URL.
	router.GET("/customers/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/customers")
	})
	router.POST("/customers/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/customers")
	})
	router.PUT("/customers/:id/", func(c *gin.Context) {
		id := c.Param("id")
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/customers/%s", id))
	})
	router.DELETE("/customers/:id/", func(c *gin.Context) {
		id := c.Param("id")
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/customers/%s", id))
	})
	router.GET("/customers/:id/", func(c *gin.Context) {
		id := c.Param("id")
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/customers/%s", id))
	})
	// Add other methods if necessary (e.g., PATCH, HEAD, OPTIONS)
}
