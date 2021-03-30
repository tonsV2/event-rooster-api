package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tonsV2/race-rooster-api/controllers"
)

type Server struct {
	Engine *gin.Engine
}

func ProvideServer(raceController controllers.RaceController) Server {
	r := gin.Default()
	r.Use(cors.Default())

	initializeRaceController(r, raceController)

	return Server{Engine: r}
}

func initializeRaceController(r *gin.Engine, raceController controllers.RaceController) {
	r.POST("/races", raceController.CreateRace)
}

func (s *Server) Run() {
	_ = s.Engine.Run()
}
