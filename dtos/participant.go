package dtos

import (
	"github.com/tonsV2/event-rooster-api/models"
)

type ParticipantDTO struct {
	ID    uint   `json:"id,string,omitempty"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type CreateParticipantDTO struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

func ToParticipantDTO(participant models.Participant) ParticipantDTO {
	return ParticipantDTO{ID: participant.ID, Name: participant.Name, Email: participant.Email, Token: participant.Token}
}
