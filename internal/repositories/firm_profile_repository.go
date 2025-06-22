package repositories

import (
	"github.com/dekkaladiwakar/black-pages-backend/internal/models"

	"gorm.io/gorm"
)

type FirmProfileRepository interface {
	Create(profile *models.FirmProfile) error
	GetByEmployerID(employerID uint) (*models.FirmProfile, error)
	Update(profile *models.FirmProfile) error
	Delete(employerID uint) error
	GetByID(id uint) (*models.FirmProfile, error)
}

type firmProfileRepository struct {
	db *gorm.DB
}

func NewFirmProfileRepository(db *gorm.DB) FirmProfileRepository {
	return &firmProfileRepository{db: db}
}

func (r *firmProfileRepository) Create(profile *models.FirmProfile) error {
	return r.db.Create(profile).Error
}

func (r *firmProfileRepository) GetByEmployerID(employerID uint) (*models.FirmProfile, error) {
	var profile models.FirmProfile
	err := r.db.Preload("Employer").Where("employer_id = ?", employerID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *firmProfileRepository) Update(profile *models.FirmProfile) error {
	return r.db.Save(profile).Error
}

func (r *firmProfileRepository) Delete(employerID uint) error {
	return r.db.Where("employer_id = ?", employerID).Delete(&models.FirmProfile{}).Error
}

func (r *firmProfileRepository) GetByID(id uint) (*models.FirmProfile, error) {
	var profile models.FirmProfile
	err := r.db.Preload("Employer").First(&profile, id).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}