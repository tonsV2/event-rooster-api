package services

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
	. "github.com/tonsV2/event-rooster-api/models"
	. "github.com/tonsV2/event-rooster-api/repositories"
)

type EventService struct {
	EventRepository EventRepository
}

func ProvideEventService(r EventRepository) EventService {
	return EventService{EventRepository: r}
}

func (r *EventService) Create(title string, date string, email string) (Event, error) {
	u4, err := uuid.NewV4()
	if err != nil {
		return Event{}, err
	}
	token := fmt.Sprint(u4)

	event := Event{Title: title, Date: date, Token: token, Email: email}
	r.EventRepository.Create(event)
	return event, nil
}
