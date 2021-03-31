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
	if err := p.db.Preload("Groups").Where("token = ?", token).Find(&event).Error; err != nil {
		return event, err
	}
	return event, nil
}

func (p *EventRepository) FindByToken(token string) (models.Event, error) {
	var event models.Event
	if err := p.db.Find(&event, "token = ?", token).Error; err != nil {
		return event, err
	}
	return event, nil
}

func (p *EventRepository) AddParticipant(event models.Event, participant models.Participant) error {
	return p.db.Model(&event).Association("Participants").Append([]*models.Participant{&participant})
}
