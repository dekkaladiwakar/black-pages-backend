package models

import (
	"time"
)

type Application struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	JobID       uint      `gorm:"not null" json:"job_id"`
	Job         Job       `gorm:"foreignKey:JobID" json:"job,omitempty"`
	JobSeekerID uint      `gorm:"not null" json:"job_seeker_id"`
	JobSeeker   JobSeeker `gorm:"foreignKey:JobSeekerID" json:"job_seeker,omitempty"`
	Status      string    `gorm:"default:'applied'" json:"status" validate:"oneof=applied shortlisted rejected selected"`
	AppliedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"applied_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}