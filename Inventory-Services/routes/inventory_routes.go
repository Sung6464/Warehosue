package routes

import (
	"Inventory-Services/controller"
	"Inventory-Services/service"
	"fmt"      // Import fmt for string formatting
	"net/http" // Import http for redirects

	"github.com/gin-gonic/gin"
)

// InventoryRoutes sets up the API routes for inventory operations.
func InventoryRoutes(router *gin.Engine) {
	inventoryController := controller.NewInventoryController(service.NewInventoryService())

	// Primary routes: define WITHOUT a trailing slash for collection endpoints
	inventoryGroup := router.Group("/inventory")
	{
		inventoryGroup.POST("", inventoryController.CreateInventory)     // Matches /inventory
		inventoryGroup.GET("", inventoryController.GetAllInventories)    // Matches /inventory
		inventoryGroup.GET("/:id", inventoryController.GetInventoryByID) // Matches /inventory/:id
		inventoryGroup.PUT("/:id", inventoryController.UpdateInventory)
		inventoryGroup.DELETE("/:id", inventoryController.DeleteInventory)
	}

	// Add explicit 301 redirects for paths that might come in WITH trailing slashes.
	router.GET("/inventory/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/inventory")
	})
	router.POST("/inventory/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/inventory")
	})
	router.PUT("/inventory/:id/", func(c *gin.Context) {
		id := c.Param("id")
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/inventory/%s", id))
	})
	router.DELETE("/inventory/:id/", func(c *gin.Context) {
		id := c.Param("id")
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/inventory/%s", id))
	})
	router.GET("/inventory/:id/", func(c *gin.Context) {
		id := c.Param("id")
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/inventory/%s", id))
	})
	// Add other methods if necessary
}
