package services

import (
	"errors"
	"time"

	"github.com/dekkaladiwakar/black-pages-backend/internal/models"
	"github.com/dekkaladiwakar/black-pages-backend/internal/repositories"

	"gorm.io/gorm"
)

type ApplyJobRequest struct {
	JobID uint `json:"job_id" binding:"required"`
}

type UpdateApplicationStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=applied shortlisted rejected selected"`
}

type ApplicationService interface {
	ApplyToJob(jobSeekerID uint, req ApplyJobRequest) (*models.Application, error)
	GetJobSeekerApplications(jobSeekerID uint) ([]models.Application, error)
	GetJobApplications(jobID uint, employerID uint) ([]models.Application, error)
	UpdateApplicationStatus(applicationID uint, employerID uint, req UpdateApplicationStatusRequest) (*models.Application, error)
	GetApplication(applicationID uint) (*models.Application, error)
	WithdrawApplication(applicationID uint, jobSeekerID uint) error
	GetApplicationStats(jobSeekerID uint) (map[string]interface{}, error)
	GetJobApplicationStats(jobID uint, employerID uint) (map[string]interface{}, error)
}

type applicationService struct {
	applicationRepo repositories.ApplicationRepository
	jobRepo         repositories.JobRepository
	jobSeekerRepo   repositories.JobSeekerRepository
	employerRepo    repositories.EmployerRepository
}

func NewApplicationService(
	applicationRepo repositories.ApplicationRepository,
	jobRepo repositories.JobRepository,
	jobSeekerRepo repositories.JobSeekerRepository,
	employerRepo repositories.EmployerRepository,
) ApplicationService {
	return &applicationService{
		applicationRepo: applicationRepo,
		jobRepo:         jobRepo,
		jobSeekerRepo:   jobSeekerRepo,
		employerRepo:    employerRepo,
	}
}

func (s *applicationService) ApplyToJob(jobSeekerID uint, req ApplyJobRequest) (*models.Application, error) {
	// Check if job seeker exists
	jobSeeker, err := s.jobSeekerRepo.GetByID(jobSeekerID)
	if err != nil {
		return nil, errors.New("job seeker profile not found")
	}

	// Check if job exists and is active
	job, err := s.jobRepo.GetByID(req.JobID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("job not found")
		}
		return nil, err
	}

	if !job.IsActive {
		return nil, errors.New("job is no longer active")
	}

	// Check if application deadline has passed
	if job.ApplicationDeadline.Before(time.Now()) {
		return nil, errors.New("application deadline has passed")
	}

	// Check if user has already applied
	existingApp, err := s.applicationRepo.GetByJobAndJobSeeker(req.JobID, jobSeekerID)
	if err == nil && existingApp != nil {
		return nil, errors.New("you have already applied to this job")
	}

	// Validate job requirements
	if job.ResumeRequired && jobSeeker.ResumeURL == "" {
		return nil, errors.New("resume is required for this job")
	}

	if job.PortfolioRequired && jobSeeker.PortfolioURL == "" {
		return nil, errors.New("portfolio is required for this job")
	}

	// Create application
	application := &models.Application{
		JobID:       req.JobID,
		JobSeekerID: jobSeekerID,
		Status:      "applied",
		AppliedAt:   time.Now(),
	}

	if err := s.applicationRepo.Create(application); err != nil {
		return nil, errors.New("failed to submit application")
	}

	// Load relationships for response
	application.Job = *job
	application.JobSeeker = *jobSeeker

	return application, nil
}

func (s *applicationService) GetJobSeekerApplications(jobSeekerID uint) ([]models.Application, error) {
	// Verify job seeker exists
	_, err := s.jobSeekerRepo.GetByID(jobSeekerID)
	if err != nil {
		return nil, errors.New("job seeker profile not found")
	}

	return s.applicationRepo.GetByJobSeekerID(jobSeekerID)
}

func (s *applicationService) GetJobApplications(jobID uint, employerID uint) ([]models.Application, error) {
	// Verify job exists and belongs to employer
	job, err := s.jobRepo.GetByID(jobID)
	if err != nil {
		return nil, errors.New("job not found")
	}

	if job.EmployerID != employerID {
		return nil, errors.New("unauthorized to view applications for this job")
	}

	return s.applicationRepo.GetByJobID(jobID)
}

func (s *applicationService) UpdateApplicationStatus(applicationID uint, employerID uint, req UpdateApplicationStatusRequest) (*models.Application, error) {
	// Get application with job info
	application, err := s.applicationRepo.GetByID(applicationID)
	if err != nil {
		return nil, errors.New("application not found")
	}

	// Verify employer owns the job
	if application.Job.EmployerID != employerID {
		return nil, errors.New("unauthorized to update this application")
	}

	// Update status
	application.Status = req.Status
	application.UpdatedAt = time.Now()

	if err := s.applicationRepo.Update(application); err != nil {
		return nil, errors.New("failed to update application status")
	}

	return application, nil
}

func (s *applicationService) GetApplication(applicationID uint) (*models.Application, error) {
	application, err := s.applicationRepo.GetByID(applicationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("application not found")
		}
		return nil, err
	}
	return application, nil
}

func (s *applicationService) WithdrawApplication(applicationID uint, jobSeekerID uint) error {
	// Get application
	application, err := s.applicationRepo.GetByID(applicationID)
	if err != nil {
		return errors.New("application not found")
	}

	// Verify job seeker owns the application
	if application.JobSeekerID != jobSeekerID {
		return errors.New("unauthorized to withdraw this application")
	}

	// Only allow withdrawal if not already processed
	if application.Status == "selected" || application.Status == "rejected" {
		return errors.New("cannot withdraw application that has been processed")
	}

	return s.applicationRepo.Delete(applicationID)
}

func (s *applicationService) GetApplicationStats(jobSeekerID uint) (map[string]interface{}, error) {
	// Verify job seeker exists
	_, err := s.jobSeekerRepo.GetByID(jobSeekerID)
	if err != nil {
		return nil, errors.New("job seeker profile not found")
	}

	// Get all applications
	applications, err := s.applicationRepo.GetByJobSeekerID(jobSeekerID)
	if err != nil {
		return nil, err
	}

	// Calculate stats
	stats := map[string]interface{}{
		"total_applications": len(applications),
		"applied":           0,
		"shortlisted":       0,
		"rejected":          0,
		"selected":          0,
		"recent_applications": applications[:min(5, len(applications))],
	}

	for _, app := range applications {
		switch app.Status {
		case "applied":
			stats["applied"] = stats["applied"].(int) + 1
		case "shortlisted":
			stats["shortlisted"] = stats["shortlisted"].(int) + 1
		case "rejected":
			stats["rejected"] = stats["rejected"].(int) + 1
		case "selected":
			stats["selected"] = stats["selected"].(int) + 1
		}
	}

	return stats, nil
}

func (s *applicationService) GetJobApplicationStats(jobID uint, employerID uint) (map[string]interface{}, error) {
	// Verify job exists and belongs to employer
	job, err := s.jobRepo.GetByID(jobID)
	if err != nil {
		return nil, errors.New("job not found")
	}

	if job.EmployerID != employerID {
		return nil, errors.New("unauthorized to view stats for this job")
	}

	// Get all applications for the job
	applications, err := s.applicationRepo.GetByJobID(jobID)
	if err != nil {
		return nil, err
	}

	// Calculate stats
	stats := map[string]interface{}{
		"total_applications": len(applications),
		"applied":           0,
		"shortlisted":       0,
		"rejected":          0,
		"selected":          0,
		"job_title":         job.Title,
		"application_deadline": job.ApplicationDeadline,
	}

	for _, app := range applications {
		switch app.Status {
		case "applied":
			stats["applied"] = stats["applied"].(int) + 1
		case "shortlisted":
			stats["shortlisted"] = stats["shortlisted"].(int) + 1
		case "rejected":
			stats["rejected"] = stats["rejected"].(int) + 1
		case "selected":
			stats["selected"] = stats["selected"].(int) + 1
		}
	}

	return stats, nil
}

