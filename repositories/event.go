package repositories

import (
	"errors"
	"github.com/tonsV2/event-rooster-api/models"
	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func ProvideEventRepository(DB *gorm.DB) EventRepository {
	return EventRepository{db: DB}
}

func (p *EventRepository) Create(event *models.Event) error {
	return p.db.Create(&event).Error
}

func (p *EventRepository) FindEventWithGroupsByToken(token string) (models.Event, error) {
	var event models.Event
	err := p.db.Preload("Groups").Where("token = ?", token).Find(&event).Error
	return event, err
}

func (p *EventRepository) FindByToken(token string) (models.Event, error) {
	var event models.Event
	err := p.db.Find(&event, "token = ?", token).Error
	return event, err
}

func (p *EventRepository) AddParticipant(event models.Event, participant models.Participant) error {
	return p.db.Model(&event).Association("Participants").Append([]*models.Participant{&participant})
}

func (p *EventRepository) FindByIdAndParticipantToken(eventId uint, participantToken string) (models.Event, error) {
	var event models.Event
	p.db.Preload("Participants", "participants.token = ?", participantToken).Find(&event, eventId)
	if len(event.Participants) < 1 {
		return event, errors.New("participant not associated with event")
	} else {
		return event, nil
	}
}

func (p *EventRepository) FindEventWithGroupsAndParticipantsByToken(token string) (models.Event, error) {
	var event models.Event
	err := p.db.Preload("Groups.Participants").Find(&event, "token = ?", token).Error
	return event, err
}
