package routes

import (
	"Warehouse-Services/controller"

	"github.com/gin-gonic/gin"
)

// SetupWarehouseRoutes sets up the API routes for warehouse operations.
func SetupWarehouseRoutes(router *gin.Engine, warehouseController *controller.WarehouseController) {
	warehouseRoutes := router.Group("/warehouses")
	{
		warehouseRoutes.POST("", warehouseController.CreateWarehouse)
		warehouseRoutes.GET("", warehouseController.GetAllWarehouses)
		warehouseRoutes.GET("/:id", warehouseController.GetWarehouse)
		warehouseRoutes.PUT("/:id", warehouseController.UpdateWarehouse)
		warehouseRoutes.GET("/:id/inventory", warehouseController.GetInventoryInWarehouse) // Get inventory in a specific warehouse

		// Routes for managing customer booking (1-to-one relationship from Warehouse perspective)
		// POST /warehouses/:id/book/:customer_id to book a warehouse for a customer
		warehouseRoutes.POST("/:id/book/:customer_id", warehouseController.BookWarehouse)
		// DELETE /warehouses/:id/unbook/:customer_id to unbook a warehouse from a customer
		warehouseRoutes.DELETE("/:id/unbook/:customer_id", warehouseController.UnbookWarehouse)
	}
}
