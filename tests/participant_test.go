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
	"time"
)

func TestAddParticipantToEventByToken(t *testing.T) {
	// TODO: This test fails because a participant with test@mail.com is already created in another test using another name
	// Either drop the database before each test... Which probably is a good idea anyway
	// Or reevaluate the whole CreateOrFind approach
	// Or , as done for now, just change the expectedName to "name0" rather than "name"

	r := gofight.New()

	server := di.BuildServer()

	eventService := getEventService()
	createdEvent, _ := eventService.Create("title", time.Now(), testEmail)

	expectedName := "name0"
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
	createdEvent, _ := eventService.Create("title", time.Now(), testEmail)

	filename := "./testdata/participants.csv"
	csvData, _ := ioutil.ReadFile(filename)

	query := gofight.H{"token": createdEvent.Token}
	uploads := []gofight.UploadFile{
		{
			Name:    "file",
			Content: csvData,
		},
	}

	r.POST("/participants/csv").
		SetQuery(query).
		SetFileFromPath(uploads).
		Run(server.Engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {

			json := r.Body.String()

			parsed := gjson.Get(json, "parsed")
			myNew := gjson.Get(json, "new")

			assert.Equal(t, 3, int(parsed.Int()))
			assert.Equal(t, 3, int(myNew.Int()))

			assert.Equal(t, http.StatusCreated, r.Code)
		})

	r.POST("/participants/csv").
		SetQuery(query).
		SetFileFromPath(uploads).
		Run(server.Engine, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {

			json := r.Body.String()

			parsed := gjson.Get(json, "parsed")
			myNew := gjson.Get(json, "new")

			assert.Equal(t, 3, int(parsed.Int()))
			assert.Equal(t, 0, int(myNew.Int()))

			assert.Equal(t, http.StatusCreated, r.Code)
		})
}

func TestAddParticipantToGroupByToken(t *testing.T) {
	r := gofight.New()

	server := di.BuildServer()

	eventService := getEventService()
	event, _ := eventService.Create("title", time.Now(), testEmail)

	groupService := getGroupService()
	group, _ := groupService.Create(event.ID, "1", time.Now(), 25)
	groupId := strconv.Itoa(int(group.ID))

	participantService := getParticipantService()
	participant, _ := participantService.CreateOrFind("name", "test@mail.com")

	_ = eventService.AddParticipantToEvent(event, participant)

	r.POST("/participants/groups").
		SetQuery(gofight.H{"groupId": groupId}).
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
