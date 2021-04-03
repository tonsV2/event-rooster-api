package controllers

import (
	"github.com/gin-gonic/gin"
)

const (
	ParseDtoFail         = "Parsing of DTO failed"
	EntityNotFound       = "Entity not found"
	UnableToCreateEntity = "Unable to create entity"
)

func handleErrorWithMessage(c *gin.Context, statusCode int, err error, message string) {
	c.JSON(statusCode, gin.H{"error": err.Error(), "message": message})
	panic(err)
}

func handleError(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, gin.H{"error": err.Error()})
	panic(err)
}
