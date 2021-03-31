//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/tonsV2/event-rooster-api/configurations"
	"github.com/tonsV2/event-rooster-api/controllers"
	"github.com/tonsV2/event-rooster-api/mail"
	"github.com/tonsV2/event-rooster-api/models"
	"github.com/tonsV2/event-rooster-api/repositories"
	"github.com/tonsV2/event-rooster-api/server"
	"github.com/tonsV2/event-rooster-api/services"
)

func BuildServer() server.Server {
	wire.Build(
		configurations.ProvideMailerConfiguration, mail.ProvideMailer,
		server.ProvideServer, models.ProvideDatabase,
		repositories.ProvideEventRepository, services.ProvideEventService, controllers.ProvideEventController,
	)
	return server.Server{}
}
