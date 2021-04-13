package dtos

import (
	"github.com/tonsV2/event-rooster-api/models"
	"time"
)

type CreateGroupDTO struct {
	GID             string    `json:"gid" binding:"required"`
	Datetime        time.Time `json:"datetime" binding:"required"`
	MaxParticipants uint      `json:"maxParticipants" binding:"required"`
}

type GroupDTO struct {
	ID              uint             `json:"id,omitempty"`
	GID             string           `json:"gid"`
	Datetime        time.Time        `json:"datetime"`
	MaxParticipants uint             `json:"maxParticipants"`
	Participants    []ParticipantDTO `json:"participants,omitempty"`
	CreatedAt       time.Time        `json:"createdAt"`
}

func ToGroupDTO(group models.Group) GroupDTO {
	return GroupDTO{ID: group.ID, GID: group.GID, Datetime: group.Datetime, MaxParticipants: group.MaxParticipants, CreatedAt: group.CreatedAt}
}

func ToGroupWithParticipantsDTO(group models.Group) GroupDTO {
	participantDtos := make([]ParticipantDTO, len(group.Participants))

	for i, participant := range group.Participants {
		participantDtos[i] = ToParticipantWithoutTokenDTO(participant)
	}

	return GroupDTO{
		ID:              group.ID,
		GID:             group.GID,
		Datetime:        group.Datetime,
		MaxParticipants: group.MaxParticipants,
		Participants:    participantDtos,
	}
}

type GroupWithParticipantsCountDTO struct {
	ID                 uint      `json:"id,omitempty"`
	GID                string    `json:"gid"`
	Datetime           time.Time `json:"datetime"`
	MaxParticipants    uint      `json:"maxParticipants"`
	ActualParticipants uint      `json:"actualParticipants"`
}

func ToGroupWithParticipantsCountDTO(groupWithParticipantsCount models.GroupWithParticipantsCount) GroupWithParticipantsCountDTO {
	return GroupWithParticipantsCountDTO{
		ID:                 groupWithParticipantsCount.ID,
		GID:                groupWithParticipantsCount.GID,
		Datetime:           groupWithParticipantsCount.Datetime,
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
