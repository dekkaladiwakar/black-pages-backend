package models

import (
	"time"
)

// Extension models for specific use cases

type StudentProfile struct {
	ID                    uint      `gorm:"primaryKey" json:"id"`
	JobSeekerID           uint      `gorm:"uniqueIndex;not null" json:"job_seeker_id"`
	JobSeeker             JobSeeker `gorm:"foreignKey:JobSeekerID" json:"job_seeker,omitempty"`
	CollegeName           string    `gorm:"not null" json:"college_name" validate:"required"`
	Degree                string    `gorm:"not null" json:"degree" validate:"required"`
	YearSemester          string    `gorm:"not null" json:"year_semester" validate:"required"`
	SoftwareProficiency   string    `gorm:"type:json" json:"software_proficiency"`   // JSON array
	PreviousInternships   string    `gorm:"type:json" json:"previous_internships"`   // JSON array
	FreelanceProjects     string    `gorm:"type:json" json:"freelance_projects"`     // JSON array
	PreferredStartMonth   string    `json:"preferred_start_month"`
	PreferredDuration     string    `json:"preferred_duration"`
	WillingToRelocate     bool      `gorm:"default:false" json:"willing_to_relocate"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type FirmProfile struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	EmployerID           uint      `gorm:"uniqueIndex;not null" json:"employer_id"`
	Employer             Employer  `gorm:"foreignKey:EmployerID" json:"employer,omitempty"`
	YearFounded          int       `json:"year_founded"`
	FirmSize             string    `json:"firm_size"`
	LegalEntityType      string    `json:"legal_entity_type"`
	PrimaryDiscipline    string    `gorm:"not null" json:"primary_discipline" validate:"required"`
	SecondaryDisciplines string    `gorm:"type:json" json:"secondary_disciplines"` // JSON array
	InstagramURL         string    `json:"instagram_url"`
	LinkedInURL          string    `json:"linkedin_url"`
	PreferredDuration    string    `json:"preferred_duration"`
	StipendRange         string    `json:"stipend_range"`
	ProjectImages        string    `gorm:"type:json" json:"project_images"` // JSON array
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}