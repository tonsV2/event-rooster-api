package controllers

import (
	"github.com/tonsV2/race-rooster-api/mail"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tonsV2/race-rooster-api/services"
)

func ProvideRaceController(r services.RaceService, m mail.Mailer) RaceController {
	return RaceController{RaceService: r, Mailer: m}
}

type RaceController struct {
	RaceService services.RaceService
	Mailer      mail.Mailer
}

type CreateRaceInput struct {
	Title string `json:"title" binding:"required"`
	Date  string `json:"date" binding:"required"`
	Email string `json:"email" binding:"required"`
}

func (r *RaceController) CreateRace(c *gin.Context) {
	var input CreateRaceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		handleError(c, err)
	}

	race, err := r.RaceService.Create(input.Title, input.Date, input.Email)
	if err != nil {
		handleError(c, err)
	}

	if err := r.Mailer.SendCreateRaceMail(race); err != nil {
		handleError(c, err)
	}

	c.JSON(http.StatusCreated, race)
}

func handleError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	panic(err)
}
