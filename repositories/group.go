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

func (g *GroupRepository) Create(group *models.Group) error {
	return g.db.Create(&group).Error
}

func (g *GroupRepository) FindGroupsWithParticipantsCountByEventId(id uint, groups *[]models.GroupWithParticipantsCount) error {
	/*
	   select id, max_participants, count(participant_id) as actual_participants
	   from groups
	            left join participant_groups rh on groups.id = rh.group_id
	   where event_id = 1
	   group by group_id
	*/

	return g.db.
		Model(models.Group{}).
		Select("id, max_participants, count(participant_id) as actual_participants").
		Where("event_id = ?", id).
		Joins("left join participant_groups pg on groups.id = pg.group_id").
		Group("group_id").
		Scan(&groups).Error
}
