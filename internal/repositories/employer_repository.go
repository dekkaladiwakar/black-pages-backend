package repositories

import (
	"github.com/dekkaladiwakar/black-pages-backend/internal/models"

	"gorm.io/gorm"
)

type EmployerRepository interface {
	Create(employer *models.Employer) error
	GetByUserID(userID uint) (*models.Employer, error)
	GetByID(id uint) (*models.Employer, error)
	Update(employer *models.Employer) error
	GetWithFirmProfile(userID uint) (*models.Employer, error)
}

type employerRepository struct {
	db *gorm.DB
}

func NewEmployerRepository(db *gorm.DB) EmployerRepository {
	return &employerRepository{db: db}
}

func (r *employerRepository) Create(employer *models.Employer) error {
	return r.db.Create(employer).Error
}

func (r *employerRepository) GetByUserID(userID uint) (*models.Employer, error) {
	var employer models.Employer
	err := r.db.Where("user_id = ?", userID).First(&employer).Error
	if err != nil {
		return nil, err
	}
	return &employer, nil
}

func (r *employerRepository) GetByID(id uint) (*models.Employer, error) {
	var employer models.Employer
	err := r.db.First(&employer, id).Error
	if err != nil {
		return nil, err
	}
	return &employer, nil
}

func (r *employerRepository) Update(employer *models.Employer) error {
	return r.db.Save(employer).Error
}

func (r *employerRepository) GetWithFirmProfile(userID uint) (*models.Employer, error) {
	var employer models.Employer
	err := r.db.Preload("FirmProfile").Where("user_id = ?", userID).First(&employer).Error
	if err != nil {
		return nil, err
	}
	return &employer, nil
}