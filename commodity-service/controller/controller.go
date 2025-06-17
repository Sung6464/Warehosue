package controller

import (
	"commodity-service/model"
	"commodity-service/service"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CommodityController handles HTTP requests related to commodities.
type CommodityController struct {
	commodityService service.CommodityService
}

// NewCommodityController creates a new instance of CommodityController.
func NewCommodityController(s service.CommodityService) *CommodityController {
	return &CommodityController{commodityService: s}
}

// CreateCommodity handles POST /commodities requests.
func (c *CommodityController) CreateCommodity(ctx *gin.Context) {
	var commodity model.Commodity
	if err := ctx.ShouldBindJSON(&commodity); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	createdCommodity, err := c.commodityService.CreateCommodity(timeoutCtx, &commodity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, createdCommodity)
}

// GetAllCommodities handles GET /commodities requests.
func (c *CommodityController) GetAllCommodities(ctx *gin.Context) {
	timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	commodities, err := c.commodityService.GetAllCommodities(timeoutCtx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, commodities)
}

// GetCommodityByID handles GET /commodities/:id requests.
func (c *CommodityController) GetCommodityByID(ctx *gin.Context) {
	id := ctx.Param("id")

	timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	commodity, err := c.commodityService.GetCommodityByID(timeoutCtx, id)
	if err != nil {
		if err.Error() == "commodity not found" || err.Error() == "invalid commodity ID format" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, commodity)
}

// UpdateCommodity handles PUT /commodities/:id requests.
func (c *CommodityController) UpdateCommodity(ctx *gin.Context) {
	id := ctx.Param("id")
	var commodity model.Commodity
	if err := ctx.ShouldBindJSON(&commodity); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	updatedCommodity, err := c.commodityService.UpdateCommodity(timeoutCtx, id, &commodity)
	if err != nil {
		if err.Error() == "commodity not found" || err.Error() == "invalid commodity ID format" || err.Error() == "commodity not found or no changes made" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, updatedCommodity)
}

// DeleteCommodity handles DELETE /commodities/:id requests.
func (c *CommodityController) DeleteCommodity(ctx *gin.Context) {
	id := ctx.Param("id")

	timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	err := c.commodityService.DeleteCommodity(timeoutCtx, id)
	if err != nil {
		if err.Error() == "commodity not found" || err.Error() == "invalid commodity ID format" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
