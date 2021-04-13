package models

import (
	"gorm.io/gorm"
	"time"
)

type Group struct {
	gorm.Model
	GID             string `gorm:"uniqueIndex:gid_unique_on_event"`
	EventID         uint   `gorm:"uniqueIndex:gid_unique_on_event"`
	Datetime        time.Time
	MaxParticipants uint
	Participants    []Participant `gorm:"many2many:participant_groups;"`
}

type GroupWithParticipantsCount struct {
	ID                 uint
	GID                string
	Datetime           time.Time
	MaxParticipants    uint
	ActualParticipants uint
}
