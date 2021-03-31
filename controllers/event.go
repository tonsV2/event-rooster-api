package controllers

import (
	models "github.com/tonsV2/event-rooster-api/dtos"
	"github.com/tonsV2/event-rooster-api/mail"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tonsV2/event-rooster-api/services"
)

func ProvideEventController(r services.EventService, m mail.Mailer) EventController {
	return EventController{EventService: r, Mailer: m}
}

type EventController struct {
	EventService services.EventService
	Mailer       mail.Mailer
}

func (r *EventController) CreateEvent(c *gin.Context) {
	var input models.CreateEventDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		handleError(c, err)
	}

	event, err := r.EventService.Create(input.Title, input.Date, input.Email)
	if err != nil {
		handleError(c, err)
	}

	if err := r.Mailer.SendCreateEventMail(event); err != nil {
		handleError(c, err)
	}

	eventDTO := models.ToEventDTO(event)
	c.JSON(http.StatusCreated, eventDTO)
}

func handleError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	panic(err)
}
