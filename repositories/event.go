package repositories

import (
	. "github.com/tonsV2/event-rooster-api/models"
	"gorm.io/gorm"
)

type EventRepository struct {
	DB *gorm.DB
}

func ProvideEventRepository(DB *gorm.DB) EventRepository {
	return EventRepository{DB: DB}
}

func (p *EventRepository) Create(event Event) Event {
	p.DB.Create(&event)
	return event
}
