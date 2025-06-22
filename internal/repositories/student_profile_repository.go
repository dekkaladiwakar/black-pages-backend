package repositories

import (
	"github.com/dekkaladiwakar/black-pages-backend/internal/models"

	"gorm.io/gorm"
)

type StudentProfileRepository interface {
	Create(profile *models.StudentProfile) error
	GetByJobSeekerID(jobSeekerID uint) (*models.StudentProfile, error)
	Update(profile *models.StudentProfile) error
	Delete(jobSeekerID uint) error
	GetByID(id uint) (*models.StudentProfile, error)
}

type studentProfileRepository struct {
	db *gorm.DB
}

func NewStudentProfileRepository(db *gorm.DB) StudentProfileRepository {
	return &studentProfileRepository{db: db}
}

func (r *studentProfileRepository) Create(profile *models.StudentProfile) error {
	return r.db.Create(profile).Error
}

func (r *studentProfileRepository) GetByJobSeekerID(jobSeekerID uint) (*models.StudentProfile, error) {
	var profile models.StudentProfile
	err := r.db.Preload("JobSeeker").Where("job_seeker_id = ?", jobSeekerID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *studentProfileRepository) Update(profile *models.StudentProfile) error {
	return r.db.Save(profile).Error
}

func (r *studentProfileRepository) Delete(jobSeekerID uint) error {
	return r.db.Where("job_seeker_id = ?", jobSeekerID).Delete(&models.StudentProfile{}).Error
}

func (r *studentProfileRepository) GetByID(id uint) (*models.StudentProfile, error) {
	var profile models.StudentProfile
	err := r.db.Preload("JobSeeker").First(&profile, id).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}