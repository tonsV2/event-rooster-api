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

func (p *EventRepository) Create(event models.Event) models.Event {
	p.db.Create(&event)
	return event
}

func (p *EventRepository) FindEventWithGroupsByToken(token string) (models.Event, error) {
	var event models.Event
	if err := p.db.Preload("Groups").Where("token = ?", token).Find(&event).Error; err != nil {
		return event, err
	}
	return event, nil
}
