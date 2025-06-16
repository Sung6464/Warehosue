package controller // The package for files in the 'controller' folder

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"Inventory-Services/model"   // FIXED: Matches go.mod module name
	"Inventory-Services/service" // FIXED: Matches go.mod module name

	"github.com/gin-gonic/gin"
)

// InventoryController handles HTTP requests for inventory.
type InventoryController struct {
	service service.InventoryService
}

// NewInventoryController creates a new instance of InventoryController.
func NewInventoryController(s service.InventoryService) *InventoryController {
	return &InventoryController{
		service: s,
	}
}

// CreateInventoryItem handles POST requests to create a new inventory item.
func (ctrl *InventoryController) CreateInventoryItem(c *gin.Context) {
	var newItem model.InventoryItem
	if err := c.ShouldBindJSON(&newItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err.Error())})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	createdItem, err := ctrl.service.CreateInventoryItem(ctx, newItem)
	if err != nil {
		if err.Error() == "warehouse ID, commodity ID, and positive quantity are required" ||
			(len(err.Error()) >= 20 && err.Error()[0:20] == "invalid warehouse ID") ||
			(len(err.Error()) >= 20 && err.Error()[0:20] == "invalid commodity ID") ||
			(len(err.Error()) >= 20 && err.Error()[0:20] == "invalid customer ID") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create inventory item: %v", err.Error())})
		}
		return
	}

	c.JSON(http.StatusCreated, createdItem)
}

// GetInventoryItem handles GET requests to retrieve a single inventory item by ID.
func (ctrl *InventoryController) GetInventoryItem(c *gin.Context) {
	itemID := c.Param("id")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	item, err := ctrl.service.GetInventoryItemByID(ctx, itemID)
	if err != nil {
		if err.Error() == "inventory item not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch inventory item: %v", err.Error())})
		}
		return
	}

	c.JSON(http.StatusOK, item)
}

// GetAllInventoryItems handles GET requests to retrieve all inventory items with filters.
func (ctrl *InventoryController) GetAllInventoryItems(c *gin.Context) {
	warehouseID := c.Query("warehouse_id")
	commodityID := c.Query("commodity_id")
	customerID := c.Query("customer_id")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	items, err := ctrl.service.GetAllInventoryItems(ctx, warehouseID, commodityID, customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch inventory items: %v", err.Error())})
		return
	}

	c.JSON(http.StatusOK, items)
}

// UpdateInventoryItem handles PUT requests to update an existing inventory item by ID.
func (ctrl *InventoryController) UpdateInventoryItem(c *gin.Context) {
	itemID := c.Param("id")

	var updatedData map[string]interface{}
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err.Error())})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	updatedItem, err := ctrl.service.UpdateInventoryItem(ctx, itemID, updatedData)
	if err != nil {
		if err.Error() == "no fields provided for update" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if err.Error() == "inventory item not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if (len(err.Error()) >= 20 && err.Error()[0:20] == "invalid warehouse ID") ||
			(len(err.Error()) >= 20 && err.Error()[0:20] == "invalid commodity ID") ||
			(len(err.Error()) >= 20 && err.Error()[0:20] == "invalid customer ID") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update inventory item: %v", err.Error())})
		}
		return
	}

	c.JSON(http.StatusOK, updatedItem)
}

// AdjustInventoryQuantity handles POST request to adjust the quantity of an inventory item.
func (ctrl *InventoryController) AdjustInventoryQuantity(c *gin.Context) {
	itemID := c.Param("id")
	var req struct {
		QuantityChange int `json:"quantity_change"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err.Error())})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	updatedItem, err := ctrl.service.AdjustInventoryQuantity(ctx, itemID, req.QuantityChange)
	if err != nil {
		if err.Error() == "inventory item not found" || err.Error() == "insufficient stock: cannot reduce quantity below zero" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to adjust inventory quantity: %v", err.Error())})
		}
		return
	}

	c.JSON(http.StatusOK, updatedItem)
}
