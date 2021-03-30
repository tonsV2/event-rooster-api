package models

import "gorm.io/gorm"

type Heat struct {
	gorm.Model
	RaceID          uint
	Datetime        string
	MaxParticipants uint
	Runners         []Runner `gorm:"many2many:runner_heats;"`
}
