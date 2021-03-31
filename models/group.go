package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	EventID         uint
	Datetime        string
	MaxParticipants uint
	Participants    []Participant `gorm:"many2many:participant_groups;"`
}
