package routes

import (
	"commodity-service/controller"
	"commodity-service/service"
	"fmt"      // Import fmt for string formatting
	"net/http" // Import http for redirects

	"github.com/gin-gonic/gin"
)

// CommodityRoutes sets up the API routes for commodity operations.
func CommodityRoutes(router *gin.Engine) {
	commodityController := controller.NewCommodityController(service.NewCommodityService())

	// Primary routes: define WITHOUT a trailing slash for collection endpoints
	commodityGroup := router.Group("/commodities")
	{
		commodityGroup.POST("", commodityController.CreateCommodity)     // Matches /commodities
		commodityGroup.GET("", commodityController.GetAllCommodities)    // Matches /commodities
		commodityGroup.GET("/:id", commodityController.GetCommodityByID) // Matches /commodities/:id
		commodityGroup.PUT("/:id", commodityController.UpdateCommodity)
		commodityGroup.DELETE("/:id", commodityController.DeleteCommodity)
	}

	// Add explicit 301 redirects for paths that might come in WITH trailing slashes.
	router.GET("/commodities/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/commodities")
	})
	router.POST("/commodities/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/commodities")
	})
	router.PUT("/commodities/:id/", func(c *gin.Context) {
		id := c.Param("id")
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/commodities/%s", id))
	})
	router.DELETE("/commodities/:id/", func(c *gin.Context) {
		id := c.Param("id")
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/commodities/%s", id))
	})
	router.GET("/commodities/:id/", func(c *gin.Context) {
		id := c.Param("id")
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/commodities/%s", id))
	})
	// Add other methods if necessary
}
