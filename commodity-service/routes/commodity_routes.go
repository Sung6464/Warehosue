package routes

import (
	"commodity-service/controller" // Fixed import path

	"github.com/gin-gonic/gin"
)

// SetupCommodityRoutes sets up the API routes for commodity operations.
func SetupCommodityRoutes(router *gin.Engine, commodityController *controller.CommodityController) {
	commodityRoutes := router.Group("/commodities")
	{
		commodityRoutes.POST("", commodityController.CreateCommodity)
		commodityRoutes.GET("", commodityController.GetAllCommodities)
		commodityRoutes.GET("/:id", commodityController.GetCommodity)
		commodityRoutes.PUT("/:id", commodityController.UpdateCommodity)
	}
}
