package tests

import (
	"github.com/appleboy/gofight/v2"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"github.com/tonsV2/event-rooster-api/di"
	"github.com/tonsV2/event-rooster-api/dtos"
	"net/http"
	"testing"
	"time"
)

func TestCreateEvent(t *testing.T) {
	r := gofight.New()

	server := di.BuildServer()

	expectedTitle := "title"
	expectedDatetime := time.Now()
	expectedEmail := testEmail

	eventDTO := dtos.CreateEventDTO{Title: expectedTitle, Datetime: expectedDatetime, Email: expectedEmail}

	r.POST("/events").
		SetJSONInterface(eventDTO).
		Run(server.Engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {

			json := r.Body.String()

			title := gjson.Get(json, "title")
			datetime := gjson.Get(json, "datetime")
			email := gjson.Get(json, "email")

			assert.Equal(t, expectedTitle, title.String())
			assert.True(t, CompareTimeToResult(expectedDatetime, datetime))
			assert.Equal(t, expectedEmail, email.String())
			assert.Equal(t, http.StatusCreated, r.Code)
		})
}

func TestFindEventWithGroupsByToken(t *testing.T) {
	r := gofight.New()

	server := di.BuildServer()

	eventService := getEventService()

	expectedTitle := "title"
	expectedDatetime := time.Now()
	expectedEmail := testEmail

	createdEvent, _ := eventService.Create(expectedTitle, expectedDatetime, expectedEmail)

	r.GET("/events/groups").
		SetQuery(gofight.H{"token": createdEvent.Token}).
		Run(server.Engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {

			json := r.Body.String()

			title := gjson.Get(json, "title")
			datetime := gjson.Get(json, "datetime")
			groups := gjson.Get(json, "groups")

			assert.Equal(t, expectedTitle, title.String())
			assert.True(t, CompareTimeToResult(expectedDatetime, datetime))
			assert.Equal(t, 0, len(groups.Array()))
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestAddGroupToEventByToken(t *testing.T) {
	r := gofight.New()

	server := di.BuildServer()
	eventService := getEventService()
	createdEvent, _ := eventService.Create("title", time.Now(), testEmail)

	expectedGid := "1"
	expectedDatetime := time.Now()
	expectedMaxParticipants := uint(25)
	groupDTO := dtos.CreateGroupDTO{GID: expectedGid, Datetime: expectedDatetime, MaxParticipants: expectedMaxParticipants}

	r.POST("/events/groups").
		SetQuery(gofight.H{"token": createdEvent.Token}).
		SetJSONInterface(groupDTO).
		Run(server.Engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {

			json := r.Body.String()

			gid := gjson.Get(json, "gid")
			datetime := gjson.Get(json, "datetime")
			maxParticipants := gjson.Get(json, "maxParticipants")

			assert.Equal(t, expectedGid, gid.String())
			assert.True(t, CompareTimeToResult(expectedDatetime, datetime))
			assert.Equal(t, expectedMaxParticipants, uint(maxParticipants.Uint()))
			assert.Equal(t, http.StatusCreated, r.Code)
		})
}

func TestGetEventWithGroupsAndParticipantsByToken(t *testing.T) {
	r := gofight.New()

	server := di.BuildServer()

	eventService := getEventService()
	event, _ := eventService.Create("title", time.Now(), testEmail)

	groupService := getGroupService()
	group0, _ := groupService.Create(event.ID, "1", time.Now(), 25)
	group1, _ := groupService.Create(event.ID, "2", time.Now(), 25)

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

			datetime := gjson.Get(json, "datetime")
			assert.True(t, CompareTimeToResult(event.Datetime, datetime))

			title := gjson.Get(json, "title")
			assert.Equal(t, event.Title, title.String())

			group0Gid := gjson.Get(json, "groups.0.gid")
			assert.Equal(t, group0.GID, group0Gid.String())

			group0Datetime := gjson.Get(json, "groups.0.datetime")
			assert.True(t, CompareTimeToResult(group0.Datetime, group0Datetime))

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

func TestFindEventParticipantsNotInAGroupByToken(t *testing.T) {
	r := gofight.New()

	server := di.BuildServer()

	eventService := getEventService()
	event, _ := eventService.Create("title", time.Now(), testEmail)

	groupService := getGroupService()
	group0, _ := groupService.Create(event.ID, "1", time.Now(), 25)

	participantService := getParticipantService()
	participant0, _ := participantService.CreateOrFind("name0", "test@mail.com")
	participant1, _ := participantService.CreateOrFind("name1", "test1@mail.com")

	_ = eventService.AddParticipantToEvent(event, participant0)
	_ = eventService.AddParticipantToEvent(event, participant1)

	_ = groupService.AddParticipant(group0, participant0)

	r.GET("/participants/not-in-groups").
		SetQuery(gofight.H{"token": event.Token}).
		Run(server.Engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {

			json := r.Body.String()

			groupsLength := gjson.Get(json, "#")
			assert.Equal(t, 1, int(groupsLength.Int()))

			name := gjson.Get(json, "0.name")
			assert.Equal(t, participant1.Name, name.String())

			id := gjson.Get(json, "0.id")
			assert.Equal(t, participant1.ID, uint(id.Uint()))

			email := gjson.Get(json, "0.email")
			assert.Equal(t, participant1.Email, email.String())

			assert.Equal(t, http.StatusOK, r.Code)
		})
}
