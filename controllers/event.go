package controllers

import (
	"github.com/tonsV2/event-rooster-api/dtos"
	"github.com/tonsV2/event-rooster-api/mail"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tonsV2/event-rooster-api/services"
)

func ProvideEventController(r services.EventService, m mail.Mailer) EventController {
	return EventController{eventService: r, mailer: m}
}

type EventController struct {
	eventService services.EventService
	mailer       mail.Mailer
}

func (e *EventController) CreateEvent(c *gin.Context) {
	var input dtos.CreateEventDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		handleError(c, err)
	}

	event, err := e.eventService.Create(input.Title, input.Date, input.Email)
	if err != nil {
		handleError(c, err)
	}

	if err := e.mailer.SendCreateEventMail(event); err != nil {
		handleError(c, err)
	}

	eventDTO := dtos.ToEventDTO(event)
	c.JSON(http.StatusCreated, eventDTO)
}

func (e *EventController) FindEventWithGroupsByToken(c *gin.Context) {
	token := c.Query("token")

	event, err := e.eventService.FindEventWithGroupsByToken(token)
	if err != nil {
		handleError(c, err)
	}

	eventDTO := dtos.ToEventWithGroupsDTO(event)
	c.JSON(http.StatusOK, eventDTO)
}

func handleError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	panic(err)
}
