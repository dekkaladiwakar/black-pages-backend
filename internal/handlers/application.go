package handlers

import (
	"net/http"
	"strconv"

	"github.com/dekkaladiwakar/black-pages-backend/internal/middleware"
	"github.com/dekkaladiwakar/black-pages-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type ApplicationHandler struct {
	applicationService services.ApplicationService
	jobSeekerService   services.JobSeekerService
	employerService    services.EmployerService
}

func NewApplicationHandler(
	applicationService services.ApplicationService,
	jobSeekerService services.JobSeekerService,
	employerService services.EmployerService,
) *ApplicationHandler {
	return &ApplicationHandler{
		applicationService: applicationService,
		jobSeekerService:   jobSeekerService,
		employerService:    employerService,
	}
}

func (h *ApplicationHandler) ApplyToJob(c *gin.Context) {
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

	var req services.ApplyJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	application, err := h.applicationService.ApplyToJob(jobSeeker.ID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Application submitted successfully",
		"data":    application,
	})
}

func (h *ApplicationHandler) GetMyApplications(c *gin.Context) {
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

	applications, err := h.applicationService.GetJobSeekerApplications(jobSeeker.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    applications,
	})
}

func (h *ApplicationHandler) GetJobApplications(c *gin.Context) {
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

	jobID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid job ID",
		})
		return
	}

	applications, err := h.applicationService.GetJobApplications(uint(jobID), employer.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    applications,
	})
}

func (h *ApplicationHandler) UpdateApplicationStatus(c *gin.Context) {
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

	applicationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid application ID",
		})
		return
	}

	var req services.UpdateApplicationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	application, err := h.applicationService.UpdateApplicationStatus(uint(applicationID), employer.ID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Application status updated successfully",
		"data":    application,
	})
}

func (h *ApplicationHandler) WithdrawApplication(c *gin.Context) {
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

	applicationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid application ID",
		})
		return
	}

	if err := h.applicationService.WithdrawApplication(uint(applicationID), jobSeeker.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Application withdrawn successfully",
	})
}

func (h *ApplicationHandler) GetMyApplicationStats(c *gin.Context) {
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

	stats, err := h.applicationService.GetApplicationStats(jobSeeker.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

func (h *ApplicationHandler) GetJobApplicationStats(c *gin.Context) {
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

	jobID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid job ID",
		})
		return
	}

	stats, err := h.applicationService.GetJobApplicationStats(uint(jobID), employer.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}