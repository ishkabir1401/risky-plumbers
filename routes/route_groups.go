package routes

import (
	"github.com/gin-gonic/gin"
	"risky-plumbers/settings"
)

func SetupRoutes(router *gin.Engine) {

	riskRouteGroup := router.Group(settings.RiskRoute)

	riskRoutes(riskRouteGroup)

}
