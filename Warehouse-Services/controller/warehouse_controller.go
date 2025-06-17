package controller

import (
	"Warehouse-Services/model"
	"Warehouse-Services/service"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// WarehouseController handles HTTP requests related to warehouses.
type WarehouseController struct {
	warehouseService service.WarehouseService
}

// NewWarehouseController creates a new instance of WarehouseController.
func NewWarehouseController(s service.WarehouseService) *WarehouseController {
	return &WarehouseController{warehouseService: s}
}

// CreateWarehouse handles POST /warehouses requests.
func (c *WarehouseController) CreateWarehouse(ctx *gin.Context) {
	var warehouse model.Warehouse
	if err := ctx.ShouldBindJSON(&warehouse); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	createdWarehouse, err := c.warehouseService.CreateWarehouse(timeoutCtx, &warehouse)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, createdWarehouse)
}

// GetAllWarehouses handles GET /warehouses requests.
func (c *WarehouseController) GetAllWarehouses(ctx *gin.Context) {
	timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	warehouses, err := c.warehouseService.GetAllWarehouses(timeoutCtx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, warehouses)
}

// GetWarehouseByID handles GET /warehouses/:id requests.
func (c *WarehouseController) GetWarehouseByID(ctx *gin.Context) {
	id := ctx.Param("id")

	timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	warehouse, err := c.warehouseService.GetWarehouseByID(timeoutCtx, id)
	if err != nil {
		if err.Error() == "warehouse not found" || err.Error() == "invalid warehouse ID format" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, warehouse)
}

// UpdateWarehouse handles PUT /warehouses/:id requests.
func (c *WarehouseController) UpdateWarehouse(ctx *gin.Context) {
	id := ctx.Param("id")
	var warehouse model.Warehouse
	if err := ctx.ShouldBindJSON(&warehouse); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	updatedWarehouse, err := c.warehouseService.UpdateWarehouse(timeoutCtx, id, &warehouse)
	if err != nil {
		if err.Error() == "warehouse not found" || err.Error() == "invalid warehouse ID format" || err.Error() == "warehouse not found or no changes made" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, updatedWarehouse)
}

// DeleteWarehouse handles DELETE /warehouses/:id requests.
func (c *WarehouseController) DeleteWarehouse(ctx *gin.Context) {
	id := ctx.Param("id")

	timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	err := c.warehouseService.DeleteWarehouse(timeoutCtx, id)
	if err != nil {
		if err.Error() == "warehouse not found" || err.Error() == "invalid warehouse ID format" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
