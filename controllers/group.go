package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/tonsV2/event-rooster-api/dtos"
	"github.com/tonsV2/event-rooster-api/services"
	"net/http"
	"strconv"
)

func ProvideGroupController(r services.EventService, g services.GroupService) GroupController {
	return GroupController{eventService: r, groupService: g}
}

type GroupController struct {
	eventService services.EventService
	groupService services.GroupService
}

func (g GroupController) GetGroupsWithParticipantsCountByEventIdAndParticipantToken(c *gin.Context) {
	eventIdStr := c.Query("eventId")
	eventId, err := strconv.ParseUint(eventIdStr, 10, 64)
	if err != nil {
		handleErrorWithMessage(c, http.StatusBadRequest, err, "Unable to parse id")
		return
	}

	participantToken := c.Query("token")

	isInEvent := g.eventService.IsParticipantInEvent(participantToken, uint(eventId))
	if !isInEvent {
		handleErrorWithMessage(c, http.StatusNotFound, err, EntityNotFound)
		return
	}

	groupsWithParticipantsCount, err := g.groupService.FindGroupsWithParticipantsCountByEventId(uint(eventId))
	if err != nil {
		handleErrorWithMessage(c, http.StatusNotFound, err, EntityNotFound)
		return
	}

	groupsDto := dtos.ToGroupsWithParticipantsCountDTO(groupsWithParticipantsCount)
	c.JSON(http.StatusOK, groupsDto)
}
