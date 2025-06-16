package routes // The package for files in the 'routes' folder

import (
	"Inventory-Services/controller" // FIXED: Matches go.mod module name

	"github.com/gin-gonic/gin"
)

// SetupInventoryRoutes sets up the API routes for inventory operations.
func SetupInventoryRoutes(router *gin.Engine, inventoryController *controller.InventoryController) {
	inventoryRoutes := router.Group("/inventory")
	{
		inventoryRoutes.POST("", inventoryController.CreateInventoryItem)
		inventoryRoutes.GET("", inventoryController.GetAllInventoryItems)
		inventoryRoutes.GET("/:id", inventoryController.GetInventoryItem)
		inventoryRoutes.PUT("/:id", inventoryController.UpdateInventoryItem)
		inventoryRoutes.POST("/:id/adjust", inventoryController.AdjustInventoryQuantity)
		inventoryRoutes.GET("/warehouses/:id", inventoryController.GetAllInventoryItems)
	}
}
