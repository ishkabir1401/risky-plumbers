package routes

import (
	"github.com/gin-gonic/gin"
	"risky-plumbers/service"
)

func riskRoutes(group *gin.RouterGroup) {
	group.GET("/risks", service.ListRisks)
	group.POST("/risks", service.CreateRisk)
	group.GET("/risks/:id", service.GetRisk)
}
