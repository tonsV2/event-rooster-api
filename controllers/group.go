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
	eventIdStr := c.Query("id")
	eventId, err := strconv.ParseUint(eventIdStr, 10, 64)
	if err != nil {
		handleError(c, err) // TODO: id not valid
	}

	participantToken := c.Query("token")

	event, err := g.eventService.FindByIdAndParticipantToken(uint(eventId), participantToken)
	if err != nil {
		handleError(c, err) // TODO: 404
	}

	groupsWithParticipantsCount, err := g.groupService.FindGroupsWithParticipantsCountByEventId(event.ID)
	if err != nil {
		handleError(c, err) // TODO: 404
	}

	groupsDto := dtos.ToGroupsWithParticipantsCountDTO(groupsWithParticipantsCount)
	c.JSON(http.StatusOK, groupsDto)
}
