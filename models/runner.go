package models

import "gorm.io/gorm"

type Runner struct {
	gorm.Model
	Name  string
	Email string `gorm:"unique;"`
	Token string `gorm:"index;"`
	Heats []Heat `gorm:"many2many:runner_heats;"`
}
