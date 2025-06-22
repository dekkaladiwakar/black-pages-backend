package services

import (
	"errors"
	"time"

	"github.com/dekkaladiwakar/black-pages-backend/internal/models"
	"github.com/dekkaladiwakar/black-pages-backend/internal/repositories"

	"gorm.io/gorm"
)

type CreateJobRequest struct {
	Title               string    `json:"title" binding:"required"`
	JobType             string    `json:"job_type" binding:"required,oneof=internship full_time contract"`
	Industry            string    `json:"industry" binding:"required"`
	TargetAudience      string    `json:"target_audience" binding:"required,oneof=students professionals any"`
	EmploymentMode      string    `json:"employment_mode" binding:"required,oneof=on_site remote hybrid"`
	StartMonth          string    `json:"start_month" binding:"required"`
	Duration            string    `json:"duration" binding:"required"`
	ApplicationDeadline time.Time `json:"application_deadline" binding:"required"`
	CompensationRange   string    `json:"compensation_range"`
	IsPaid              bool      `json:"is_paid"`
	City                string    `json:"city" binding:"required"`
	State               string    `json:"state" binding:"required"`
	RequiredSkills      []string  `json:"required_skills" binding:"required,min=1"`
	MinExperience       string    `json:"min_experience"`
	PortfolioRequired   bool      `json:"portfolio_required"`
	ResumeRequired      bool      `json:"resume_required"`
	Description         string    `json:"description" binding:"required"`
	AboutTeam           string    `json:"about_team"`
	ContactEmail        string    `json:"contact_email" binding:"required,email"`
}

type UpdateJobRequest struct {
	Title               string    `json:"title"`
	JobType             string    `json:"job_type" binding:"omitempty,oneof=internship full_time contract"`
	Industry            string    `json:"industry"`
	TargetAudience      string    `json:"target_audience" binding:"omitempty,oneof=students professionals any"`
	EmploymentMode      string    `json:"employment_mode" binding:"omitempty,oneof=on_site remote hybrid"`
	StartMonth          string    `json:"start_month"`
	Duration            string    `json:"duration"`
	ApplicationDeadline *time.Time `json:"application_deadline"`
	CompensationRange   string    `json:"compensation_range"`
	IsPaid              *bool     `json:"is_paid"`
	City                string    `json:"city"`
	State               string    `json:"state"`
	RequiredSkills      []string  `json:"required_skills"`
	MinExperience       string    `json:"min_experience"`
	PortfolioRequired   *bool     `json:"portfolio_required"`
	ResumeRequired      *bool     `json:"resume_required"`
	Description         string    `json:"description"`
	AboutTeam           string    `json:"about_team"`
	ContactEmail        string    `json:"contact_email" binding:"omitempty,email"`
	IsActive            *bool     `json:"is_active"`
}

type JobFilters struct {
	Industry        string `form:"industry"`
	JobType         string `form:"job_type"`
	City            string `form:"city"`
	TargetAudience  string `form:"target_audience"`
	EmploymentMode  string `form:"employment_mode"`
	IsPaid          *bool  `form:"is_paid"`
	IsActive        *bool  `form:"is_active"`
	OrderBy         string `form:"order_by"`
	OrderDirection  string `form:"order_direction"`
	Limit           int    `form:"limit"`
}

type FilterOptions struct {
	Industries       []string `json:"industries"`
	JobTypes         []string `json:"job_types"`
	TargetAudiences  []string `json:"target_audiences"`
	EmploymentModes  []string `json:"employment_modes"`
	Cities           []string `json:"cities"`
}

type JobService interface {
	CreateJob(employerID uint, req CreateJobRequest) (*models.Job, error)
	GetJob(id uint) (*models.Job, error)
	UpdateJob(employerID uint, jobID uint, req UpdateJobRequest) (*models.Job, error)
	DeleteJob(employerID uint, jobID uint) error
	GetEmployerJobs(employerID uint, filters JobFilters) ([]models.Job, error)
	GetAllJobs(filters JobFilters) ([]models.Job, error)
	ToggleJobStatus(employerID uint, jobID uint) (*models.Job, error)
	GetEmployerDashboardStats(employerID uint) (map[string]interface{}, error)
	GetFilterOptions() (*FilterOptions, error)
}

type jobService struct {
	jobRepo      repositories.JobRepository
	employerRepo repositories.EmployerRepository
}

func NewJobService(jobRepo repositories.JobRepository, employerRepo repositories.EmployerRepository) JobService {
	return &jobService{
		jobRepo:      jobRepo,
		employerRepo: employerRepo,
	}
}

func (s *jobService) CreateJob(employerID uint, req CreateJobRequest) (*models.Job, error) {
	employer, err := s.employerRepo.GetByID(employerID)
	if err != nil {
		return nil, errors.New("employer not found")
	}

	if req.ApplicationDeadline.Before(time.Now()) {
		return nil, errors.New("application deadline cannot be in the past")
	}

	skillsJSON := ""
	if len(req.RequiredSkills) > 0 {
		skillsJSON = `["` + req.RequiredSkills[0]
		for i := 1; i < len(req.RequiredSkills); i++ {
			skillsJSON += `","` + req.RequiredSkills[i]
		}
		skillsJSON += `"]`
	}

	job := &models.Job{
		EmployerID:          employerID,
		Title:               req.Title,
		JobType:             req.JobType,
		Industry:            req.Industry,
		TargetAudience:      req.TargetAudience,
		EmploymentMode:      req.EmploymentMode,
		StartMonth:          req.StartMonth,
		Duration:            req.Duration,
		ApplicationDeadline: req.ApplicationDeadline,
		CompensationRange:   req.CompensationRange,
		IsPaid:              req.IsPaid,
		City:                req.City,
		State:               req.State,
		RequiredSkills:      skillsJSON,
		MinExperience:       req.MinExperience,
		PortfolioRequired:   req.PortfolioRequired,
		ResumeRequired:      req.ResumeRequired,
		Description:         req.Description,
		AboutTeam:           req.AboutTeam,
		ContactEmail:        req.ContactEmail,
		IsActive:            true,
	}

	if err := s.jobRepo.Create(job); err != nil {
		return nil, errors.New("failed to create job")
	}

	job.Employer = *employer
	return job, nil
}

func (s *jobService) GetJob(id uint) (*models.Job, error) {
	job, err := s.jobRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("job not found")
		}
		return nil, err
	}
	return job, nil
}

func (s *jobService) UpdateJob(employerID uint, jobID uint, req UpdateJobRequest) (*models.Job, error) {
	job, err := s.jobRepo.GetByID(jobID)
	if err != nil {
		return nil, errors.New("job not found")
	}

	if job.EmployerID != employerID {
		return nil, errors.New("unauthorized to update this job")
	}

	if req.Title != "" {
		job.Title = req.Title
	}
	if req.JobType != "" {
		job.JobType = req.JobType
	}
	if req.Industry != "" {
		job.Industry = req.Industry
	}
	if req.TargetAudience != "" {
		job.TargetAudience = req.TargetAudience
	}
	if req.EmploymentMode != "" {
		job.EmploymentMode = req.EmploymentMode
	}
	if req.StartMonth != "" {
		job.StartMonth = req.StartMonth
	}
	if req.Duration != "" {
		job.Duration = req.Duration
	}
	if req.ApplicationDeadline != nil {
		if req.ApplicationDeadline.Before(time.Now()) {
			return nil, errors.New("application deadline cannot be in the past")
		}
		job.ApplicationDeadline = *req.ApplicationDeadline
	}
	if req.CompensationRange != "" {
		job.CompensationRange = req.CompensationRange
	}
	if req.IsPaid != nil {
		job.IsPaid = *req.IsPaid
	}
	if req.City != "" {
		job.City = req.City
	}
	if req.State != "" {
		job.State = req.State
	}
	if len(req.RequiredSkills) > 0 {
		skillsJSON := `["` + req.RequiredSkills[0]
		for i := 1; i < len(req.RequiredSkills); i++ {
			skillsJSON += `","` + req.RequiredSkills[i]
		}
		skillsJSON += `"]`
		job.RequiredSkills = skillsJSON
	}
	if req.MinExperience != "" {
		job.MinExperience = req.MinExperience
	}
	if req.PortfolioRequired != nil {
		job.PortfolioRequired = *req.PortfolioRequired
	}
	if req.ResumeRequired != nil {
		job.ResumeRequired = *req.ResumeRequired
	}
	if req.Description != "" {
		job.Description = req.Description
	}
	if req.AboutTeam != "" {
		job.AboutTeam = req.AboutTeam
	}
	if req.ContactEmail != "" {
		job.ContactEmail = req.ContactEmail
	}
	if req.IsActive != nil {
		job.IsActive = *req.IsActive
	}

	if err := s.jobRepo.Update(job); err != nil {
		return nil, errors.New("failed to update job")
	}

	return job, nil
}

func (s *jobService) DeleteJob(employerID uint, jobID uint) error {
	job, err := s.jobRepo.GetByID(jobID)
	if err != nil {
		return errors.New("job not found")
	}

	if job.EmployerID != employerID {
		return errors.New("unauthorized to delete this job")
	}

	return s.jobRepo.Delete(jobID)
}

func (s *jobService) GetEmployerJobs(employerID uint, filters JobFilters) ([]models.Job, error) {
	repoFilters := repositories.JobFilters{
		EmployerID:     employerID,
		Industry:       filters.Industry,
		JobType:        filters.JobType,
		City:           filters.City,
		TargetAudience: filters.TargetAudience,
		EmploymentMode: filters.EmploymentMode,
		IsPaid:         filters.IsPaid,
		IsActive:       filters.IsActive,
		OrderBy:        filters.OrderBy,
		OrderDirection: filters.OrderDirection,
		Limit:          filters.Limit,
	}

	return s.jobRepo.GetWithFilters(repoFilters)
}

func (s *jobService) GetAllJobs(filters JobFilters) ([]models.Job, error) {
	// For public job browsing, always filter to only active jobs
	activeOnly := true
	
	repoFilters := repositories.JobFilters{
		Industry:       filters.Industry,
		JobType:        filters.JobType,
		City:           filters.City,
		TargetAudience: filters.TargetAudience,
		EmploymentMode: filters.EmploymentMode,
		IsPaid:         filters.IsPaid,
		IsActive:       &activeOnly,  // Force active jobs only for public browsing
		OrderBy:        filters.OrderBy,
		OrderDirection: filters.OrderDirection,
		Limit:          filters.Limit,
	}

	return s.jobRepo.GetWithFilters(repoFilters)
}

func (s *jobService) ToggleJobStatus(employerID uint, jobID uint) (*models.Job, error) {
	job, err := s.jobRepo.GetByID(jobID)
	if err != nil {
		return nil, errors.New("job not found")
	}

	if job.EmployerID != employerID {
		return nil, errors.New("unauthorized to modify this job")
	}

	job.IsActive = !job.IsActive

	if err := s.jobRepo.Update(job); err != nil {
		return nil, errors.New("failed to update job status")
	}

	return job, nil
}

func (s *jobService) GetEmployerDashboardStats(employerID uint) (map[string]interface{}, error) {
	totalJobs, err := s.jobRepo.CountByEmployerID(employerID)
	if err != nil {
		return nil, err
	}

	jobs, err := s.jobRepo.GetByEmployerID(employerID)
	if err != nil {
		return nil, err
	}

	activeJobs := int64(0)
	for _, job := range jobs {
		if job.IsActive {
			activeJobs++
		}
	}

	stats := map[string]interface{}{
		"total_jobs":  totalJobs,
		"active_jobs": activeJobs,
		"recent_jobs": jobs[:min(5, len(jobs))], // Last 5 jobs
	}

	return stats, nil
}

func (s *jobService) GetFilterOptions() (*FilterOptions, error) {
	// Query distinct values from database
	industries, err := s.jobRepo.GetDistinctIndustries()
	if err != nil {
		return nil, err
	}

	cities, err := s.jobRepo.GetDistinctCities()
	if err != nil {
		return nil, err
	}

	// Static options for enum fields
	jobTypes := []string{"internship", "full_time", "contract"}
	targetAudiences := []string{"students", "professionals", "any"}
	employmentModes := []string{"on_site", "remote", "hybrid"}
	
	// Static industries with fallback for empty database
	staticIndustries := []string{"Architecture", "Interior Design", "Urban Planning", "Construction", "Landscape Architecture"}
	if len(industries) == 0 {
		industries = staticIndustries
	}

	return &FilterOptions{
		Industries:      industries,
		JobTypes:        jobTypes,
		TargetAudiences: targetAudiences,
		EmploymentModes: employmentModes,
		Cities:          cities,
	}, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}