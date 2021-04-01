package services

import (
	. "github.com/tonsV2/event-rooster-api/models"
	. "github.com/tonsV2/event-rooster-api/repositories"
)

type GroupService struct {
	groupRepository GroupRepository
}

func ProvideGroupService(r GroupRepository) GroupService {
	return GroupService{groupRepository: r}
}

func (g *GroupService) Create(eventId uint, datetime string, maxParticipants uint) (Group, error) {
	group := Group{EventID: eventId, Datetime: datetime, MaxParticipants: maxParticipants}
	err := g.groupRepository.Create(&group)
	return group, err
}

func (g *GroupService) FindGroupsWithParticipantsCountByEventId(id uint) ([]GroupWithParticipantsCount, error) {
	var groups []GroupWithParticipantsCount
	err := g.groupRepository.FindGroupsWithParticipantsCountByEventId(id, &groups)
	if err != nil {
		return groups, nil
	} else {
		return groups, err
	}
}
