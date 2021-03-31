package models

import "gorm.io/gorm"

type Event struct {
	gorm.Model
	Title   string
	Date    string
	Email   string
	Token   string   `gorm:"index;"`
	Groups  []Group  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Runners []Runner `gorm:"many2many:event_runners;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
