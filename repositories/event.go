package repositories

import (
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
	err := p.db.Preload("Groups").Where("token = ?", token).First(&event).Error
	return event, err
}

func (p *EventRepository) FindByToken(token string) (models.Event, error) {
	var event models.Event
	err := p.db.First(&event, "token = ?", token).Error
	return event, err
}

func (p *EventRepository) FindById(eventId uint) (models.Event, error) {
	var event models.Event
	err := p.db.First(&event, eventId).Error
	return event, err
}

func (p *EventRepository) AddParticipant(event models.Event, participant models.Participant) error {
	return p.db.Model(&event).Association("Participants").Append([]*models.Participant{&participant})
}

func (p *EventRepository) IsParticipantInEvent(participantToken string, eventId uint) bool {
	var event models.Event
	p.db.Preload("Participants", "participants.token = ?", participantToken).First(&event, eventId)
	return len(event.Participants) > 0
}

func (p *EventRepository) FindEventWithGroupsAndParticipantsByToken(token string) (models.Event, error) {
	var event models.Event
	err := p.db.Preload("Groups.Participants").First(&event, "token = ?", token).Error
	return event, err
}

func (p *EventRepository) FindEventParticipantsNotInAGroupByToken(token string) ([]models.Participant, error) {
	var event models.Event
	err := p.db.First(&event, "token = ?", token).Error

	groupIds := p.db.
		Model(&models.Group{}).
		Select("id").
		Where("event_id = ?", event.ID)

	/* Participants associated with an event but not in a group
	   select id
	   from participants
	   	inner join participant_groups pg on participants.id = pg.participant_id
	   where group_id in (1, 2) -- groupIds
	*/
	groupedParticipantIds := p.db.
		Model(&models.Participant{}).
		Select("ID").
		Joins("inner join participant_groups pg on participants.id = pg.participant_id").
		Where("group_id in (?)", groupIds)

	err = p.db.
		Preload("Participants", "id not in (?)", groupedParticipantIds).
		First(&event, event.ID).Error

	return event.Participants, err
}
