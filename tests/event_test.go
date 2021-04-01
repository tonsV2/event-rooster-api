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

			data := []byte(r.Body.String())

			title := gjson.GetBytes(data, "title")
			date := gjson.GetBytes(data, "date")
			email := gjson.GetBytes(data, "email")

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

			data := []byte(r.Body.String())

			title := gjson.GetBytes(data, "title")
			date := gjson.GetBytes(data, "date")
			groups := gjson.GetBytes(data, "groups")

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

			data := []byte(r.Body.String())

			datetime := gjson.GetBytes(data, "datetime")
			maxParticipants := gjson.GetBytes(data, "maxParticipants")

			assert.Equal(t, expectedDatetime, datetime.String())
			assert.Equal(t, expectedMaxParticipants, uint(maxParticipants.Uint()))
			assert.Equal(t, http.StatusCreated, r.Code)
		})
}
