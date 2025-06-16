package routes

import (
	"api-gateway/controller" // FIXED: Corrected import path

	"github.com/gin-gonic/gin"
)

// SetupGatewayRoutes sets up the API routes for the API Gateway.
func SetupGatewayRoutes(router *gin.Engine, gatewayController *controller.GatewayController) {
	apiGroup := router.Group("/api")
	{
		apiGroup.Any("/customers/*proxyPath", gatewayController.ProxyToCustomerService)
		apiGroup.Any("/warehouses/*proxyPath", gatewayController.ProxyToWarehouseService)
		apiGroup.Any("/commodities/*proxyPath", gatewayController.ProxyToCommoditiesService)
		apiGroup.Any("/inventory/*proxyPath", gatewayController.ProxyToInventoryService)
	}
}
