package services

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
	. "github.com/tonsV2/event-rooster-api/models"
	. "github.com/tonsV2/event-rooster-api/repositories"
	"time"
)

type EventService struct {
	eventRepository EventRepository
}

func ProvideEventService(r EventRepository) EventService {
	return EventService{eventRepository: r}
}

func (r *EventService) Create(title string, datetime time.Time, email string) (Event, error) {
	u4, err := uuid.NewV4()
	if err != nil {
		return Event{}, err
	}
	token := fmt.Sprint(u4)

	event := Event{Title: title, Datetime: datetime, Token: token, Email: email}
	err = r.eventRepository.Create(&event)
	return event, err
}

func (r *EventService) FindEventWithGroupsByToken(token string) (Event, error) {
	return r.eventRepository.FindEventWithGroupsByToken(token)
}

func (r *EventService) FindByToken(token string) (Event, error) {
	return r.eventRepository.FindByToken(token)
}

func (r *EventService) FindById(eventId uint) (Event, error) {
	return r.eventRepository.FindById(eventId)
}

func (r *EventService) AddParticipantToEvent(event Event, participant Participant) error {
	return r.eventRepository.AddParticipant(event, participant)
}

func (r *EventService) IsParticipantInEvent(participantToken string, eventId uint) bool {
	return r.eventRepository.IsParticipantInEvent(participantToken, eventId)
}

func (r *EventService) FindEventWithGroupsAndParticipantsByToken(token string) (Event, error) {
	return r.eventRepository.FindEventWithGroupsAndParticipantsByToken(token)
}

func (r *EventService) FindEventParticipantsNotInAGroupByToken(token string) ([]Participant, error) {
	return r.eventRepository.FindEventParticipantsNotInAGroupByToken(token)
}
