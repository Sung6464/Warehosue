package routes

import (
	"Warehouse-Services/controller"
	"Warehouse-Services/service"
	"fmt"      // Import fmt for string formatting
	"net/http" // Import http for redirects

	"github.com/gin-gonic/gin"
)

// WarehouseRoutes sets up the API routes for warehouse operations.
func WarehouseRoutes(router *gin.Engine) {
	warehouseController := controller.NewWarehouseController(service.NewWarehouseService())

	// Primary routes: define WITHOUT a trailing slash for collection endpoints
	warehouseGroup := router.Group("/warehouses")
	{
		warehouseGroup.POST("", warehouseController.CreateWarehouse)     // Matches /warehouses
		warehouseGroup.GET("", warehouseController.GetAllWarehouses)     // Matches /warehouses
		warehouseGroup.GET("/:id", warehouseController.GetWarehouseByID) // Matches /warehouses/:id
		warehouseGroup.PUT("/:id", warehouseController.UpdateWarehouse)
		warehouseGroup.DELETE("/:id", warehouseController.DeleteWarehouse)
	}

	// Add explicit 301 redirects for paths that might come in WITH trailing slashes.
	router.GET("/warehouses/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/warehouses")
	})
	router.POST("/warehouses/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/warehouses")
	})
	router.PUT("/warehouses/:id/", func(c *gin.Context) {
		id := c.Param("id")
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/warehouses/%s", id))
	})
	router.DELETE("/warehouses/:id/", func(c *gin.Context) {
		id := c.Param("id")
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/warehouses/%s", id))
	})
	router.GET("/warehouses/:id/", func(c *gin.Context) {
		id := c.Param("id")
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/warehouses/%s", id))
	})
	// Add other methods if necessary
}
