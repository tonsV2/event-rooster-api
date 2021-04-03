package controllers

import (
	"github.com/tonsV2/event-rooster-api/dtos"
	"github.com/tonsV2/event-rooster-api/mail"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tonsV2/event-rooster-api/services"
)

func ProvideEventController(r services.EventService, g services.GroupService, m mail.Mailer) EventController {
	return EventController{eventService: r, groupService: g, mailer: m}
}

type EventController struct {
	eventService services.EventService
	groupService services.GroupService
	mailer       mail.Mailer
}

func (e *EventController) CreateEvent(c *gin.Context) {
	var input dtos.CreateEventDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		handleErrorWithMessage(c, http.StatusBadRequest, err, ParseDtoFail)
	}

	event, err := e.eventService.Create(input.Title, input.Datetime, input.Email)
	if err != nil {
		handleErrorWithMessage(c, http.StatusBadRequest, err, UnableToCreateEntity)
	}

	if err := e.mailer.SendCreateEventMail(event); err != nil {
		handleError(c, http.StatusBadRequest, err)
	}

	eventDTO := dtos.ToEventDTO(event)
	c.JSON(http.StatusCreated, eventDTO)
}

func (e *EventController) FindEventWithGroupsByToken(c *gin.Context) {
	token := c.Query("token")

	event, err := e.eventService.FindEventWithGroupsByToken(token)
	if err != nil {
		handleErrorWithMessage(c, http.StatusNotFound, err, EntityNotFound)
	}

	eventDTO := dtos.ToEventWithGroupsDTO(event)
	c.JSON(http.StatusOK, eventDTO)
}

func (e *EventController) AddGroupToEventByToken(c *gin.Context) {
	token := c.Query("token")

	var input dtos.CreateGroupDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		handleErrorWithMessage(c, http.StatusBadRequest, err, ParseDtoFail)
	}

	event, err := e.eventService.FindByToken(token)
	if err != nil {
		handleErrorWithMessage(c, http.StatusNotFound, err, EntityNotFound)
	}

	group, err := e.groupService.Create(event.ID, input.Datetime, input.MaxParticipants)
	if err != nil {
		handleErrorWithMessage(c, http.StatusBadRequest, err, UnableToCreateEntity)
	}

	groupDTO := dtos.ToGroupDTO(group)
	c.JSON(http.StatusCreated, groupDTO)
}

func (e *EventController) GetEventWithGroupsAndParticipantsByToken(c *gin.Context) {
	token := c.Query("token")

	event, err := e.eventService.FindEventWithGroupsAndParticipantsByToken(token)
	if err != nil {
		handleErrorWithMessage(c, http.StatusNotFound, err, EntityNotFound)
	}

	eventDTO := dtos.ToEventWithGroupsAndParticipantsDTO(event)
	c.JSON(http.StatusOK, eventDTO)
}

func (e *EventController) FindEventParticipantsNotInAGroupByToken(c *gin.Context) {
	token := c.Query("token")

	participants, err := e.eventService.FindEventParticipantsNotInAGroupByToken(token)
	if err != nil {
		handleErrorWithMessage(c, http.StatusNotFound, err, EntityNotFound)
	}

	participantDTOS := dtos.ToParticipantsWithoutTokenDTO(participants)
	c.JSON(http.StatusOK, participantDTOS)
}
