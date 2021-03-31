package models

import "github.com/tonsV2/race-rooster-api/models"

type RaceDTO struct {
	ID    uint   `json:"id,string,omitempty"`
	Title string `json:"title"`
	Date  string `json:"date"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type CreateRaceDTO struct {
	Title string `json:"title" binding:"required"`
	Date  string `json:"date" binding:"required"`
	Email string `json:"email" binding:"required"`
}

func ToRaceDTO(race models.Race) RaceDTO {
	return RaceDTO{ID: race.ID, Title: race.Title, Date: race.Date, Email: race.Email, Token: race.Token}
}

func ToRaceDTOs(races []models.Race) []RaceDTO {
	raceDtos := make([]RaceDTO, len(races))

	for i, race := range races {
		raceDtos[i] = ToRaceDTO(race)
	}

	return raceDtos
}
