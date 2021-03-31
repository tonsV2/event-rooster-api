package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	EventID         uint
	Datetime        string
	MaxParticipants uint
	Runners         []Runner `gorm:"many2many:runner_groups;"`
}
