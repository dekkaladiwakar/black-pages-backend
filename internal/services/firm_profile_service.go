package services

import (
	"errors"

	"github.com/dekkaladiwakar/black-pages-backend/internal/models"
	"github.com/dekkaladiwakar/black-pages-backend/internal/repositories"
	"github.com/dekkaladiwakar/black-pages-backend/internal/utils"

	"gorm.io/gorm"
)

type CreateFirmProfileRequest struct {
	YearFounded          int      `json:"year_founded"`
	FirmSize             string   `json:"firm_size"`
	LegalEntityType      string   `json:"legal_entity_type"`
	PrimaryDiscipline    string   `json:"primary_discipline" binding:"required"`
	SecondaryDisciplines []string `json:"secondary_disciplines"`
	InstagramURL         string   `json:"instagram_url"`
	LinkedInURL          string   `json:"linkedin_url"`
	PreferredDuration    string   `json:"preferred_duration"`
	StipendRange         string   `json:"stipend_range"`
	ProjectImages        []string `json:"project_images"`
}

type UpdateFirmProfileRequest struct {
	YearFounded          *int     `json:"year_founded"`
	FirmSize             string   `json:"firm_size"`
	LegalEntityType      string   `json:"legal_entity_type"`
	PrimaryDiscipline    string   `json:"primary_discipline"`
	SecondaryDisciplines []string `json:"secondary_disciplines"`
	InstagramURL         string   `json:"instagram_url"`
	LinkedInURL          string   `json:"linkedin_url"`
	PreferredDuration    string   `json:"preferred_duration"`
	StipendRange         string   `json:"stipend_range"`
	ProjectImages        []string `json:"project_images"`
}

type FirmProfileService interface {
	CreateProfile(employerID uint, req CreateFirmProfileRequest) (*models.FirmProfile, error)
	GetProfile(employerID uint) (*models.FirmProfile, error)
	UpdateProfile(employerID uint, req UpdateFirmProfileRequest) (*models.FirmProfile, error)
	DeleteProfile(employerID uint) error
}

type firmProfileService struct {
	firmProfileRepo repositories.FirmProfileRepository
	employerRepo    repositories.EmployerRepository
}

func NewFirmProfileService(
	firmProfileRepo repositories.FirmProfileRepository,
	employerRepo repositories.EmployerRepository,
) FirmProfileService {
	return &firmProfileService{
		firmProfileRepo: firmProfileRepo,
		employerRepo:    employerRepo,
	}
}

func (s *firmProfileService) CreateProfile(employerID uint, req CreateFirmProfileRequest) (*models.FirmProfile, error) {
	// Verify employer exists
	employer, err := s.employerRepo.GetByID(employerID)
	if err != nil {
		return nil, errors.New("employer profile not found")
	}

	// Check if employer is actually a firm
	if employer.EmployerType != "firm" {
		return nil, errors.New("firm profile can only be created for firm employers")
	}

	// Check if firm profile already exists
	existingProfile, err := s.firmProfileRepo.GetByEmployerID(employerID)
	if err == nil && existingProfile != nil {
		return nil, errors.New("firm profile already exists")
	}

	// Convert arrays to JSON strings
	disciplinesJSON := utils.ArrayToJSON(req.SecondaryDisciplines)
	imagesJSON := utils.ArrayToJSON(req.ProjectImages)

	// Create firm profile
	profile := &models.FirmProfile{
		EmployerID:           employerID,
		YearFounded:          req.YearFounded,
		FirmSize:             req.FirmSize,
		LegalEntityType:      req.LegalEntityType,
		PrimaryDiscipline:    req.PrimaryDiscipline,
		SecondaryDisciplines: disciplinesJSON,
		InstagramURL:         req.InstagramURL,
		LinkedInURL:          req.LinkedInURL,
		PreferredDuration:    req.PreferredDuration,
		StipendRange:         req.StipendRange,
		ProjectImages:        imagesJSON,
	}

	if err := s.firmProfileRepo.Create(profile); err != nil {
		return nil, errors.New("failed to create firm profile")
	}

	// Load relationships for response
	profile.Employer = *employer
	return profile, nil
}

func (s *firmProfileService) GetProfile(employerID uint) (*models.FirmProfile, error) {
	// Verify employer exists
	_, err := s.employerRepo.GetByID(employerID)
	if err != nil {
		return nil, errors.New("employer profile not found")
	}

	profile, err := s.firmProfileRepo.GetByEmployerID(employerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("firm profile not found")
		}
		return nil, err
	}
	return profile, nil
}

func (s *firmProfileService) UpdateProfile(employerID uint, req UpdateFirmProfileRequest) (*models.FirmProfile, error) {
	// Get existing profile
	profile, err := s.firmProfileRepo.GetByEmployerID(employerID)
	if err != nil {
		return nil, errors.New("firm profile not found")
	}

	// Update fields if provided
	if req.YearFounded != nil {
		profile.YearFounded = *req.YearFounded
	}
	if req.FirmSize != "" {
		profile.FirmSize = req.FirmSize
	}
	if req.LegalEntityType != "" {
		profile.LegalEntityType = req.LegalEntityType
	}
	if req.PrimaryDiscipline != "" {
		profile.PrimaryDiscipline = req.PrimaryDiscipline
	}
	if len(req.SecondaryDisciplines) > 0 {
		profile.SecondaryDisciplines = utils.ArrayToJSON(req.SecondaryDisciplines)
	}
	if req.InstagramURL != "" {
		profile.InstagramURL = req.InstagramURL
	}
	if req.LinkedInURL != "" {
		profile.LinkedInURL = req.LinkedInURL
	}
	if req.PreferredDuration != "" {
		profile.PreferredDuration = req.PreferredDuration
	}
	if req.StipendRange != "" {
		profile.StipendRange = req.StipendRange
	}
	if len(req.ProjectImages) > 0 {
		profile.ProjectImages = utils.ArrayToJSON(req.ProjectImages)
	}

	if err := s.firmProfileRepo.Update(profile); err != nil {
		return nil, errors.New("failed to update firm profile")
	}

	return profile, nil
}

func (s *firmProfileService) DeleteProfile(employerID uint) error {
	// Verify profile exists
	_, err := s.firmProfileRepo.GetByEmployerID(employerID)
	if err != nil {
		return errors.New("firm profile not found")
	}

	return s.firmProfileRepo.Delete(employerID)
}