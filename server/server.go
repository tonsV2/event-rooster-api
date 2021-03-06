package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tonsV2/event-rooster-api/controllers"
)

type Server struct {
	Engine *gin.Engine
}

func ProvideServer(
	healthController controllers.HealthController,
	eventController controllers.EventController,
	participantController controllers.ParticipantController,
	groupController controllers.GroupController) Server {

	r := gin.Default()
	r.Use(cors.Default())

	initializeEventController(r, healthController, eventController, participantController, groupController)

	return Server{Engine: r}
}

func initializeEventController(r *gin.Engine,
	healthController controllers.HealthController,
	eventController controllers.EventController,
	participantController controllers.ParticipantController,
	groupController controllers.GroupController) {

	r.GET("/health", healthController.GetHealthStatus)
	r.GET("/", healthController.GetHealthStatus)

	r.POST("/events", eventController.CreateEvent)
	r.GET("/events", eventController.GetEventWithGroupsAndParticipantsByToken)
	r.GET("/events/groups", eventController.FindEventWithGroupsByToken)
	r.POST("/events/groups", eventController.AddGroupToEventByToken)

	r.POST("/participants", participantController.AddParticipantToEventByToken)
	r.POST("/participants/csv", participantController.AddParticipantsCSVToEventByToken)
	r.POST("/participants/groups", participantController.AddParticipantToGroupByToken)
	r.GET("/participants/not-in-groups", eventController.FindEventParticipantsNotInAGroupByToken)

	r.GET("/events/groups-count", groupController.GetGroupsWithParticipantsCountByEventIdAndParticipantToken)
}

func (s *Server) Run() {
	_ = s.Engine.Run()
}
