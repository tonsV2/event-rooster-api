package models

import "gorm.io/gorm"

type Runner struct {
	gorm.Model
	Name   string
	Email  string  `gorm:"unique;"`
	Token  string  `gorm:"index;"`
	Groups []Group `gorm:"many2many:runner_groups;"`
}
