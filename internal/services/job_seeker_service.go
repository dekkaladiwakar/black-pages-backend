package services

import (
	"errors"

	"github.com/dekkaladiwakar/black-pages-backend/internal/models"
	"github.com/dekkaladiwakar/black-pages-backend/internal/repositories"

	"gorm.io/gorm"
)

type CreateJobSeekerRequest struct {
	FullName       string   `json:"full_name" binding:"required"`
	JobSeekerType  string   `json:"job_seeker_type" binding:"required,oneof=student professional freelancer"`
	CurrentCity    string   `json:"current_city" binding:"required"`
	Phone          string   `json:"phone" binding:"required"`
	DesiredField   string   `json:"desired_field" binding:"required"`
	ResumeURL      string   `json:"resume_url" binding:"required"`
	PortfolioURL   string   `json:"portfolio_url"`
	Skills         []string `json:"skills"`
}

type UpdateJobSeekerRequest struct {
	FullName       string   `json:"full_name"`
	JobSeekerType  string   `json:"job_seeker_type" binding:"omitempty,oneof=student professional freelancer"`
	CurrentCity    string   `json:"current_city"`
	Phone          string   `json:"phone"`
	DesiredField   string   `json:"desired_field"`
	ResumeURL      string   `json:"resume_url"`
	PortfolioURL   string   `json:"portfolio_url"`
	Skills         []string `json:"skills"`
}

type JobSeekerService interface {
	CreateProfile(userID uint, req CreateJobSeekerRequest) (*models.JobSeeker, error)
	GetProfile(userID uint) (*models.JobSeeker, error)
	UpdateProfile(userID uint, req UpdateJobSeekerRequest) (*models.JobSeeker, error)
	GetProfileWithExtensions(userID uint) (*models.JobSeeker, error)
}

type jobSeekerService struct {
	jobSeekerRepo repositories.JobSeekerRepository
	userRepo      repositories.UserRepository
}

func NewJobSeekerService(jobSeekerRepo repositories.JobSeekerRepository, userRepo repositories.UserRepository) JobSeekerService {
	return &jobSeekerService{
		jobSeekerRepo: jobSeekerRepo,
		userRepo:      userRepo,
	}
}

func (s *jobSeekerService) CreateProfile(userID uint, req CreateJobSeekerRequest) (*models.JobSeeker, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.UserType != "job_seeker" {
		return nil, errors.New("user is not a job seeker")
	}

	existing, _ := s.jobSeekerRepo.GetByUserID(userID)
	if existing != nil {
		return nil, errors.New("job seeker profile already exists")
	}

	skillsJSON := ""
	if len(req.Skills) > 0 {
		// Convert skills array to JSON string - simplified for MVP
		skillsJSON = `["` + req.Skills[0]
		for i := 1; i < len(req.Skills); i++ {
			skillsJSON += `","` + req.Skills[i]
		}
		skillsJSON += `"]`
	}

	jobSeeker := &models.JobSeeker{
		UserID:        userID,
		FullName:      req.FullName,
		JobSeekerType: req.JobSeekerType,
		CurrentCity:   req.CurrentCity,
		Phone:         req.Phone,
		DesiredField:  req.DesiredField,
		ResumeURL:     req.ResumeURL,
		PortfolioURL:  req.PortfolioURL,
		Skills:        skillsJSON,
	}

	if err := s.jobSeekerRepo.Create(jobSeeker); err != nil {
		return nil, errors.New("failed to create job seeker profile")
	}

	return jobSeeker, nil
}

func (s *jobSeekerService) GetProfile(userID uint) (*models.JobSeeker, error) {
	jobSeeker, err := s.jobSeekerRepo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("job seeker profile not found")
		}
		return nil, err
	}
	return jobSeeker, nil
}

func (s *jobSeekerService) UpdateProfile(userID uint, req UpdateJobSeekerRequest) (*models.JobSeeker, error) {
	jobSeeker, err := s.jobSeekerRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("job seeker profile not found")
	}

	if req.FullName != "" {
		jobSeeker.FullName = req.FullName
	}
	if req.JobSeekerType != "" {
		jobSeeker.JobSeekerType = req.JobSeekerType
	}
	if req.CurrentCity != "" {
		jobSeeker.CurrentCity = req.CurrentCity
	}
	if req.Phone != "" {
		jobSeeker.Phone = req.Phone
	}
	if req.DesiredField != "" {
		jobSeeker.DesiredField = req.DesiredField
	}
	if req.ResumeURL != "" {
		jobSeeker.ResumeURL = req.ResumeURL
	}
	if req.PortfolioURL != "" {
		jobSeeker.PortfolioURL = req.PortfolioURL
	}
	if len(req.Skills) > 0 {
		skillsJSON := `["` + req.Skills[0]
		for i := 1; i < len(req.Skills); i++ {
			skillsJSON += `","` + req.Skills[i]
		}
		skillsJSON += `"]`
		jobSeeker.Skills = skillsJSON
	}

	if err := s.jobSeekerRepo.Update(jobSeeker); err != nil {
		return nil, errors.New("failed to update job seeker profile")
	}

	return jobSeeker, nil
}

func (s *jobSeekerService) GetProfileWithExtensions(userID uint) (*models.JobSeeker, error) {
	return s.jobSeekerRepo.GetWithStudentProfile(userID)
}