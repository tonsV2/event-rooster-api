package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ProvideHealthController() HealthController {
	return HealthController{}
}

type HealthController struct {
}

func (e *HealthController) GetHealthStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "UP"})
}
