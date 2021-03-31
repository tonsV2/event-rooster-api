package controllers

import (
	"github.com/tonsV2/event-rooster-api/dtos"
	"github.com/tonsV2/event-rooster-api/mail"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tonsV2/event-rooster-api/services"
)

func ProvideParticipantController(r services.EventService, p services.ParticipantService, m mail.Mailer) ParticipantController {
	return ParticipantController{eventService: r, participantService: p, mailer: m}
}

type ParticipantController struct {
	eventService       services.EventService
	participantService services.ParticipantService
	mailer             mail.Mailer
}

func (p *ParticipantController) AddParticipantToEventByToken(c *gin.Context) {
	token := c.Query("token")

	var input dtos.CreateParticipantDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		handleError(c, err)
	}

	event, err := p.eventService.FindByToken(token)
	if err != nil {
		handleError(c, err) // TODO: 404
	}

	participant, err := p.participantService.CreateOrFind(input.Name, input.Email)
	if err != nil {
		handleError(c, err) // TODO: idk?
	}

	err = p.eventService.AddParticipantToEvent(event, participant)
	if err != nil {
		handleError(c, err) // TODO: idk?
	}

	if err := p.mailer.SendWelcomeParticipantMail(event, participant); err != nil {
		handleError(c, err) // TODO: Mail error...
	}

	participantDTO := dtos.ToParticipantDTO(participant)
	c.JSON(http.StatusCreated, participantDTO)
}
