package services

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
	. "github.com/tonsV2/race-rooster-api/models"
	. "github.com/tonsV2/race-rooster-api/repositories"
)

type RaceService struct {
	RaceRepository RaceRepository
}

func ProvideRaceService(r RaceRepository) RaceService {
	return RaceService{RaceRepository: r}
}

func (r *RaceService) Create(title string, date string) (Race, error) {
	u4, err := uuid.NewV4()
	if err != nil {
		return Race{}, err
	}
	token := fmt.Sprint(u4)

	race := Race{Title: title, Date: date, Token: token}
	r.RaceRepository.Create(race)
	return race, nil
}
