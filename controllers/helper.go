package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func handleError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	panic(err)
}
