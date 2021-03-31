package services

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
	. "github.com/tonsV2/event-rooster-api/models"
	. "github.com/tonsV2/event-rooster-api/repositories"
)

type ParticipantService struct {
	participantRepository ParticipantRepository
}

func ProvideParticipantService(r ParticipantRepository) ParticipantService {
	return ParticipantService{participantRepository: r}
}

func (r *ParticipantService) CreateOrFind(name string, email string) (Participant, error) {
	u4, err := uuid.NewV4()
	if err != nil {
		return Participant{}, err
	}
	token := fmt.Sprint(u4)

	participant := Participant{Name: name, Email: email, Token: token}
	err = r.participantRepository.Create(&participant)
	if err != nil && err.Error() == "UNIQUE constraint failed: participants.email" {
		return r.participantRepository.FindByEmail(email)
	}
	return participant, err
}
