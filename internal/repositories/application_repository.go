package repositories

import (
	"github.com/dekkaladiwakar/black-pages-backend/internal/models"

	"gorm.io/gorm"
)

type ApplicationRepository interface {
	Create(application *models.Application) error
	GetByID(id uint) (*models.Application, error)
	GetByJobSeekerID(jobSeekerID uint) ([]models.Application, error)
	GetByJobID(jobID uint) ([]models.Application, error)
	GetByJobAndJobSeeker(jobID, jobSeekerID uint) (*models.Application, error)
	Update(application *models.Application) error
	Delete(id uint) error
	CountByJobID(jobID uint) (int64, error)
	CountByJobSeekerID(jobSeekerID uint) (int64, error)
}

type applicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &applicationRepository{db: db}
}

func (r *applicationRepository) Create(application *models.Application) error {
	return r.db.Create(application).Error
}

func (r *applicationRepository) GetByID(id uint) (*models.Application, error) {
	var application models.Application
	err := r.db.Preload("Job").Preload("JobSeeker").First(&application, id).Error
	if err != nil {
		return nil, err
	}
	return &application, nil
}

func (r *applicationRepository) GetByJobSeekerID(jobSeekerID uint) ([]models.Application, error) {
	var applications []models.Application
	err := r.db.Preload("Job").Preload("Job.Employer").
		Where("job_seeker_id = ?", jobSeekerID).
		Order("applied_at DESC").
		Find(&applications).Error
	return applications, err
}

func (r *applicationRepository) GetByJobID(jobID uint) ([]models.Application, error) {
	var applications []models.Application
	err := r.db.Preload("JobSeeker").
		Where("job_id = ?", jobID).
		Order("applied_at DESC").
		Find(&applications).Error
	return applications, err
}

func (r *applicationRepository) GetByJobAndJobSeeker(jobID, jobSeekerID uint) (*models.Application, error) {
	var application models.Application
	err := r.db.Where("job_id = ? AND job_seeker_id = ?", jobID, jobSeekerID).First(&application).Error
	if err != nil {
		return nil, err
	}
	return &application, nil
}

func (r *applicationRepository) Update(application *models.Application) error {
	return r.db.Save(application).Error
}

func (r *applicationRepository) Delete(id uint) error {
	return r.db.Delete(&models.Application{}, id).Error
}

func (r *applicationRepository) CountByJobID(jobID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Application{}).Where("job_id = ?", jobID).Count(&count).Error
	return count, err
}

func (r *applicationRepository) CountByJobSeekerID(jobSeekerID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Application{}).Where("job_seeker_id = ?", jobSeekerID).Count(&count).Error
	return count, err
}