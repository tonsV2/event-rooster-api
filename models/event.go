package models

import (
	"gorm.io/gorm"
	"time"
)

type Event struct {
	gorm.Model
	Title        string
	Datetime     time.Time
	Email        string
	Token        string        `gorm:"index;"`
	Groups       []Group       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Participants []Participant `gorm:"many2many:event_participants;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
