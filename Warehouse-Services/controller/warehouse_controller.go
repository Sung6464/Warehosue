package controller

import (
	"context" // Explicit import for context.WithTimeout
	"fmt"
	"net/http"
	"time"

	"Warehouse-Services/model"
	"Warehouse-Services/service"

	"github.com/gin-gonic/gin"
)

// WarehouseController handles HTTP requests for warehouses.
type WarehouseController struct {
	service service.WarehouseService
}

// NewWarehouseController creates a new instance of WarehouseController.
func NewWarehouseController(s service.WarehouseService) *WarehouseController {
	return &WarehouseController{
		service: s,
	}
}

// CreateWarehouse handles POST requests to create a new warehouse.
func (ctrl *WarehouseController) CreateWarehouse(c *gin.Context) {
	var newWarehouse model.Warehouse
	if err := c.ShouldBindJSON(&newWarehouse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err.Error())})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	createdWarehouse, err := ctrl.service.CreateWarehouse(ctx, newWarehouse)
	if err != nil {
		if err.Error() == "warehouse name is required" || err.Error() == "warehouse location (address) is required" ||
			(len(err.Error()) >= 20 && err.Error()[0:20] == "invalid commodity ID") ||
			(len(err.Error()) >= 20 && err.Error()[0:20] == "invalid customer ID") ||
			err.Error() == "warehouse is already booked by another customer" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create warehouse: %v", err.Error())})
		}
		return
	}

	c.JSON(http.StatusCreated, createdWarehouse)
}

// GetWarehouse handles GET requests to retrieve a single warehouse by ID.
func (ctrl *WarehouseController) GetWarehouse(c *gin.Context) {
	warehouseID := c.Param("id")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	warehouse, err := ctrl.service.GetWarehouseByID(ctx, warehouseID)
	if err != nil {
		if err.Error() == "warehouse not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch warehouse: %v", err.Error())})
		}
		return
	}

	c.JSON(http.StatusOK, warehouse)
}

// GetAllWarehouses handles GET requests to retrieve all warehouses.
func (ctrl *WarehouseController) GetAllWarehouses(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	warehouses, err := ctrl.service.GetAllWarehouses(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch all warehouses: %v", err.Error())})
		return
	}

	c.JSON(http.StatusOK, warehouses)
}

// UpdateWarehouse handles PUT requests to update an existing warehouse by ID.
func (ctrl *WarehouseController) UpdateWarehouse(c *gin.Context) {
	warehouseID := c.Param("id")

	var updatedData map[string]interface{}
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err.Error())})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	updatedWarehouse, err := ctrl.service.UpdateWarehouse(ctx, warehouseID, updatedData)
	if err != nil {
		if err.Error() == "no fields provided for update" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if err.Error() == "warehouse not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if (len(err.Error()) >= 20 && err.Error()[0:20] == "invalid commodity ID") ||
			(len(err.Error()) >= 20 && err.Error()[0:20] == "invalid customer ID") ||
			err.Error() == "warehouse is already booked by another customer" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update warehouse: %v", err.Error())})
		}
		return
	}

	c.JSON(http.StatusOK, updatedWarehouse)
}

// GetInventoryInWarehouse handles GET requests to retrieve all inventory items for a specific warehouse.
func (ctrl *WarehouseController) GetInventoryInWarehouse(c *gin.Context) {
	warehouseID := c.Param("id")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	inventoryItems, err := ctrl.service.GetInventoryInWarehouse(ctx, warehouseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch inventory for warehouse: %v", err.Error())})
		return
	}

	c.JSON(http.StatusOK, inventoryItems)
}

// BookWarehouse handles POST requests to book a warehouse for a specific customer.
func (ctrl *WarehouseController) BookWarehouse(c *gin.Context) {
	warehouseID := c.Param("id")
	customerID := c.Param("customer_id")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	updatedWarehouse, err := ctrl.service.BookWarehouseForCustomer(ctx, warehouseID, customerID)
	if err != nil {
		if err.Error() == "warehouse not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if (len(err.Error()) >= 20 && err.Error()[0:20] == "invalid customer ID") ||
			(len(err.Error()) >= 20 && err.Error()[0:20] == "warehouse is already booked") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to book warehouse: %v", err.Error())})
		}
		return
	}

	c.JSON(http.StatusOK, updatedWarehouse)
}

// UnbookWarehouse handles DELETE requests to unbook a warehouse from a specific customer.
func (ctrl *WarehouseController) UnbookWarehouse(c *gin.Context) {
	warehouseID := c.Param("id")
	customerID := c.Param("customer_id")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	updatedWarehouse, err := ctrl.service.UnbookWarehouseFromCustomer(ctx, warehouseID, customerID)
	if err != nil {
		if err.Error() == "warehouse not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if len(err.Error()) >= 20 && err.Error()[0:20] == "warehouse is not booked" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to unbook warehouse: %v", err.Error())})
		}
		return
	}

	c.JSON(http.StatusOK, updatedWarehouse)
}
