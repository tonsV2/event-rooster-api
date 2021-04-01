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

func ProvideParticipantController(r services.EventService, p services.ParticipantService, m mail.Mailer) ParticipantController {
	return ParticipantController{eventService: r, participantService: p, mailer: m}
}

type ParticipantController struct {
	eventService       services.EventService
	participantService services.ParticipantService
	mailer             mail.Mailer
}

func (p *ParticipantController) AddParticipantToEventByToken(c *gin.Context) {
	token := c.Query("token")

	var input dtos.CreateParticipantDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		handleError(c, err)
	}

	event, err := p.eventService.FindByToken(token)
	if err != nil {
		handleError(c, err) // TODO: 404
	}

	participant := p.addParticipantToEvent(c, event, input.Name, input.Email)

	participantDTO := dtos.ToParticipantDTO(participant)
	c.JSON(http.StatusCreated, participantDTO)
}

func (p *ParticipantController) AddParticipantsCSVToEventByToken(c *gin.Context) {
	token := c.Query("token")

	event, err := p.eventService.FindByToken(token)
	if err != nil {
		handleError(c, err) // TODO: 404
	}

	file, err := c.FormFile("file")
	if err != nil {
		handleError(c, err)
	}

	u4, err := uuid.NewV4()
	if err != nil {
		handleError(c, err)
	}
	newFileName := "/tmp/" + fmt.Sprint(u4)

	if err := c.SaveUploadedFile(file, newFileName); err != nil {
		handleError(c, err)
	}

	in, err := os.Open(newFileName)
	if err != nil {
		handleError(c, err)
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
		handleError(c, err)
	}

	c.JSON(http.StatusCreated, fmt.Sprintf("%d participants parsed", count))
}

func (p *ParticipantController) addParticipantToEvent(c *gin.Context, event models.Event, name string, email string) models.Participant {
	participant, err := p.participantService.CreateOrFind(name, email)
	if err != nil {
		handleError(c, err) // TODO: idk?
	}

	err = p.eventService.AddParticipantToEvent(event, participant)
	if err != nil {
		handleError(c, err) // TODO: idk?
	}

	if err := p.mailer.SendWelcomeParticipantMail(event, participant); err != nil {
		handleError(c, err) // TODO: Mail error...
	}

	return participant
}
