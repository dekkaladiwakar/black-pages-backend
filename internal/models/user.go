package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
	PasswordHash string    `gorm:"not null" json:"-"` // Don't include in JSON responses
	UserType     string    `gorm:"not null" json:"user_type" validate:"required,oneof=job_seeker employer"`
	IsVerified   bool      `gorm:"default:false" json:"is_verified"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	
	// Relationships
	JobSeeker *JobSeeker `gorm:"foreignKey:UserID" json:"job_seeker,omitempty"`
	Employer  *Employer  `gorm:"foreignKey:UserID" json:"employer,omitempty"`
}

type JobSeeker struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	User          User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	FullName      string    `gorm:"not null" json:"full_name" validate:"required"`
	JobSeekerType string    `gorm:"not null" json:"job_seeker_type" validate:"required,oneof=student professional freelancer"`
	CurrentCity   string    `gorm:"not null" json:"current_city" validate:"required"`
	Phone         string    `gorm:"not null" json:"phone" validate:"required"`
	DesiredField  string    `gorm:"not null" json:"desired_field" validate:"required"`
	ResumeURL     string    `gorm:"not null" json:"resume_url" validate:"required"`
	PortfolioURL  string    `json:"portfolio_url"`
	Skills        string    `gorm:"type:json" json:"skills"` // JSON array as string
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	
	// Relationships
	StudentProfile *StudentProfile `gorm:"foreignKey:JobSeekerID" json:"student_profile,omitempty"`
	Applications   []Application   `gorm:"foreignKey:JobSeekerID" json:"applications,omitempty"`
}

type Employer struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	UserID             uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	User               User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CompanyName        string    `gorm:"not null" json:"company_name" validate:"required"`
	EmployerType       string    `gorm:"not null" json:"employer_type" validate:"required,oneof=firm corporation startup"`
	Industry           string    `gorm:"not null" json:"industry" validate:"required"`
	PrimaryPhone       string    `gorm:"not null" json:"primary_phone" validate:"required"`
	ContactPerson      string    `gorm:"not null" json:"contact_person" validate:"required"`
	ContactPersonDesig string    `gorm:"not null" json:"contact_person_desig" validate:"required"`
	City               string    `gorm:"not null" json:"city" validate:"required"`
	State              string    `gorm:"not null" json:"state" validate:"required"`
	PinCode            string    `gorm:"size:6;not null" json:"pin_code" validate:"required,len=6,numeric"`
	WebsiteURL         string    `gorm:"not null" json:"website_url" validate:"required,url"`
	LogoURL            string    `json:"logo_url"`
	IsHiring           bool      `gorm:"default:false" json:"is_hiring"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	
	// Relationships
	FirmProfile *FirmProfile `gorm:"foreignKey:EmployerID" json:"firm_profile,omitempty"`
	Jobs        []Job        `gorm:"foreignKey:EmployerID" json:"jobs,omitempty"`
}

// BeforeCreate hooks for validation
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Additional validation logic can be added here
	return nil
}

func (js *JobSeeker) BeforeCreate(tx *gorm.DB) error {
	// Additional validation logic can be added here
	return nil
}