package handlers

import (
	"net/http"

	"github.com/dekkaladiwakar/black-pages-backend/internal/middleware"
	"github.com/dekkaladiwakar/black-pages-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type ProfileExtensionHandler struct {
	studentProfileService services.StudentProfileService
	firmProfileService    services.FirmProfileService
	jobSeekerService      services.JobSeekerService
	employerService       services.EmployerService
}

func NewProfileExtensionHandler(
	studentProfileService services.StudentProfileService,
	firmProfileService services.FirmProfileService,
	jobSeekerService services.JobSeekerService,
	employerService services.EmployerService,
) *ProfileExtensionHandler {
	return &ProfileExtensionHandler{
		studentProfileService: studentProfileService,
		firmProfileService:    firmProfileService,
		jobSeekerService:      jobSeekerService,
		employerService:       employerService,
	}
}

// Student Profile Extension Handlers

func (h *ProfileExtensionHandler) CreateStudentProfile(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User not authenticated",
		})
		return
	}

	// Get job seeker profile
	jobSeeker, err := h.jobSeekerService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Job seeker profile not found. Create profile first.",
		})
		return
	}

	var req services.CreateStudentProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	profile, err := h.studentProfileService.CreateProfile(jobSeeker.ID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Student profile created successfully",
		"data":    profile,
	})
}

func (h *ProfileExtensionHandler) GetStudentProfile(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User not authenticated",
		})
		return
	}

	// Get job seeker profile
	jobSeeker, err := h.jobSeekerService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Job seeker profile not found",
		})
		return
	}

	profile, err := h.studentProfileService.GetProfile(jobSeeker.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profile,
	})
}

func (h *ProfileExtensionHandler) UpdateStudentProfile(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User not authenticated",
		})
		return
	}

	// Get job seeker profile
	jobSeeker, err := h.jobSeekerService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Job seeker profile not found",
		})
		return
	}

	var req services.UpdateStudentProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	profile, err := h.studentProfileService.UpdateProfile(jobSeeker.ID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Student profile updated successfully",
		"data":    profile,
	})
}

func (h *ProfileExtensionHandler) DeleteStudentProfile(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User not authenticated",
		})
		return
	}

	// Get job seeker profile
	jobSeeker, err := h.jobSeekerService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Job seeker profile not found",
		})
		return
	}

	if err := h.studentProfileService.DeleteProfile(jobSeeker.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Student profile deleted successfully",
	})
}

// Firm Profile Extension Handlers

func (h *ProfileExtensionHandler) CreateFirmProfile(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User not authenticated",
		})
		return
	}

	// Get employer profile
	employer, err := h.employerService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Employer profile not found. Create profile first.",
		})
		return
	}

	var req services.CreateFirmProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	profile, err := h.firmProfileService.CreateProfile(employer.ID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Firm profile created successfully",
		"data":    profile,
	})
}

func (h *ProfileExtensionHandler) GetFirmProfile(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User not authenticated",
		})
		return
	}

	// Get employer profile
	employer, err := h.employerService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Employer profile not found",
		})
		return
	}

	profile, err := h.firmProfileService.GetProfile(employer.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profile,
	})
}

func (h *ProfileExtensionHandler) UpdateFirmProfile(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User not authenticated",
		})
		return
	}

	// Get employer profile
	employer, err := h.employerService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Employer profile not found",
		})
		return
	}

	var req services.UpdateFirmProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	profile, err := h.firmProfileService.UpdateProfile(employer.ID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Firm profile updated successfully",
		"data":    profile,
	})
}

func (h *ProfileExtensionHandler) DeleteFirmProfile(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User not authenticated",
		})
		return
	}

	// Get employer profile
	employer, err := h.employerService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Employer profile not found",
		})
		return
	}

	if err := h.firmProfileService.DeleteProfile(employer.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Firm profile deleted successfully",
	})
}

// Enhanced Profile Endpoints with Extensions

func (h *ProfileExtensionHandler) GetJobSeekerProfileWithExtensions(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User not authenticated",
		})
		return
	}

	// Get base job seeker profile
	jobSeeker, err := h.jobSeekerService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Job seeker profile not found",
		})
		return
	}

	response := gin.H{
		"success": true,
		"data": gin.H{
			"job_seeker": jobSeeker,
		},
	}

	// Try to get student profile if job seeker is a student
	if jobSeeker.JobSeekerType == "student" {
		studentProfile, err := h.studentProfileService.GetProfile(jobSeeker.ID)
		if err == nil {
			response["data"].(gin.H)["student_profile"] = studentProfile
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProfileExtensionHandler) GetEmployerProfileWithExtensions(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User not authenticated",
		})
		return
	}

	// Get base employer profile
	employer, err := h.employerService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Employer profile not found",
		})
		return
	}

	response := gin.H{
		"success": true,
		"data": gin.H{
			"employer": employer,
		},
	}

	// Try to get firm profile if employer is a firm
	if employer.EmployerType == "firm" {
		firmProfile, err := h.firmProfileService.GetProfile(employer.ID)
		if err == nil {
			response["data"].(gin.H)["firm_profile"] = firmProfile
		}
	}

	c.JSON(http.StatusOK, response)
}