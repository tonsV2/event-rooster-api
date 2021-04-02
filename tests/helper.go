package tests

import (
	"github.com/tidwall/gjson"
	"github.com/tonsV2/event-rooster-api/models"
	"github.com/tonsV2/event-rooster-api/repositories"
	"github.com/tonsV2/event-rooster-api/services"
	"time"
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

func CompareTimeToResult(mytime time.Time, result gjson.Result) bool {
	parse, _ := time.Parse(time.RFC3339, result.String())
	return mytime.Format(time.RFC3339) == parse.Format(time.RFC3339)
}
