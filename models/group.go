package models

import (
	"gorm.io/gorm"
	"time"
)

type Group struct {
	gorm.Model
	EventID         uint
	Datetime        time.Time
	MaxParticipants uint
	Participants    []Participant `gorm:"many2many:participant_groups;"`
}

type GroupWithParticipantsCount struct {
	ID                 uint
	MaxParticipants    uint
	ActualParticipants uint
}
