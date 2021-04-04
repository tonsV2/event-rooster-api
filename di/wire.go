//+build wireinject

package di

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
		server.ProvideServer,
		controllers.ProvideHealthController,
		models.ProvideDatabase,
		configurations.ProvideMailerConfiguration, mail.ProvideMailer,
		repositories.ProvideEventRepository, services.ProvideEventService, controllers.ProvideEventController,
		repositories.ProvideGroupRepository, services.ProvideGroupService, controllers.ProvideGroupController,
		repositories.ProvideParticipantRepository, services.ProvideParticipantService, controllers.ProvideParticipantController,
	)
	return server.Server{}
}
