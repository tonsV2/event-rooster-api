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

func (r *GroupService) Create(eventId uint, datetime string, maxParticipants uint) Group {
	group := Group{EventID: eventId, Datetime: datetime, MaxParticipants: maxParticipants}
	r.groupRepository.Create(&group)
	return group
}
