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
	err = r.eventRepository.Create(&event)
	return event, err
}

func (r *EventService) FindEventWithGroupsByToken(token string) (Event, error) {
	return r.eventRepository.FindEventWithGroupsByToken(token)
}

func (r *EventService) FindByToken(token string) (Event, error) {
	return r.eventRepository.FindByToken(token)
}

func (r *EventService) AddParticipantToEvent(event Event, participant Participant) error {
	return r.eventRepository.AddParticipant(event, participant)
}

func (r *EventService) FindByIdAndParticipantToken(eventId uint, participantToken string) (Event, error) {
	return r.eventRepository.FindByIdAndParticipantToken(eventId, participantToken)
}

func (r *EventService) FindEventWithGroupsAndParticipantsByToken(token string) (Event, error) {
	return r.eventRepository.FindEventWithGroupsAndParticipantsByToken(token)
}
