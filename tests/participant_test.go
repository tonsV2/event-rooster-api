package tests

import (
	"github.com/appleboy/gofight/v2"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"github.com/tonsV2/event-rooster-api/di"
	"github.com/tonsV2/event-rooster-api/dtos"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

func TestAddParticipantToEventByToken(t *testing.T) {
	r := gofight.New()

	server := di.BuildServer()

	eventService := getEventService()
	createdEvent, _ := eventService.Create("title", "datetime", testEmail)

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
	createdEvent, _ := eventService.Create("title", "datetime", testEmail)

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

			assert.Equal(t, "\"3 participants parsed\"", body)
			assert.Equal(t, http.StatusCreated, r.Code)
		})
}

func TestAddParticipantToGroupByToken(t *testing.T) {
	r := gofight.New()

	server := di.BuildServer()

	eventService := getEventService()
	event, _ := eventService.Create("title", "datetime", testEmail)

	groupService := getGroupService()
	group, _ := groupService.Create(event.ID, "datetime", 25)
	groupId := strconv.Itoa(int(group.ID))

	participantService := getParticipantService()
	participant, _ := participantService.CreateOrFind("name", "test@mail.com")

	_ = eventService.AddParticipantToEvent(event, participant)

	r.POST("/participants/groups").
		SetQuery(gofight.H{"id": groupId}).
		SetQuery(gofight.H{"token": participant.Token}).
		Run(server.Engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {

			json := r.Body.String()

			id := gjson.Get(json, "id")
			name := gjson.Get(json, "name")
			email := gjson.Get(json, "email")
			token := gjson.Get(json, "token")

			assert.Equal(t, participant.ID, uint(id.Uint()))
			assert.Equal(t, participant.Name, name.String())
			assert.Equal(t, participant.Email, email.String())
			assert.Equal(t, participant.Token, token.String())
			assert.Equal(t, http.StatusCreated, r.Code)
		})
}
