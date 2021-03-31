// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package di

import (
	"github.com/tonsV2/event-rooster-api/configurations"
	"github.com/tonsV2/event-rooster-api/controllers"
	"github.com/tonsV2/event-rooster-api/mail"
	"github.com/tonsV2/event-rooster-api/models"
	"github.com/tonsV2/event-rooster-api/repositories"
	"github.com/tonsV2/event-rooster-api/server"
	"github.com/tonsV2/event-rooster-api/services"
)

// Injectors from wire.go:

func BuildServer() server.Server {
	db := models.ProvideDatabase()
	eventRepository := repositories.ProvideEventRepository(db)
	eventService := services.ProvideEventService(eventRepository)
	groupRepository := repositories.ProvideGroupRepository(db)
	groupService := services.ProvideGroupService(groupRepository)
	mailerConfiguration := configurations.ProvideMailerConfiguration()
	mailer := mail.ProvideMailer(mailerConfiguration)
	eventController := controllers.ProvideEventController(eventService, groupService, mailer)
	participantRepository := repositories.ProvideParticipantRepository(db)
	participantService := services.ProvideParticipantService(participantRepository)
	participantController := controllers.ProvideParticipantController(eventService, participantService, mailer)
	serverServer := server.ProvideServer(eventController, participantController)
	return serverServer
}