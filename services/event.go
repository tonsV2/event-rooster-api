package services

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
	. "github.com/tonsV2/event-rooster-api/models"
	. "github.com/tonsV2/event-rooster-api/repositories"
)

type EventService struct {
	eventRepository EventRepository
}

func ProvideEventService(r EventRepository) EventService {
	return EventService{eventRepository: r}
}

func (r *EventService) Create(title string, date string, email string) (Event, error) {
	u4, err := uuid.NewV4()
	if err != nil {
		return Event{}, err
	}
	token := fmt.Sprint(u4)

	event := Event{Title: title, Date: date, Token: token, Email: email}
	r.eventRepository.Create(event)
	return event, nil
}

func (r *EventService) FindEventWithGroupsByToken(token string) (Event, error) {
	return r.eventRepository.FindEventWithGroupsByToken(token)
}
