package controllers

import (
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nu7hatch/gouuid"
	"github.com/tonsV2/event-rooster-api/dtos"
	"github.com/tonsV2/event-rooster-api/mail"
	"github.com/tonsV2/event-rooster-api/models"
	"github.com/tonsV2/event-rooster-api/services"
	"io"
	"log"
	"net/http"
	"os"
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

func (p *ParticipantController) AddParticipantsCSVToEventByToken(c *gin.Context) {
	token := c.Query("token")

	event, err := p.eventService.FindByToken(token)
	if err != nil {
		handleErrorWithMessage(c, http.StatusNotFound, err, EntityNotFound)
	}

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

	count := 0
	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		name := record[0]
		email := record[1]

		p.addParticipantToEvent(c, event, name, email)
		count++
	}

	if err := os.Remove(newFileName); err != nil {
		handleError(c, http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusCreated, fmt.Sprintf("%d participants parsed", count))
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

	// Confirm participant is associated with event
	event, err := p.eventService.FindByIdAndParticipantToken(group.EventID, token)
	if err != nil {
		handleErrorWithMessage(c, http.StatusNotFound, err, EntityNotFound)
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
