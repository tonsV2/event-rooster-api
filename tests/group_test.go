package tests

import (
	"github.com/appleboy/gofight/v2"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"github.com/tonsV2/event-rooster-api/di"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestGetGroupsWithParticipantsCountByEventIdAndParticipantToken(t *testing.T) {
	r := gofight.New()

	server := di.BuildServer()

	eventService := getEventService()
	event, _ := eventService.Create("title", time.Now(), testEmail)

	groupService := getGroupService()
	expectedGid := "1"
	expectedDatetime := time.Now()
	expectedActualMaxParticipants := uint(25)

	group, _ := groupService.Create(event.ID, expectedGid, expectedDatetime, expectedActualMaxParticipants)

	participantService := getParticipantService()
	participant, _ := participantService.CreateOrFind("name", testEmail)

	_ = eventService.AddParticipantToEvent(event, participant)

	r.GET("/events/groups-count").
		SetQuery(gofight.H{
			"eventId": strconv.Itoa(int(event.ID)),
			"token":   participant.Token},
		).
		Run(server.Engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			json := r.Body.String()

			groups := gjson.Get(json, "#")
			assert.Equal(t, uint(1), uint(groups.Uint()))

			id := gjson.Get(json, "0.id")
			assert.Equal(t, group.ID, uint(id.Uint()))

			maxParticipants := gjson.Get(json, "0.maxParticipants")
			assert.Equal(t, expectedActualMaxParticipants, uint(maxParticipants.Uint()))

			gid := gjson.Get(json, "0.gid")
			assert.Equal(t, expectedGid, gid.String())

			actualParticipants := gjson.Get(json, "0.actualParticipants")
			assert.Equal(t, uint(0), uint(actualParticipants.Uint()))

			assert.Equal(t, http.StatusOK, r.Code)
		})
}
