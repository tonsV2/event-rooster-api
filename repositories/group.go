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
		Select("id, max_participants, count(participant_id) as actual_participants, g_id, datetime").
		Where("event_id = ?", id).
		Joins("left join participant_groups pg on groups.id = pg.group_id").
		Group("groups.id").
		Order("g_id").
		Scan(&groups).Error
}

func (g *GroupRepository) FindById(id string) (models.Group, error) {
	var group models.Group
	err := g.db.First(&group, id).Error
	return group, err
}

func (g *GroupRepository) AddParticipant(group models.Group, participant models.Participant) error {
	return g.db.Model(&group).Association("Participants").Append([]*models.Participant{&participant})
}

func (g *GroupRepository) RemoveParticipant(group models.Group, participant models.Participant) error {
	return g.db.Model(&group).Association("Participants").Delete([]*models.Participant{&participant})
}

func (g *GroupRepository) FindParticipantGroups(group models.Group, participant models.Participant) ([]models.Group, error) {
	eventGroupIds := g.db.
		Model(&models.Group{}).
		Select("id").
		Where("event_id = ?", group.EventID)

	participantGroupIds := g.db.
		Table("participant_groups").
		Select("group_id").
		Where("participant_id = ? and group_id in (?)", participant.ID, eventGroupIds)

	var groups []models.Group
	err := g.db.
		Model(&models.Group{}).
		Where("id in (?)", participantGroupIds).
		Find(&groups).Error

	return groups, err
}
