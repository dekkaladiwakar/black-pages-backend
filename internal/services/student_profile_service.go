package services

import (
	"errors"

	"github.com/dekkaladiwakar/black-pages-backend/internal/models"
	"github.com/dekkaladiwakar/black-pages-backend/internal/repositories"
	"github.com/dekkaladiwakar/black-pages-backend/internal/utils"

	"gorm.io/gorm"
)

type CreateStudentProfileRequest struct {
	CollegeName         string   `json:"college_name" binding:"required"`
	Degree              string   `json:"degree" binding:"required"`
	YearSemester        string   `json:"year_semester" binding:"required"`
	SoftwareProficiency []string `json:"software_proficiency"`
	PreviousInternships []string `json:"previous_internships"`
	FreelanceProjects   []string `json:"freelance_projects"`
	PreferredStartMonth string   `json:"preferred_start_month"`
	PreferredDuration   string   `json:"preferred_duration"`
	WillingToRelocate   bool     `json:"willing_to_relocate"`
}

type UpdateStudentProfileRequest struct {
	CollegeName         string   `json:"college_name"`
	Degree              string   `json:"degree"`
	YearSemester        string   `json:"year_semester"`
	SoftwareProficiency []string `json:"software_proficiency"`
	PreviousInternships []string `json:"previous_internships"`
	FreelanceProjects   []string `json:"freelance_projects"`
	PreferredStartMonth string   `json:"preferred_start_month"`
	PreferredDuration   string   `json:"preferred_duration"`
	WillingToRelocate   *bool    `json:"willing_to_relocate"`
}

type StudentProfileService interface {
	CreateProfile(jobSeekerID uint, req CreateStudentProfileRequest) (*models.StudentProfile, error)
	GetProfile(jobSeekerID uint) (*models.StudentProfile, error)
	UpdateProfile(jobSeekerID uint, req UpdateStudentProfileRequest) (*models.StudentProfile, error)
	DeleteProfile(jobSeekerID uint) error
}

type studentProfileService struct {
	studentProfileRepo repositories.StudentProfileRepository
	jobSeekerRepo      repositories.JobSeekerRepository
}

func NewStudentProfileService(
	studentProfileRepo repositories.StudentProfileRepository,
	jobSeekerRepo repositories.JobSeekerRepository,
) StudentProfileService {
	return &studentProfileService{
		studentProfileRepo: studentProfileRepo,
		jobSeekerRepo:      jobSeekerRepo,
	}
}

func (s *studentProfileService) CreateProfile(jobSeekerID uint, req CreateStudentProfileRequest) (*models.StudentProfile, error) {
	// Verify job seeker exists
	jobSeeker, err := s.jobSeekerRepo.GetByID(jobSeekerID)
	if err != nil {
		return nil, errors.New("job seeker profile not found")
	}

	// Check if job seeker is actually a student
	if jobSeeker.JobSeekerType != "student" {
		return nil, errors.New("student profile can only be created for student job seekers")
	}

	// Check if student profile already exists
	existingProfile, err := s.studentProfileRepo.GetByJobSeekerID(jobSeekerID)
	if err == nil && existingProfile != nil {
		return nil, errors.New("student profile already exists")
	}

	// Convert arrays to JSON strings
	softwareJSON := utils.ArrayToJSON(req.SoftwareProficiency)
	internshipsJSON := utils.ArrayToJSON(req.PreviousInternships)
	projectsJSON := utils.ArrayToJSON(req.FreelanceProjects)

	// Create student profile
	profile := &models.StudentProfile{
		JobSeekerID:         jobSeekerID,
		CollegeName:         req.CollegeName,
		Degree:              req.Degree,
		YearSemester:        req.YearSemester,
		SoftwareProficiency: softwareJSON,
		PreviousInternships: internshipsJSON,
		FreelanceProjects:   projectsJSON,
		PreferredStartMonth: req.PreferredStartMonth,
		PreferredDuration:   req.PreferredDuration,
		WillingToRelocate:   req.WillingToRelocate,
	}

	if err := s.studentProfileRepo.Create(profile); err != nil {
		return nil, errors.New("failed to create student profile")
	}

	// Load relationships for response
	profile.JobSeeker = *jobSeeker
	return profile, nil
}

func (s *studentProfileService) GetProfile(jobSeekerID uint) (*models.StudentProfile, error) {
	// Verify job seeker exists
	_, err := s.jobSeekerRepo.GetByID(jobSeekerID)
	if err != nil {
		return nil, errors.New("job seeker profile not found")
	}

	profile, err := s.studentProfileRepo.GetByJobSeekerID(jobSeekerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("student profile not found")
		}
		return nil, err
	}
	return profile, nil
}

func (s *studentProfileService) UpdateProfile(jobSeekerID uint, req UpdateStudentProfileRequest) (*models.StudentProfile, error) {
	// Get existing profile
	profile, err := s.studentProfileRepo.GetByJobSeekerID(jobSeekerID)
	if err != nil {
		return nil, errors.New("student profile not found")
	}

	// Update fields if provided
	if req.CollegeName != "" {
		profile.CollegeName = req.CollegeName
	}
	if req.Degree != "" {
		profile.Degree = req.Degree
	}
	if req.YearSemester != "" {
		profile.YearSemester = req.YearSemester
	}
	if len(req.SoftwareProficiency) > 0 {
		profile.SoftwareProficiency = utils.ArrayToJSON(req.SoftwareProficiency)
	}
	if len(req.PreviousInternships) > 0 {
		profile.PreviousInternships = utils.ArrayToJSON(req.PreviousInternships)
	}
	if len(req.FreelanceProjects) > 0 {
		profile.FreelanceProjects = utils.ArrayToJSON(req.FreelanceProjects)
	}
	if req.PreferredStartMonth != "" {
		profile.PreferredStartMonth = req.PreferredStartMonth
	}
	if req.PreferredDuration != "" {
		profile.PreferredDuration = req.PreferredDuration
	}
	if req.WillingToRelocate != nil {
		profile.WillingToRelocate = *req.WillingToRelocate
	}

	if err := s.studentProfileRepo.Update(profile); err != nil {
		return nil, errors.New("failed to update student profile")
	}

	return profile, nil
}

func (s *studentProfileService) DeleteProfile(jobSeekerID uint) error {
	// Verify profile exists
	_, err := s.studentProfileRepo.GetByJobSeekerID(jobSeekerID)
	if err != nil {
		return errors.New("student profile not found")
	}

	return s.studentProfileRepo.Delete(jobSeekerID)
}

