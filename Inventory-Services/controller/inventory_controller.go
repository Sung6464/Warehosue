package controller

import (
	"Inventory-Services/model"
	"Inventory-Services/service"
	"context" // Added context import
	"net/http"
	"time" // Added time import

	"github.com/gin-gonic/gin"
)

// InventoryController handles HTTP requests related to inventory.
type InventoryController struct {
	inventoryService service.InventoryService
}

// NewInventoryController creates a new instance of InventoryController.
func NewInventoryController(s service.InventoryService) *InventoryController {
	return &InventoryController{inventoryService: s}
}

// CreateInventory handles POST /inventory requests.
func (c *InventoryController) CreateInventory(ctx *gin.Context) {
	var inventory model.Inventory
	if err := ctx.ShouldBindJSON(&inventory); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	createdInventory, err := c.inventoryService.CreateInventory(timeoutCtx, &inventory)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, createdInventory)
}

// GetAllInventories handles GET /inventory requests.
func (c *InventoryController) GetAllInventories(ctx *gin.Context) {
	timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	inventories, err := c.inventoryService.GetAllInventories(timeoutCtx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, inventories)
}

// GetInventoryByID handles GET /inventory/:id requests.
func (c *InventoryController) GetInventoryByID(ctx *gin.Context) {
	id := ctx.Param("id")

	timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	inventory, err := c.inventoryService.GetInventoryByID(timeoutCtx, id)
	if err != nil {
		if err.Error() == "inventory not found" || err.Error() == "invalid inventory ID format" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, inventory)
}

// UpdateInventory handles PUT /inventory/:id requests.
func (c *InventoryController) UpdateInventory(ctx *gin.Context) {
	id := ctx.Param("id")
	var inventory model.Inventory
	if err := ctx.ShouldBindJSON(&inventory); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	updatedInventory, err := c.inventoryService.UpdateInventory(timeoutCtx, id, &inventory)
	if err != nil {
		if err.Error() == "inventory not found" || err.Error() == "invalid inventory ID format" || err.Error() == "inventory not found or no changes made" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, updatedInventory)
}

// DeleteInventory handles DELETE /inventory/:id requests.
func (c *InventoryController) DeleteInventory(ctx *gin.Context) {
	id := ctx.Param("id")

	timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	err := c.inventoryService.DeleteInventory(timeoutCtx, id)
	if err != nil {
		if err.Error() == "inventory not found" || err.Error() == "invalid inventory ID format" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
