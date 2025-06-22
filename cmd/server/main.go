package main

import (
	"log"
	"os"

	"github.com/dekkaladiwakar/black-pages-backend/internal/handlers"
	"github.com/dekkaladiwakar/black-pages-backend/internal/middleware"
	"github.com/dekkaladiwakar/black-pages-backend/internal/repositories"
	"github.com/dekkaladiwakar/black-pages-backend/internal/services"
	"github.com/dekkaladiwakar/black-pages-backend/internal/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	utils.InitDatabase()
	defer utils.CloseDB()

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	router := gin.Default()

	// Setup CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:3000",                    // Local development
		"https://blackpages.up.railway.app",       // Railway production
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.AllowCredentials = true // Allow credentials for auth headers
	router.Use(cors.New(config))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		response := gin.H{
			"status":  "ok",
			"message": "Black Pages API is running",
			"version": "1.0.0",
		}

		// Check database health
		if err := utils.DatabaseHealthCheck(); err != nil {
			c.JSON(503, gin.H{
				"status":   "error",
				"message":  "Database connection failed",
				"database": err.Error(),
			})
			return
		}

		response["database"] = "connected"
		c.JSON(200, response)
	})

	// Initialize repositories and services
	userRepo := repositories.NewUserRepository(utils.GetDB())
	jobSeekerRepo := repositories.NewJobSeekerRepository(utils.GetDB())
	employerRepo := repositories.NewEmployerRepository(utils.GetDB())
	jobRepo := repositories.NewJobRepository(utils.GetDB())
	applicationRepo := repositories.NewApplicationRepository(utils.GetDB())
	studentProfileRepo := repositories.NewStudentProfileRepository(utils.GetDB())
	firmProfileRepo := repositories.NewFirmProfileRepository(utils.GetDB())
	
	authService := services.NewAuthService(userRepo)
	jobSeekerService := services.NewJobSeekerService(jobSeekerRepo, userRepo)
	employerService := services.NewEmployerService(employerRepo, userRepo)
	jobService := services.NewJobService(jobRepo, employerRepo)
	applicationService := services.NewApplicationService(applicationRepo, jobRepo, jobSeekerRepo, employerRepo)
	studentProfileService := services.NewStudentProfileService(studentProfileRepo, jobSeekerRepo)
	firmProfileService := services.NewFirmProfileService(firmProfileRepo, employerRepo)
	
	storageService := services.NewMockS3Service()
	fileService := services.NewFileService(storageService)
	
	authHandler := handlers.NewAuthHandler(authService)
	jobSeekerHandler := handlers.NewJobSeekerHandler(jobSeekerService)
	employerHandler := handlers.NewEmployerHandler(employerService)
	uploadHandler := handlers.NewUploadHandler(fileService, jobSeekerService)
	jobHandler := handlers.NewJobHandler(jobService, employerService)
	applicationHandler := handlers.NewApplicationHandler(applicationService, jobSeekerService, employerService)
	profileExtensionHandler := handlers.NewProfileExtensionHandler(studentProfileService, firmProfileService, jobSeekerService, employerService)

	// API routes group
	api := router.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		// Authentication routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/me", middleware.AuthRequired(), authHandler.GetMe)
		}

		// Job Seeker routes
		jobSeekers := api.Group("/job-seekers")
		jobSeekers.Use(middleware.AuthRequired())
		jobSeekers.Use(middleware.RequireRole("job_seeker"))
		{
			jobSeekers.POST("/profile", jobSeekerHandler.CreateProfile)
			jobSeekers.GET("/profile", jobSeekerHandler.GetProfile)
			jobSeekers.PUT("/profile", jobSeekerHandler.UpdateProfile)
			jobSeekers.GET("/profile/full", profileExtensionHandler.GetJobSeekerProfileWithExtensions)
			
			// Student profile extension routes
			jobSeekers.POST("/student-profile", profileExtensionHandler.CreateStudentProfile)
			jobSeekers.GET("/student-profile", profileExtensionHandler.GetStudentProfile)
			jobSeekers.PUT("/student-profile", profileExtensionHandler.UpdateStudentProfile)
			jobSeekers.DELETE("/student-profile", profileExtensionHandler.DeleteStudentProfile)
		}

		// Employer routes
		employers := api.Group("/employers")
		employers.Use(middleware.AuthRequired())
		employers.Use(middleware.RequireRole("employer"))
		{
			employers.POST("/profile", employerHandler.CreateProfile)
			employers.GET("/profile", employerHandler.GetProfile)
			employers.PUT("/profile", employerHandler.UpdateProfile)
			employers.GET("/profile/full", profileExtensionHandler.GetEmployerProfileWithExtensions)
			
			// Firm profile extension routes
			employers.POST("/firm-profile", profileExtensionHandler.CreateFirmProfile)
			employers.GET("/firm-profile", profileExtensionHandler.GetFirmProfile)
			employers.PUT("/firm-profile", profileExtensionHandler.UpdateFirmProfile)
			employers.DELETE("/firm-profile", profileExtensionHandler.DeleteFirmProfile)
		}

		// Upload routes (job seekers only)
		upload := api.Group("/upload")
		upload.Use(middleware.AuthRequired())
		upload.Use(middleware.RequireRole("job_seeker"))
		{
			upload.POST("/resume", uploadHandler.UploadResume)
			upload.POST("/portfolio", uploadHandler.UploadPortfolio)
		}

		// Public job routes (anyone can browse jobs)
		jobs := api.Group("/jobs")
		{
			jobs.GET("", jobHandler.GetAllJobs)          // Browse all jobs with filters
			jobs.GET("/:id", jobHandler.GetJob)          // Get specific job details
		}
		
		// Job filter options endpoint (separate to avoid route conflicts)
		api.GET("/jobs/filters", jobHandler.GetJobFilterOptions)

		// Employer job management routes
		employerJobs := api.Group("/employers/jobs")
		employerJobs.Use(middleware.AuthRequired())
		employerJobs.Use(middleware.RequireRole("employer"))
		{
			employerJobs.POST("", jobHandler.CreateJob)                    // Create new job
			employerJobs.GET("", jobHandler.GetEmployerJobs)               // Get employer's jobs
			employerJobs.GET("/:id", jobHandler.GetJob)                    // Get specific job
			employerJobs.PUT("/:id", jobHandler.UpdateJob)                 // Update job
			employerJobs.DELETE("/:id", jobHandler.DeleteJob)              // Delete job
			employerJobs.PUT("/:id/toggle", jobHandler.ToggleJobStatus)    // Toggle active/inactive
		}

		// Employer dashboard
		employerDashboard := api.Group("/employers/dashboard")
		employerDashboard.Use(middleware.AuthRequired())
		employerDashboard.Use(middleware.RequireRole("employer"))
		{
			employerDashboard.GET("", jobHandler.GetEmployerDashboard)     // Dashboard stats
		}

		// Job Seeker Application routes
		applications := api.Group("/applications")
		applications.Use(middleware.AuthRequired())
		applications.Use(middleware.RequireRole("job_seeker"))
		{
			applications.POST("", applicationHandler.ApplyToJob)                    // Apply to job
			applications.GET("", applicationHandler.GetMyApplications)              // Get my applications
			applications.GET("/stats", applicationHandler.GetMyApplicationStats)    // Get application stats
			applications.DELETE("/:id", applicationHandler.WithdrawApplication)     // Withdraw application
		}

		// Employer Application Management routes
		employerApplications := api.Group("/employers/jobs/:id/applications")
		employerApplications.Use(middleware.AuthRequired())
		employerApplications.Use(middleware.RequireRole("employer"))
		{
			employerApplications.GET("", applicationHandler.GetJobApplications)     // Get job applications
			employerApplications.GET("/stats", applicationHandler.GetJobApplicationStats) // Get application stats for job
		}

		// Application status management (employers only)
		applicationStatus := api.Group("/applications/:id")
		applicationStatus.Use(middleware.AuthRequired())
		applicationStatus.Use(middleware.RequireRole("employer"))
		{
			applicationStatus.PUT("/status", applicationHandler.UpdateApplicationStatus) // Update application status
		}
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Fatal(router.Run(":" + port))
}