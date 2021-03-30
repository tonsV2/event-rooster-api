package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	. "github.com/tonsV2/race-rooster-api/services"
)

func ProvideRaceController(r RaceService) RaceController {
	return RaceController{RaceService: r}
}

type RaceController struct {
	RaceService RaceService
}

type CreateRaceInput struct {
	Title string `json:"title" binding:"required"`
	Date  string `json:"date" binding:"required"`
	Email string `json:"email" binding:"required"`
}

func (r *RaceController) CreateRace(c *gin.Context) {
	var input CreateRaceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	race, err := r.RaceService.Create(input.Title, input.Date, input.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, race)
}
