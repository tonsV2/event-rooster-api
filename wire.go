//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/tonsV2/race-rooster-api/controllers"
	"github.com/tonsV2/race-rooster-api/models"
	"github.com/tonsV2/race-rooster-api/repositories"
	"github.com/tonsV2/race-rooster-api/server"
	"github.com/tonsV2/race-rooster-api/services"
)

func BuildServer() server.Server {
	wire.Build(server.ProvideServer, models.ProvideDatabase,
		repositories.ProvideRaceRepository, services.ProvideRaceService, controllers.ProvideRaceController,
	)
	return server.Server{}
}
