package repositories

import (
	"github.com/dekkaladiwakar/black-pages-backend/internal/models"

	"gorm.io/gorm"
)

type JobSeekerRepository interface {
	Create(jobSeeker *models.JobSeeker) error
	GetByUserID(userID uint) (*models.JobSeeker, error)
	GetByID(id uint) (*models.JobSeeker, error)
	Update(jobSeeker *models.JobSeeker) error
	Delete(id uint) error
	GetWithStudentProfile(userID uint) (*models.JobSeeker, error)
}

type jobSeekerRepository struct {
	db *gorm.DB
}

func NewJobSeekerRepository(db *gorm.DB) JobSeekerRepository {
	return &jobSeekerRepository{db: db}
}

func (r *jobSeekerRepository) Create(jobSeeker *models.JobSeeker) error {
	return r.db.Create(jobSeeker).Error
}

func (r *jobSeekerRepository) GetByUserID(userID uint) (*models.JobSeeker, error) {
	var jobSeeker models.JobSeeker
	err := r.db.Where("user_id = ?", userID).First(&jobSeeker).Error
	if err != nil {
		return nil, err
	}
	return &jobSeeker, nil
}

func (r *jobSeekerRepository) GetByID(id uint) (*models.JobSeeker, error) {
	var jobSeeker models.JobSeeker
	err := r.db.First(&jobSeeker, id).Error
	if err != nil {
		return nil, err
	}
	return &jobSeeker, nil
}

func (r *jobSeekerRepository) Update(jobSeeker *models.JobSeeker) error {
	return r.db.Save(jobSeeker).Error
}

func (r *jobSeekerRepository) Delete(id uint) error {
	return r.db.Delete(&models.JobSeeker{}, id).Error
}

func (r *jobSeekerRepository) GetWithStudentProfile(userID uint) (*models.JobSeeker, error) {
	var jobSeeker models.JobSeeker
	err := r.db.Preload("StudentProfile").Where("user_id = ?", userID).First(&jobSeeker).Error
	if err != nil {
		return nil, err
	}
	return &jobSeeker, nil
}