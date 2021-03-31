package repositories

import (
	"github.com/tonsV2/event-rooster-api/models"
	"gorm.io/gorm"
)

type GroupRepository struct {
	db *gorm.DB
}

func ProvideGroupRepository(DB *gorm.DB) GroupRepository {
	return GroupRepository{db: DB}
}

func (p *GroupRepository) Create(group *models.Group) error {
	return p.db.Create(&group).Error
}
