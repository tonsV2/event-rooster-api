package repositories

import (
	"github.com/tonsV2/event-rooster-api/models"
	"gorm.io/gorm"
)

type ParticipantRepository struct {
	db *gorm.DB
}

func ProvideParticipantRepository(DB *gorm.DB) ParticipantRepository {
	return ParticipantRepository{db: DB}
}

func (p *ParticipantRepository) Create(participant *models.Participant) error {
	return p.db.Create(&participant).Error
}

func (p *ParticipantRepository) FindByEmail(email string) (models.Participant, error) {
	var participant models.Participant
	err := p.db.Find(&participant, "email = ?", email).Error
	return participant, err
}

func (p *ParticipantRepository) FindByToken(token string) (models.Participant, error) {
	var participant models.Participant
	err := p.db.Find(&participant, "token = ?", token).Error
	return participant, err
}
