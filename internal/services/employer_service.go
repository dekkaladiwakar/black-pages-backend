package services

import (
	"errors"

	"github.com/dekkaladiwakar/black-pages-backend/internal/models"
	"github.com/dekkaladiwakar/black-pages-backend/internal/repositories"

	"gorm.io/gorm"
)

type CreateEmployerRequest struct {
	CompanyName        string `json:"company_name" binding:"required"`
	EmployerType       string `json:"employer_type" binding:"required,oneof=firm corporation startup"`
	Industry           string `json:"industry" binding:"required"`
	PrimaryPhone       string `json:"primary_phone" binding:"required"`
	ContactPerson      string `json:"contact_person" binding:"required"`
	ContactPersonDesig string `json:"contact_person_desig" binding:"required"`
	City               string `json:"city" binding:"required"`
	State              string `json:"state" binding:"required"`
	PinCode            string `json:"pin_code" binding:"required,len=6"`
	WebsiteURL         string `json:"website_url" binding:"required,url"`
	LogoURL            string `json:"logo_url"`
}

type UpdateEmployerRequest struct {
	CompanyName        string `json:"company_name"`
	EmployerType       string `json:"employer_type" binding:"omitempty,oneof=firm corporation startup"`
	Industry           string `json:"industry"`
	PrimaryPhone       string `json:"primary_phone"`
	ContactPerson      string `json:"contact_person"`
	ContactPersonDesig string `json:"contact_person_desig"`
	City               string `json:"city"`
	State              string `json:"state"`
	PinCode            string `json:"pin_code" binding:"omitempty,len=6"`
	WebsiteURL         string `json:"website_url" binding:"omitempty,url"`
	LogoURL            string `json:"logo_url"`
	IsHiring           *bool  `json:"is_hiring"`
}

type EmployerService interface {
	CreateProfile(userID uint, req CreateEmployerRequest) (*models.Employer, error)
	GetProfile(userID uint) (*models.Employer, error)
	UpdateProfile(userID uint, req UpdateEmployerRequest) (*models.Employer, error)
	GetProfileWithExtensions(userID uint) (*models.Employer, error)
}

type employerService struct {
	employerRepo repositories.EmployerRepository
	userRepo     repositories.UserRepository
}

func NewEmployerService(employerRepo repositories.EmployerRepository, userRepo repositories.UserRepository) EmployerService {
	return &employerService{
		employerRepo: employerRepo,
		userRepo:     userRepo,
	}
}

func (s *employerService) CreateProfile(userID uint, req CreateEmployerRequest) (*models.Employer, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.UserType != "employer" {
		return nil, errors.New("user is not an employer")
	}

	existing, _ := s.employerRepo.GetByUserID(userID)
	if existing != nil {
		return nil, errors.New("employer profile already exists")
	}

	employer := &models.Employer{
		UserID:             userID,
		CompanyName:        req.CompanyName,
		EmployerType:       req.EmployerType,
		Industry:           req.Industry,
		PrimaryPhone:       req.PrimaryPhone,
		ContactPerson:      req.ContactPerson,
		ContactPersonDesig: req.ContactPersonDesig,
		City:               req.City,
		State:              req.State,
		PinCode:            req.PinCode,
		WebsiteURL:         req.WebsiteURL,
		LogoURL:            req.LogoURL,
		IsHiring:           false,
	}

	if err := s.employerRepo.Create(employer); err != nil {
		return nil, errors.New("failed to create employer profile")
	}

	return employer, nil
}

func (s *employerService) GetProfile(userID uint) (*models.Employer, error) {
	employer, err := s.employerRepo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("employer profile not found")
		}
		return nil, err
	}
	return employer, nil
}

func (s *employerService) UpdateProfile(userID uint, req UpdateEmployerRequest) (*models.Employer, error) {
	employer, err := s.employerRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("employer profile not found")
	}

	if req.CompanyName != "" {
		employer.CompanyName = req.CompanyName
	}
	if req.EmployerType != "" {
		employer.EmployerType = req.EmployerType
	}
	if req.Industry != "" {
		employer.Industry = req.Industry
	}
	if req.PrimaryPhone != "" {
		employer.PrimaryPhone = req.PrimaryPhone
	}
	if req.ContactPerson != "" {
		employer.ContactPerson = req.ContactPerson
	}
	if req.ContactPersonDesig != "" {
		employer.ContactPersonDesig = req.ContactPersonDesig
	}
	if req.City != "" {
		employer.City = req.City
	}
	if req.State != "" {
		employer.State = req.State
	}
	if req.PinCode != "" {
		employer.PinCode = req.PinCode
	}
	if req.WebsiteURL != "" {
		employer.WebsiteURL = req.WebsiteURL
	}
	if req.LogoURL != "" {
		employer.LogoURL = req.LogoURL
	}
	if req.IsHiring != nil {
		employer.IsHiring = *req.IsHiring
	}

	if err := s.employerRepo.Update(employer); err != nil {
		return nil, errors.New("failed to update employer profile")
	}

	return employer, nil
}

func (s *employerService) GetProfileWithExtensions(userID uint) (*models.Employer, error) {
	return s.employerRepo.GetWithFirmProfile(userID)
}