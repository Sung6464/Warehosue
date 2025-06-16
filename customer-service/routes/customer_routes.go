package routes

import (
	"Customer-Services/controller" // Fixed import path

	"github.com/gin-gonic/gin"
)

// SetupCustomerRoutes sets up the API routes for customer operations.
func SetupCustomerRoutes(router *gin.Engine, customerController *controller.CustomerController) {
	customerRoutes := router.Group("/customers")
	{
		customerRoutes.POST("", customerController.CreateCustomer)
		customerRoutes.GET("", customerController.GetAllCustomers) // Handles optional ?warehouse_id= filter
		customerRoutes.GET("/:id", customerController.GetCustomer)
		customerRoutes.PUT("/:id", customerController.UpdateCustomer)
		// New routes for many-to-many customer-warehouse mapping
		// POST /customers/:id/warehouses/:warehouse_id to add a warehouse to a customer
		customerRoutes.POST("/:id/warehouses/:warehouse_id", customerController.AddWarehouseToCustomer)
		// DELETE /customers/:id/warehouses/:warehouse_id to remove a warehouse from a customer
		customerRoutes.DELETE("/:id/warehouses/:warehouse_id", customerController.RemoveWarehouseFromCustomer)
	}
}
