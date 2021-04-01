package dtos

import (
	"github.com/tonsV2/event-rooster-api/models"
)

type GroupDTO struct {
	ID              uint   `json:"id,string,omitempty"`
	Datetime        string `json:"datetime"`
	MaxParticipants uint   `json:"maxParticipants"`
}

type CreateGroupDTO struct {
	Datetime        string `json:"datetime" binding:"required"`
	MaxParticipants uint   `json:"maxParticipants" binding:"required"`
}

func ToGroupDTO(group models.Group) GroupDTO {
	return GroupDTO{ID: group.ID, Datetime: group.Datetime, MaxParticipants: group.MaxParticipants}
}

type GroupWithParticipantsCountDTO struct {
	ID                 uint `json:"id,string,omitempty"`
	MaxParticipants    uint `json:"maxParticipants"`
	ActualParticipants uint `json:"actualParticipants"`
}

func ToGroupWithParticipantsCountDTO(groupWithParticipantsCount models.GroupWithParticipantsCount) GroupWithParticipantsCountDTO {
	return GroupWithParticipantsCountDTO{
		ID:                 groupWithParticipantsCount.ID,
		MaxParticipants:    groupWithParticipantsCount.MaxParticipants,
		ActualParticipants: groupWithParticipantsCount.ActualParticipants,
	}
}

func ToGroupsWithParticipantsCountDTO(groupsWithParticipantsCount []models.GroupWithParticipantsCount) []GroupWithParticipantsCountDTO {
	groupDtos := make([]GroupWithParticipantsCountDTO, len(groupsWithParticipantsCount))

	for i, group := range groupsWithParticipantsCount {
		groupDtos[i] = ToGroupWithParticipantsCountDTO(group)
	}
	return groupDtos
}
