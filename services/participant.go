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

func (p *ParticipantService) CreateOrFind(name string, email string) (Participant, error) {
	u4, err := uuid.NewV4()
	if err != nil {
		return Participant{}, err
	}
	token := fmt.Sprint(u4)

	participant := Participant{Name: name, Email: email, Token: token}
	err = p.participantRepository.Create(&participant)
	// TODO: Better handling of error
	// SQLite and Postgresql returns different error messages
	// SQLite: UNIQUE constraint failed: participants.email
	// Postgresql: duplicate key value violates unique constraint "participants_email_key"
	// So currently the below code tries p.participantRepository.FindByEmail() in case of any kind of error from p.participantRepository.Create()
	// It should only be for errors on "unique constraint", other kind of errors should be handled differently
	if err != nil /*&& err.Error() == "UNIQUE constraint failed: participants.email"*/ {
		return p.participantRepository.FindByEmail(email)
	}
	return participant, err
}

func (p *ParticipantService) Create(name string, email string) (Participant, error) {
	u4, err := uuid.NewV4()
	if err != nil {
		return Participant{}, err
	}
	token := fmt.Sprint(u4)

	participant := Participant{Name: name, Email: email, Token: token}
	err = p.participantRepository.Create(&participant)
	return participant, err
}

func (p *ParticipantService) FindByToken(token string) (Participant, error) {
	return p.participantRepository.FindByToken(token)
}

func (p *ParticipantService) FindByEmail(email string) (Participant, error) {
	return p.participantRepository.FindByEmail(email)
}
