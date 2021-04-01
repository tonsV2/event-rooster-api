package tests

import (
	"github.com/appleboy/gofight/v2"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"github.com/tonsV2/event-rooster-api/di"
	"github.com/tonsV2/event-rooster-api/dtos"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestAddParticipantToEventByToken(t *testing.T) {
	r := gofight.New()

	server := di.BuildServer()
	eventService := getEventService()
	createdEvent, _ := eventService.Create("title", "date", testEmail)

	expectedName := "name"
	expectedEmail := "test@mail.com"
	participantDTO := dtos.CreateParticipantDTO{Name: expectedName, Email: expectedEmail}

	r.POST("/participants").
		SetQuery(gofight.H{"token": createdEvent.Token}).
		SetJSONInterface(participantDTO).
		Run(server.Engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {

			json := r.Body.String()

			name := gjson.Get(json, "name")
			email := gjson.Get(json, "email")

			assert.Equal(t, expectedName, name.String())
			assert.Equal(t, expectedEmail, email.String())
			assert.Equal(t, http.StatusCreated, r.Code)
		})
}

func TestAddParticipantsCSVToEventByToken(t *testing.T) {
	r := gofight.New()

	server := di.BuildServer()

	eventService := getEventService()
	createdEvent, _ := eventService.Create("title", "date", testEmail)

	filename := "./testdata/participants.csv"
	csvData, _ := ioutil.ReadFile(filename)

	r.POST("/participants/csv").
		SetQuery(gofight.H{"token": createdEvent.Token}).
		SetFileFromPath([]gofight.UploadFile{
			{
				Name:    "file",
				Content: csvData,
			},
		}).
		Run(server.Engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {

			body := r.Body.String()

			assert.Equal(t, "\"2 participants parsed\"", body)
			assert.Equal(t, http.StatusCreated, r.Code)
		})
}
