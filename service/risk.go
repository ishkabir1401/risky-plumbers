package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"risky-plumbers/model"
	"risky-plumbers/settings"
	"risky-plumbers/utils"
)

func ListRisks(c *gin.Context) {
	settings.RiskStoreStruct.Mux.Lock()
	defer settings.RiskStoreStruct.Mux.Unlock()

	var risks = make([]model.Risk, 0)
	for _, risk := range settings.RiskStoreStruct.Risks {
		risks = append(risks, *risk)
	}
	c.JSON(http.StatusOK, gin.H{"data": risks})
}

func CreateRisk(c *gin.Context) {
	var req model.Risk
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if req.State == "" || req.Title == "" || req.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}
	if req.State != "open" && req.State != "closed" && req.State != "accepted" && req.State != "investigating" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state value"})
		return
	}
	req.ID = utils.GenerateUUid()
	settings.RiskStoreStruct.Mux.Lock()
	settings.RiskStoreStruct.Risks[req.ID] = &req
	defer settings.RiskStoreStruct.Mux.Unlock()

	c.JSON(http.StatusOK, gin.H{"data": req})
}

func GetRisk(c *gin.Context) {
	id := c.Param("id")
	settings.RiskStoreStruct.Mux.Lock()
	risk, exists := settings.RiskStoreStruct.Risks[id]
	defer settings.RiskStoreStruct.Mux.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Risk not found"})
		return
	}
	c.JSON(http.StatusOK, risk)
}
