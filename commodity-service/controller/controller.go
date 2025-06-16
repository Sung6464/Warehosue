package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"commodity-service/model"   // Fixed import path
	"commodity-service/service" // Fixed import path

	"github.com/gin-gonic/gin"
)

// CommodityController handles HTTP requests for commodities.
type CommodityController struct {
	service service.CommodityService
}

// NewCommodityController creates a new instance of CommodityController.
func NewCommodityController(s service.CommodityService) *CommodityController {
	return &CommodityController{
		service: s,
	}
}

// CreateCommodity handles POST requests to create a new commodity.
func (ctrl *CommodityController) CreateCommodity(c *gin.Context) {
	var newCommodity model.Commodity
	if err := c.ShouldBindJSON(&newCommodity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err.Error())})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	createdCommodity, err := ctrl.service.CreateCommodity(ctx, newCommodity)
	if err != nil {
		if err.Error() == "commodity name is required" || err.Error() == "commodity amount must be positive" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create commodity: %v", err.Error())})
		}
		return
	}

	c.JSON(http.StatusCreated, createdCommodity)
}

// GetCommodity handles GET requests to retrieve a single commodity by ID.
func (ctrl *CommodityController) GetCommodity(c *gin.Context) {
	commodityID := c.Param("id")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	commodity, err := ctrl.service.GetCommodityByID(ctx, commodityID)
	if err != nil {
		if err.Error() == "commodity not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch commodity: %v", err.Error())})
		}
		return
	}

	c.JSON(http.StatusOK, commodity)
}

// GetAllCommodities handles GET requests to retrieve all commodities.
func (ctrl *CommodityController) GetAllCommodities(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	commodities, err := ctrl.service.GetAllCommodities(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch all commodities: %v", err.Error())})
		return
	}

	c.JSON(http.StatusOK, commodities)
}

// UpdateCommodity handles PUT requests to update an existing commodity by ID.
func (ctrl *CommodityController) UpdateCommodity(c *gin.Context) {
	commodityID := c.Param("id")

	var updatedData map[string]interface{}
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err.Error())})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	updatedCommodity, err := ctrl.service.UpdateCommodity(ctx, commodityID, updatedData)
	if err != nil {
		if err.Error() == "no fields provided for update" || err.Error() == "commodity amount cannot be negative" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if err.Error() == "commodity not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update commodity: %v", err.Error())})
		}
		return
	}

	c.JSON(http.StatusOK, updatedCommodity)
}
