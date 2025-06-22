package models

import (
	"time"
)

type Job struct {
	ID                  uint            `gorm:"primaryKey" json:"id"`
	EmployerID          uint            `gorm:"not null" json:"employer_id"`
	Employer            Employer        `gorm:"foreignKey:EmployerID" json:"employer,omitempty"`
	Title               string          `gorm:"not null" json:"title" validate:"required"`
	JobType             string          `gorm:"not null" json:"job_type" validate:"required,oneof=internship full_time contract"`
	Industry            string          `gorm:"not null" json:"industry" validate:"required"`
	TargetAudience      string          `gorm:"not null" json:"target_audience" validate:"required,oneof=students professionals any"`
	EmploymentMode      string          `gorm:"not null" json:"employment_mode" validate:"required,oneof=on_site remote hybrid"`
	StartMonth          string          `gorm:"not null" json:"start_month" validate:"required"`
	Duration            string          `gorm:"not null" json:"duration" validate:"required"`
	ApplicationDeadline time.Time       `gorm:"not null" json:"application_deadline" validate:"required"`
	CompensationRange   string          `json:"compensation_range"`
	IsPaid              bool            `gorm:"not null" json:"is_paid"`
	City                string          `gorm:"not null" json:"city" validate:"required"`
	State               string          `gorm:"not null" json:"state" validate:"required"`
	RequiredSkills      string          `gorm:"type:json;not null" json:"required_skills" validate:"required"`
	MinExperience       string          `json:"min_experience"`
	PortfolioRequired   bool            `gorm:"not null" json:"portfolio_required"`
	ResumeRequired      bool            `gorm:"default:true" json:"resume_required"`
	Description         string          `gorm:"type:text;not null" json:"description" validate:"required"`
	AboutTeam           string          `gorm:"type:text" json:"about_team"`
	ContactEmail        string          `gorm:"not null" json:"contact_email" validate:"required,email"`
	IsActive            bool            `gorm:"default:true" json:"is_active"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
	
	// Relationships
	Applications []Application `gorm:"foreignKey:JobID" json:"applications,omitempty"`
}