package dtos

import "github.com/tonsV2/event-rooster-api/models"

type EventDTO struct {
	ID       uint   `json:"id,string,omitempty"`
	Title    string `json:"title"`
	Datetime string `json:"datetime"`
	Email    string `json:"email"`
	Token    string `json:"token,omitempty"`
}

type CreateEventDTO struct {
	Title    string `json:"title" binding:"required"`
	Datetime string `json:"datetime" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type EventWithGroupsDTO struct {
	ID       uint       `json:"id,string,omitempty"`
	Title    string     `json:"title"`
	Datetime string     `json:"datetime"`
	Groups   []GroupDTO `json:"groups"`
}

func ToEventWithGroupsDTO(event models.Event) EventWithGroupsDTO {
	groupDtos := make([]GroupDTO, len(event.Groups))

	for i, group := range event.Groups {
		groupDtos[i] = ToGroupDTO(group)
	}

	return EventWithGroupsDTO{ID: event.ID, Title: event.Title, Datetime: event.Datetime, Groups: groupDtos}
}

func ToEventWithGroupsAndParticipantsDTO(event models.Event) EventWithGroupsDTO {
	groupDtos := make([]GroupDTO, len(event.Groups))

	for i, group := range event.Groups {
		groupDtos[i] = ToGroupWithParticipantsDTO(group)
	}

	return EventWithGroupsDTO{ID: event.ID, Title: event.Title, Datetime: event.Datetime, Groups: groupDtos}
}

func ToEventDTO(event models.Event) EventDTO {
	return EventDTO{ID: event.ID, Title: event.Title, Datetime: event.Datetime, Email: event.Email, Token: event.Token}
}
