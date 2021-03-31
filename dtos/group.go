package dtos

import (
	"github.com/tonsV2/event-rooster-api/models"
)

type GroupDTO struct {
	ID              uint `json:"id,string,omitempty"`
	Datetime        string
	MaxParticipants uint
}

func ToGroupDTO(group models.Group) GroupDTO {
	return GroupDTO{ID: group.ID, Datetime: group.Datetime, MaxParticipants: group.MaxParticipants}
}
