package repositories

import (
	. "github.com/tonsV2/race-rooster-api/models"
	"gorm.io/gorm"
)

type RaceRepository struct {
	DB *gorm.DB
}

func ProvideRaceRepository(DB *gorm.DB) RaceRepository {
	return RaceRepository{DB: DB}
}

func (p *RaceRepository) Create(race Race) Race {
	p.DB.Create(&race)
	return race
}
