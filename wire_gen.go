// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"github.com/tonsV2/race-rooster-api/configurations"
	"github.com/tonsV2/race-rooster-api/controllers"
	"github.com/tonsV2/race-rooster-api/mail"
	"github.com/tonsV2/race-rooster-api/models"
	"github.com/tonsV2/race-rooster-api/repositories"
	"github.com/tonsV2/race-rooster-api/server"
	"github.com/tonsV2/race-rooster-api/services"
)

// Injectors from wire.go:

func BuildServer() server.Server {
	db := models.ProvideDatabase()
	raceRepository := repositories.ProvideRaceRepository(db)
	raceService := services.ProvideRaceService(raceRepository)
	mailerConfiguration := configurations.ProvideMailerConfiguration()
	mailer := mail.ProvideMailer(mailerConfiguration)
	raceController := controllers.ProvideRaceController(raceService, mailer)
	serverServer := server.ProvideServer(raceController)
	return serverServer
}
