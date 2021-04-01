package tests

import (
	"github.com/appleboy/gofight/v2"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"github.com/tonsV2/event-rooster-api/di"
	"github.com/tonsV2/event-rooster-api/dtos"
	"net/http"
	"testing"
)

func TestCreateEvent(t *testing.T) {
	r := gofight.New()

	server := di.BuildServer()

	expectedTitle := "title"
	expectedDate := "date"
	expectedEmail := testEmail

	eventDTO := dtos.CreateEventDTO{Title: expectedTitle, Date: expectedDate, Email: expectedEmail}

	r.POST("/events").
		SetJSONInterface(eventDTO).
		Run(server.Engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {

			json := r.Body.String()

			title := gjson.Get(json, "title")
			date := gjson.Get(json, "date")
			email := gjson.Get(json, "email")

			assert.Equal(t, expectedTitle, title.String())
			assert.Equal(t, expectedDate, date.String())
			assert.Equal(t, expectedEmail, email.String())
			assert.Equal(t, http.StatusCreated, r.Code)
		})
}

func TestFindEventWithGroupsByToken(t *testing.T) {
	r := gofight.New()

	server := di.BuildServer()

	eventService := getEventService()

	expectedTitle := "title"
	expectedDate := "date"
	expectedEmail := testEmail

	createdEvent, _ := eventService.Create(expectedTitle, expectedDate, expectedEmail)

	r.GET("/events/groups").
		SetQuery(gofight.H{"token": createdEvent.Token}).
		Run(server.Engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {

			json := r.Body.String()

			title := gjson.Get(json, "title")
			date := gjson.Get(json, "date")
			groups := gjson.Get(json, "groups")

			assert.Equal(t, expectedTitle, title.String())
			assert.Equal(t, expectedDate, date.String())
			assert.Equal(t, 0, len(groups.Array()))
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestAddGroupToEventByToken(t *testing.T) {
	r := gofight.New()

	server := di.BuildServer()
	eventService := getEventService()
	createdEvent, _ := eventService.Create("title", "date", testEmail)

	expectedDatetime := "datetime"
	expectedMaxParticipants := uint(25)
	groupDTO := dtos.CreateGroupDTO{Datetime: expectedDatetime, MaxParticipants: expectedMaxParticipants}

	r.POST("/events/groups").
		SetQuery(gofight.H{"token": createdEvent.Token}).
		SetJSONInterface(groupDTO).
		Run(server.Engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {

			json := r.Body.String()

			datetime := gjson.Get(json, "datetime")
			maxParticipants := gjson.Get(json, "maxParticipants")

			assert.Equal(t, expectedDatetime, datetime.String())
			assert.Equal(t, expectedMaxParticipants, uint(maxParticipants.Uint()))
			assert.Equal(t, http.StatusCreated, r.Code)
		})
}

func TestGetEventWithGroupsAndParticipantsByToken(t *testing.T) {
	r := gofight.New()

	server := di.BuildServer()

	eventService := getEventService()
	event, _ := eventService.Create("title", "date", testEmail)

	groupService := getGroupService()
	group0, _ := groupService.Create(event.ID, "datetime0", 25)
	group1, _ := groupService.Create(event.ID, "datetime1", 25)

	participantService := getParticipantService()
	participant0, _ := participantService.CreateOrFind("name0", "test@mail.com")
	participant1, _ := participantService.CreateOrFind("name1", "test1@mail.com")
	participant2, _ := participantService.CreateOrFind("name2", "test2@mail.com")

	_ = eventService.AddParticipantToEvent(event, participant0)
	_ = eventService.AddParticipantToEvent(event, participant1)
	_ = eventService.AddParticipantToEvent(event, participant2)

	_ = groupService.AddParticipant(group0, participant0)
	_ = groupService.AddParticipant(group0, participant1)
	_ = groupService.AddParticipant(group1, participant2)

	r.GET("/events").
		SetQuery(gofight.H{"token": event.Token}).
		Run(server.Engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {

			json := r.Body.String()

			date := gjson.Get(json, "date")
			assert.Equal(t, event.Date, date.String())

			title := gjson.Get(json, "title")
			assert.Equal(t, event.Title, title.String())

			group0Datetime := gjson.Get(json, "groups.0.datetime")
			assert.Equal(t, group0.Datetime, group0Datetime.String())

			group0MaxParticipants := gjson.Get(json, "groups.0.maxParticipants")
			assert.Equal(t, group0.MaxParticipants, uint(group0MaxParticipants.Uint()))

			group0Id := gjson.Get(json, "groups.0.id")
			assert.Equal(t, group0.ID, uint(group0Id.Uint()))

			group0Participants0Email := gjson.Get(json, "groups.0.participants.0.email")
			assert.Equal(t, participant0.Email, group0Participants0Email.String())

			groupsLength := gjson.Get(json, "groups.#")
			assert.Equal(t, 2, int(groupsLength.Int()))

			group0ParticipantsLength := gjson.Get(json, "groups.0.participants.#")
			assert.Equal(t, 2, int(group0ParticipantsLength.Int()))

			group1ParticipantsLength := gjson.Get(json, "groups.1.participants.#")
			assert.Equal(t, 1, int(group1ParticipantsLength.Int()))

			assert.Equal(t, http.StatusOK, r.Code)
		})
}
