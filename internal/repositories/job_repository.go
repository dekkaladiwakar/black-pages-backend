package repositories

import (
	"github.com/dekkaladiwakar/black-pages-backend/internal/models"

	"gorm.io/gorm"
)

type JobFilters struct {
	EmployerID      uint   
	Industry        string 
	JobType         string 
	City            string 
	TargetAudience  string 
	EmploymentMode  string 
	IsPaid          *bool  
	IsActive        *bool  
	OrderBy         string 
	OrderDirection  string 
	Limit           int    
}

type JobRepository interface {
	Create(job *models.Job) error
	GetByID(id uint) (*models.Job, error)
	GetByEmployerID(employerID uint) ([]models.Job, error)
	Update(job *models.Job) error
	Delete(id uint) error
	GetAll() ([]models.Job, error)
	GetWithFilters(filters JobFilters) ([]models.Job, error)
	CountByEmployerID(employerID uint) (int64, error)
	GetDistinctIndustries() ([]string, error)
	GetDistinctCities() ([]string, error)
}

type jobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) JobRepository {
	return &jobRepository{db: db}
}

func (r *jobRepository) Create(job *models.Job) error {
	return r.db.Create(job).Error
}

func (r *jobRepository) GetByID(id uint) (*models.Job, error) {
	var job models.Job
	err := r.db.Preload("Employer").First(&job, id).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *jobRepository) GetByEmployerID(employerID uint) ([]models.Job, error) {
	var jobs []models.Job
	err := r.db.Where("employer_id = ?", employerID).Order("created_at DESC").Find(&jobs).Error
	return jobs, err
}

func (r *jobRepository) Update(job *models.Job) error {
	return r.db.Save(job).Error
}

func (r *jobRepository) Delete(id uint) error {
	return r.db.Delete(&models.Job{}, id).Error
}

func (r *jobRepository) GetAll() ([]models.Job, error) {
	var jobs []models.Job
	err := r.db.Preload("Employer").Where("is_active = ?", true).Order("created_at DESC").Find(&jobs).Error
	return jobs, err
}

func (r *jobRepository) GetWithFilters(filters JobFilters) ([]models.Job, error) {
	builder := NewJobQueryBuilder(r.db)
	
	if filters.EmployerID > 0 {
		builder = builder.WithEmployerID(filters.EmployerID)
	}
	
	builder = builder.WithIndustry(filters.Industry).
		WithJobType(filters.JobType).
		WithCity(filters.City).
		WithTargetAudience(filters.TargetAudience).
		WithEmploymentMode(filters.EmploymentMode).
		OrderBy(filters.OrderBy, filters.OrderDirection).
		Limit(filters.Limit)

	if filters.IsPaid != nil {
		builder = builder.WithPaidStatus(*filters.IsPaid)
	}

	// Only apply IsActive filter if explicitly provided
	builder = builder.WithIsActive(filters.IsActive)

	var jobs []models.Job
	err := builder.Build().Preload("Employer").Find(&jobs).Error
	return jobs, err
}

func (r *jobRepository) CountByEmployerID(employerID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Job{}).Where("employer_id = ?", employerID).Count(&count).Error
	return count, err
}

type JobQueryBuilder struct {
	query *gorm.DB
}

func NewJobQueryBuilder(db *gorm.DB) *JobQueryBuilder {
	return &JobQueryBuilder{
		query: db.Model(&models.Job{}),
	}
}

func (b *JobQueryBuilder) WithIndustry(industry string) *JobQueryBuilder {
	if industry != "" {
		b.query = b.query.Where("industry = ?", industry)
	}
	return b
}

func (b *JobQueryBuilder) WithJobType(jobType string) *JobQueryBuilder {
	if jobType != "" {
		b.query = b.query.Where("job_type = ?", jobType)
	}
	return b
}

func (b *JobQueryBuilder) WithCity(city string) *JobQueryBuilder {
	if city != "" {
		b.query = b.query.Where("city ILIKE ?", "%"+city+"%")
	}
	return b
}

func (b *JobQueryBuilder) WithTargetAudience(audience string) *JobQueryBuilder {
	if audience != "" {
		b.query = b.query.Where("target_audience = ? OR target_audience = ?", audience, "any")
	}
	return b
}

func (b *JobQueryBuilder) WithEmploymentMode(mode string) *JobQueryBuilder {
	if mode != "" {
		b.query = b.query.Where("employment_mode = ?", mode)
	}
	return b
}

func (b *JobQueryBuilder) WithPaidStatus(isPaid bool) *JobQueryBuilder {
	b.query = b.query.Where("is_paid = ?", isPaid)
	return b
}

func (b *JobQueryBuilder) WithIsActive(isActive *bool) *JobQueryBuilder {
	if isActive != nil {
		b.query = b.query.Where("is_active = ?", *isActive)
	}
	return b
}

func (b *JobQueryBuilder) WithEmployerID(employerID uint) *JobQueryBuilder {
	if employerID > 0 {
		b.query = b.query.Where("employer_id = ?", employerID)
	}
	return b
}

func (b *JobQueryBuilder) OrderBy(field string, direction string) *JobQueryBuilder {
	if field != "" && direction != "" {
		b.query = b.query.Order(field + " " + direction)
	} else {
		b.query = b.query.Order("created_at DESC")
	}
	return b
}

func (b *JobQueryBuilder) Limit(limit int) *JobQueryBuilder {
	if limit > 0 {
		b.query = b.query.Limit(limit)
	}
	return b
}

func (b *JobQueryBuilder) Build() *gorm.DB {
	return b.query
}

func (r *jobRepository) GetDistinctIndustries() ([]string, error) {
	var industries []string
	err := r.db.Model(&models.Job{}).
		Where("is_active = ?", true).
		Distinct("industry").
		Where("industry != ''").
		Pluck("industry", &industries).Error
	return industries, err
}

func (r *jobRepository) GetDistinctCities() ([]string, error) {
	var cities []string
	err := r.db.Model(&models.Job{}).
		Where("is_active = ?", true).
		Distinct("city").
		Where("city != ''").
		Pluck("city", &cities).Error
	return cities, err
}