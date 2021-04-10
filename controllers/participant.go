package controllers

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nu7hatch/gouuid"
	"github.com/tonsV2/event-rooster-api/dtos"
	"github.com/tonsV2/event-rooster-api/mail"
	"github.com/tonsV2/event-rooster-api/models"
	"github.com/tonsV2/event-rooster-api/services"
	"net/http"
	"os"
	"strings"
)

func ProvideParticipantController(r services.EventService, p services.ParticipantService, m mail.Mailer, g services.GroupService) ParticipantController {
	return ParticipantController{eventService: r, participantService: p, mailer: m, groupService: g}
}

type ParticipantController struct {
	eventService       services.EventService
	participantService services.ParticipantService
	groupService       services.GroupService
	mailer             mail.Mailer
}

func (p *ParticipantController) AddParticipantToEventByToken(c *gin.Context) {
	token := c.Query("token")

	var input dtos.CreateParticipantDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		handleErrorWithMessage(c, http.StatusBadRequest, err, ParseDtoFail)
	}

	event, err := p.eventService.FindByToken(token)
	if err != nil {
		handleErrorWithMessage(c, http.StatusNotFound, err, EntityNotFound)
	}

	participant := p.addParticipantToEvent(c, event, input.Name, input.Email)

	participantDTO := dtos.ToParticipantDTO(participant)
	c.JSON(http.StatusCreated, participantDTO)
}

func (p *ParticipantController) addParticipantToEvent(c *gin.Context, event models.Event, name string, email string) models.Participant {
	participant, err := p.participantService.CreateOrFind(name, email)
	if err != nil {
		handleError(c, http.StatusBadRequest, err)
	}

	err = p.eventService.AddParticipantToEvent(event, participant)
	if err != nil {
		handleError(c, http.StatusBadRequest, err)
	}

	// TODO: Don't send mail if participant already added to event
	if err := p.mailer.SendWelcomeParticipantMail(event, participant); err != nil {
		handleError(c, http.StatusBadRequest, err)
	}

	return participant
}

func (p *ParticipantController) AddParticipantsCSVToEventByToken(c *gin.Context) {
	token := c.Query("token")

	event, err := p.eventService.FindByToken(token)
	if err != nil {
		handleErrorWithMessage(c, http.StatusNotFound, err, EntityNotFound)
	}

	records := p.parseCSV(c, err)

	participants := p.addParticipantsToEvent(c, event, records)

	if err := p.mailer.SendWelcomeParticipantMails(event, participants); err != nil {
		handleError(c, http.StatusBadRequest, err)
	}

	body := gin.H{"parsed": len(records), "new": len(participants)}
	c.JSON(http.StatusCreated, body)
}

func (p *ParticipantController) parseCSV(c *gin.Context, err error) [][]string {
	file, err := c.FormFile("file")
	if err != nil {
		handleError(c, http.StatusInternalServerError, err)
	}

	u4, err := uuid.NewV4()
	if err != nil {
		handleError(c, http.StatusInternalServerError, err)
	}
	newFileName := "/tmp/" + fmt.Sprint(u4)

	if err := c.SaveUploadedFile(file, newFileName); err != nil {
		handleError(c, http.StatusInternalServerError, err)
	}

	in, err := os.Open(newFileName)
	if err != nil {
		handleError(c, http.StatusInternalServerError, err)
	}

	r := csv.NewReader(in)
	records, err := r.ReadAll()

	if err := os.Remove(newFileName); err != nil {
		handleError(c, http.StatusInternalServerError, err)
	}
	return records
}

func (p *ParticipantController) addParticipantsToEvent(c *gin.Context, event models.Event, records [][]string) []models.Participant {
	var participants []models.Participant

	for _, record := range records {
		name := record[0]
		email := record[1]

		participant, err := p.participantService.CreateOrFind(name, email)
		if err != nil {
			handleError(c, http.StatusBadRequest, err)
		}

		isInEvent := p.eventService.IsParticipantInEvent(participant.Token, event.ID)
		if isInEvent {
			continue
		}

		err = p.eventService.AddParticipantToEvent(event, participant)
		if err != nil {
			handleError(c, http.StatusBadRequest, err)
		}

		participants = append(participants, participant)
	}

	return participants
}

func (p *ParticipantController) AddParticipantToGroupByToken(c *gin.Context) {
	token := c.Query("token")
	groupId := c.Query("groupId")

	participant, err := p.participantService.FindByToken(token)
	if err != nil {
		handleErrorWithMessage(c, http.StatusNotFound, err, EntityNotFound)
	}

	group, err := p.groupService.FindById(groupId)
	if err != nil {
		handleErrorWithMessage(c, http.StatusNotFound, err, EntityNotFound)
	}

	isInEvent := p.eventService.IsParticipantInEvent(token, group.EventID)
	if !isInEvent {
		handleError(c, http.StatusNotFound, errors.New(strings.ToLower(EntityNotFound)))
		return
	}

	event, err := p.eventService.FindById(group.EventID)
	if err != nil {
		handleError(c, http.StatusNotFound, err)
	}

	if p.isGroupFull(event, group.ID) {
		handleError(c, http.StatusForbidden, errors.New("group is full"))
	}

	err = p.groupService.AddParticipant(group, participant)
	if err != nil {
		handleError(c, http.StatusBadRequest, err)
	}

	if err := p.mailer.SendWelcomeToGroupMail(event, group, participant); err != nil {
		handleError(c, http.StatusBadRequest, err)
	}

	participantDTO := dtos.ToParticipantDTO(participant)
	c.JSON(http.StatusCreated, participantDTO)
}

func (p *ParticipantController) isGroupFull(event models.Event, groupId uint) bool {
	groupParticipantsCounts, err := p.groupService.FindGroupsWithParticipantsCountByEventId(event.ID)
	if err != nil {
		panic(err.Error())
	}
	for _, group := range groupParticipantsCounts {
		if group.ID == groupId {
			maxParticipants := group.MaxParticipants
			actualParticipants := group.ActualParticipants
			return !(maxParticipants > actualParticipants)
		}
	}
	panic(errors.New("no matching group found"))
}
