package handlers

import (
	"net/http"

	"github.com/dekkaladiwakar/black-pages-backend/internal/middleware"
	"github.com/dekkaladiwakar/black-pages-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	fileService      services.FileService
	jobSeekerService services.JobSeekerService
}

func NewUploadHandler(fileService services.FileService, jobSeekerService services.JobSeekerService) *UploadHandler {
	return &UploadHandler{
		fileService:      fileService,
		jobSeekerService: jobSeekerService,
	}
}

func (h *UploadHandler) UploadResume(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User not authenticated",
		})
		return
	}

	file, err := c.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "No file uploaded or invalid form data",
		})
		return
	}

	url, err := h.fileService.UploadResume(userID, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	_, err = h.jobSeekerService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Job seeker profile not found. Create profile first.",
		})
		return
	}

	updateReq := services.UpdateJobSeekerRequest{
		ResumeURL: url,
	}
	
	updatedProfile, err := h.jobSeekerService.UpdateProfile(userID, updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update profile with resume URL",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Resume uploaded successfully",
		"data": gin.H{
			"url":     url,
			"profile": updatedProfile,
		},
	})
}

func (h *UploadHandler) UploadPortfolio(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User not authenticated",
		})
		return
	}

	file, err := c.FormFile("portfolio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "No file uploaded or invalid form data",
		})
		return
	}

	url, err := h.fileService.UploadPortfolio(userID, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	_, err = h.jobSeekerService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Job seeker profile not found. Create profile first.",
		})
		return
	}

	updateReq := services.UpdateJobSeekerRequest{
		PortfolioURL: url,
	}
	
	updatedProfile, err := h.jobSeekerService.UpdateProfile(userID, updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update profile with portfolio URL",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Portfolio uploaded successfully",
		"data": gin.H{
			"url":     url,
			"profile": updatedProfile,
		},
	})
}