package models

import "gorm.io/gorm"

type Participant struct {
	gorm.Model
	Name   string
	Email  string  `gorm:"unique;"`
	Token  string  `gorm:"index;"`
	Groups []Group `gorm:"many2many:participant_groups;"`
}
