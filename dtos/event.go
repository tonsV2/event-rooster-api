package models

import "github.com/tonsV2/event-rooster-api/models"

type EventDTO struct {
	ID    uint   `json:"id,string,omitempty"`
	Title string `json:"title"`
	Date  string `json:"date"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type CreateEventDTO struct {
	Title string `json:"title" binding:"required"`
	Date  string `json:"date" binding:"required"`
	Email string `json:"email" binding:"required"`
}

func ToEventDTO(event models.Event) EventDTO {
	return EventDTO{ID: event.ID, Title: event.Title, Date: event.Date, Email: event.Email, Token: event.Token}
}

func ToEventDTOs(events []models.Event) []EventDTO {
	eventDtos := make([]EventDTO, len(events))

	for i, event := range events {
		eventDtos[i] = ToEventDTO(event)
	}

	return eventDtos
}
