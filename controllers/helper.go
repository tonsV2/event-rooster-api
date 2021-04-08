package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
)

const (
	ParseDtoFail         = "Parsing of DTO failed"
	EntityNotFound       = "Entity not found"
	UnableToCreateEntity = "Unable to create entity"
)

func handleErrorWithMessage(c *gin.Context, statusCode int, err error, message string) {
	log.Println(err)
	c.JSON(statusCode, gin.H{"error": err.Error(), "message": message})
}

func handleError(c *gin.Context, statusCode int, err error) {
	log.Println(err)
	c.JSON(statusCode, gin.H{"error": err.Error()})
}
