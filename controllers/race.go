package controllers

import (
	models "github.com/tonsV2/race-rooster-api/dtos"
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

func (r *RaceController) CreateRace(c *gin.Context) {
	var input models.CreateRaceDTO
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

	raceDTO := models.ToRaceDTO(race)
	c.JSON(http.StatusCreated, raceDTO)
}

func handleError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	panic(err)
}
