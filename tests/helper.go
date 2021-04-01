package tests

import (
	"github.com/tonsV2/event-rooster-api/models"
	"github.com/tonsV2/event-rooster-api/repositories"
	"github.com/tonsV2/event-rooster-api/services"
)

var testEmail = "test@mail.com"

func getEventService() services.EventService {
	db := models.ProvideDatabase()
	eventRepository := repositories.ProvideEventRepository(db)
	eventService := services.ProvideEventService(eventRepository)
	return eventService
}

func getGroupService() services.GroupService {
	db := models.ProvideDatabase()
	groupRepository := repositories.ProvideGroupRepository(db)
	groupService := services.ProvideGroupService(groupRepository)
	return groupService
}

func getParticipantService() services.ParticipantService {
	db := models.ProvideDatabase()
	participantRepository := repositories.ProvideParticipantRepository(db)
	participantService := services.ProvideParticipantService(participantRepository)
	return participantService
}
