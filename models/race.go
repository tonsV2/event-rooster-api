package models

import "gorm.io/gorm"

type Race struct {
	gorm.Model
	Title   string
	Date    string
	Email   string
	Token   string   `gorm:"index;"`
	Heats   []Heat   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Runners []Runner `gorm:"many2many:race_runners;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
