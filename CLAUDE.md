# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Black Pages is a generic job marketplace platform designed to connect job seekers with employers. The project is architected as a modular platform that can be extended for specific industries and use cases. The initial implementation focuses on architecture students and firms but is designed to scale to other industries.

## Current Implementation Status

**✅ COMPLETED:**
- Go project structure with proper module initialization (`go.mod` with all dependencies)
- Docker development environment (PostgreSQL + Go app containers)
- Complete database models with GORM and comprehensive validation tags:
  - Core models: User, JobSeeker, Employer, Job, Application
  - Extension models: StudentProfile, FirmProfile
- Basic Gin server with CORS, health checks, and API structure
- Database connection utilities with auto-migrations
- Project directory structure following Go conventions

**✅ COMPLETED PHASES:**

**Phase 1-3: Core Platform (COMPLETE)**
- ✅ Full authentication system (register, login, JWT, middleware)
- ✅ Complete profile management system:
  - ✅ JobSeeker profiles (create, read, update) with skills JSON
  - ✅ Employer profiles (create, read, update) with company data
  - ✅ Role-based route protection working
  - ✅ All endpoints tested with real data
- ✅ Service Layer + Repository pattern fully implemented
- ✅ Database operations confirmed working across all models

**Phase 4: Job Management System (COMPLETE)**
- ✅ Job repository with CRUD operations and Builder pattern for filtering
- ✅ Job service with comprehensive business logic and validation
- ✅ Job handlers for all CRUD endpoints (create, read, update, delete, toggle status)
- ✅ Public job browsing with advanced filtering (industry, location, type, etc.)
- ✅ Employer job management dashboard with statistics
- ✅ All job endpoints tested and working

**Phase 5: Application System (COMPLETE)**
- ✅ Application repository with full CRUD operations
- ✅ Application service with business logic (apply, withdraw, status updates)
- ✅ Application handlers for job seekers and employers
- ✅ Job application workflow (apply → shortlist → select/reject)
- ✅ Application statistics and tracking
- ✅ All application endpoints tested and working

**Phase 6: Student/Firm Profile Extensions (COMPLETE)**
- ✅ StudentProfile repository and service for academic/experience data
- ✅ FirmProfile repository and service for architecture firm specifics
- ✅ Profile extension handlers with full CRUD operations
- ✅ Enhanced "full" profile endpoints showing base + extension data
- ✅ Type validation (students only get student profiles, firms only get firm profiles)
- ✅ JSON array handling for skills, projects, and disciplines
- ✅ All extension endpoints tested and working

**Phase 7: Critical Business Logic Testing (COMPLETE)**
- ✅ Service layer test infrastructure with testify framework
- ✅ Critical validation logic tests (100% passing):
  - ✅ Application deadline validation (past/future deadlines)
  - ✅ User type validation (student/professional/freelancer restrictions)
  - ✅ Employer type validation (firm/corporation/startup restrictions)
  - ✅ Application status validation (withdrawal rules)
  - ✅ Required field validation (resume/portfolio requirements)
  - ✅ JSON array conversion utilities
- ✅ Comprehensive test coverage for all critical business rules

**Phase 8: Project Structure Refactoring (COMPLETE)**
- ✅ Models split into domain files: `user.go`, `job.go`, `application.go`, `profile.go`
- ✅ Explicit migration files for all 7 models (production-ready with indexes)
- ✅ Migration runner tool with version tracking (`cmd/migrate/main.go`)
- ✅ Tests moved to proper `test/services/` directory structure
- ✅ Removed unused folders (`static/`, `templates/`, `tmp/`)
- ✅ Updated imports and utility functions (shared `utils.ArrayToJSON`)
- ✅ All functionality verified working after restructure
- ✅ AutoMigrate removed - using explicit migrations

**Phase 9: Frontend Development (COMPLETE)**
- ✅ Next.js 15.3.4 with TypeScript and Tailwind CSS setup
- ✅ Complete authentication system (login/register with validation)
- ✅ Dynamic job filtering with backend API integration
- ✅ Responsive dashboard layouts for employers and job seekers
- ✅ Onboarding flow with skip option and profile completion blocking
- ✅ Reusable component architecture (UI, forms, layouts)
- ✅ Profile management system with view/edit functionality
- ✅ Action blocking system with helpful user guidance
- ✅ Mobile-first responsive design following brand guidelines
- ✅ Professional navigation and user experience flows

**Phase 10: Job Creation & Management Frontend (COMPLETE)**
- ✅ Complete job posting form with comprehensive validation
- ✅ Real-time form validation system (useFormValidation hook)
- ✅ Employer job management dashboard (/dashboard/employer/jobs)
- ✅ Job creation flow with success state and user control
- ✅ Authentication-aware navigation system throughout app
- ✅ Smart breadcrumb navigation based on user type
- ✅ Form validation for all required fields and business rules
- ✅ Skills management with add/remove functionality
- ✅ Date validation preventing past deadlines
- ✅ Success flow with "Post Another Job" and navigation options

**Phase 11: Navigation System Overhaul (COMPLETE)**
- ✅ Authentication-aware navigation hook (useAuthNavigation)
- ✅ Smart "Home" routing based on user authentication status
- ✅ Role-based navigation (job seekers vs employers get different flows)
- ✅ Fixed all 404 navigation issues and broken links
- ✅ Context-aware breadcrumbs in job detail pages
- ✅ Proper routing between public jobs and employer job management
- ✅ Dashboard quick actions now navigate to correct pages

**✅ PHASE 11 COMPLETE - Inactive Jobs Display Issue RESOLVED:**

**🐛 ISSUE IDENTIFIED:**
User reported inactive jobs are not showing in the "Inactive Jobs" tab even though 2 jobs were deactivated and appear as inactive in the dashboard recent jobs section.

**🔍 DEBUGGING PROCESS:**
- ✅ Identified backend API inconsistency:
  - Dashboard endpoint (`/api/employers/dashboard`) uses `GetByEmployerID()` → Returns ALL jobs ✅
  - Jobs page endpoint (`/api/employers/jobs`) uses `GetWithFilters()` → Returns empty array ❌
- ✅ **ROOT CAUSE FOUND:** Backend filtering logic issue
  - When we removed hardcoded `WHERE is_active = true` from QueryBuilder
  - GetWithFilters() was not handling null IsActive filter properly
  - Frontend receives empty array instead of all jobs
- ✅ **BACKEND FIXES IMPLEMENTED:**
  - Added `IsActive` field to JobFilters struct (repository + service layers)
  - Added `WithIsActive()` method to JobQueryBuilder
  - Updated GetEmployerJobs to support active/inactive filtering
  - Updated GetAllJobs to force active-only for public browsing
  - Fixed GetWithFilters to handle null IsActive filter correctly

**🔧 RESOLUTION STEPS:**
- ✅ Backend changes deployed and container restarted
- ✅ Added detailed API response logging to frontend and backend
- ✅ **CRITICAL DISCOVERY:** Docker container was using old built image without code changes
- ✅ **SOLUTION:** Rebuilt Docker image with `docker compose build app`
- ✅ **VERIFICATION:** API now correctly returns all jobs (2 inactive jobs) via `/api/employers/jobs`

**📋 FINAL STATUS:**
✅ **ISSUE COMPLETELY RESOLVED** - The inactive jobs display functionality now works correctly:
1. ✅ Backend API correctly returns all jobs (both active and inactive)
2. ✅ Frontend receives proper job data arrays
3. ✅ Inactive jobs tab can now display the 2 deactivated job postings
4. ✅ Employer job management system fully functional

**✅ PHASE 12 COMPLETE - Job Seeker Frontend (COMPLETE):**

✅ **COMPLETE:** Job Seeker Frontend System with comprehensive application workflow:
- ✅ **Complete Job Application Workflow:** `/jobs/[id]/apply` page with profile completion gating
- ✅ **Job Seeker Applications Management:** `/dashboard/job-seeker/applications` with filtering and status tracking
- ✅ **Application Detail Pages:** `/dashboard/job-seeker/applications/[id]` with full application and job information
- ✅ **Enhanced Job Seeker Profile Management:** `/dashboard/job-seeker/profile` with form validation and skills management
- ✅ **Authentication-Aware Apply Buttons:** Smart routing based on user type and authentication status
- ✅ **Profile Completion Validation:** Required fields checking before allowing job applications
- ✅ **Application Status Tracking:** Real-time status updates (applied, shortlisted, selected, rejected)
- ✅ **Application Withdrawal:** One-click withdrawal with confirmation for active applications
- ✅ **Comprehensive Form Validation:** Real-time validation with user-friendly error messaging
- ✅ **Mobile-First Responsive Design:** Professional UI following brand guidelines
- ✅ **Complete API Integration:** All job seeker endpoints working with proper error handling

**✅ PHASE 12.1 COMPLETE - Job Seeker Frontend Bug Fixes:**

🐛 **CRITICAL BUGS RESOLVED:**
- ✅ **Authentication Race Condition Fixed:** Apply button was redirecting authenticated users to login due to race condition in authentication loading state
- ✅ **Profile Data Mismatch Fixed:** Dashboard and profile page now use consistent API endpoints and data parsing for profile completion status
- ✅ **API Method Bug Fixed:** Fixed undefined HTTP method in API requests causing authentication failures
- ✅ **Navigation System Enhanced:** Apply buttons now properly handle authentication state and user types
- ✅ **Login Redirect Handling:** Login page now properly handles redirect parameters after successful authentication
- ✅ **Role-Based UI Elements:** "Post a Job" CTA now hidden for job seekers, only visible for employers/unauthenticated users
- ✅ **Backend API Validation:** Comprehensive testing confirmed backend is working perfectly with proper validation

🔧 **TECHNICAL IMPROVEMENTS:**
- **Enhanced Error Handling:** Added proper loading state checks in authentication-dependent components
- **Consistent Data Flow:** Standardized API calls between dashboard and profile pages
- **Race Condition Prevention:** Added authentication loading state guards in critical user flows
- **URL Encoding:** Fixed redirect parameter encoding for special characters
- **Smart Navigation:** Enhanced useAuthNavigation hook with proper loading states

**📋 NEXT PRIORITIES:**

**Phase 13: Enhanced Features**
1. S3 file upload integration (replace mock service)
2. Extended test coverage (integration and repository tests)
3. Production deployment optimization
4. Performance optimization and caching

**🎯 CHECKPOINT STATUS (Phase 12 Complete - FULL PLATFORM COMPLETE):**
The job marketplace platform now includes a **COMPLETE, PRODUCTION-READY** full-stack application with:
- **COMPLETE:** Full user authentication and profile management (both backend + frontend)
- **COMPLETE:** Full employer job management system (create, edit, activate/deactivate, applications)
- **COMPLETE:** Complete job seeker application workflow (apply, track, withdraw)
- **COMPLETE:** Job seeker profile management with validation and skills
- **COMPLETE:** Application management system with status workflows for both sides
- **COMPLETE:** Authentication-aware navigation system throughout the app
- **COMPLETE:** Smart routing based on user type and authentication status
- **COMPLETE:** Comprehensive form validation system with real-time feedback
- **COMPLETE:** Profile completion gating for job applications
- **COMPLETE:** Application status tracking and withdrawal functionality
- **COMPLETE:** Docker deployment with proper container rebuilding workflow
- Professional Next.js frontend with TypeScript
- Responsive dashboard system for employers and job seekers
- Complete onboarding flow with profile completion tracking
- Reusable component architecture for scalability
- Mobile-first design following brand guidelines
- Action blocking system ensuring data quality
- All backend APIs working with frontend integration
- Student profile extensions (college, degree, internships, projects)
- Firm profile extensions (disciplines, social media, project images)
- Enhanced "full" profile endpoints with automatic extension loading
- Critical business logic testing with 100% pass rate
- Validated type safety and security rules
- Test infrastructure ready for expansion
- Production-ready project structure with domain separation
- Explicit migration system with version control
- Clean codebase organization for team collaboration

**🚀 WORKING API ENDPOINTS:**

**Authentication:**
- POST `/api/auth/register` - User registration
- POST `/api/auth/login` - User login
- GET `/api/auth/me` - Get current user (protected)

**Job Seeker Profile Management:**
- POST `/api/job-seekers/profile` - Create job seeker profile (protected)
- GET `/api/job-seekers/profile` - Get job seeker profile (protected)
- PUT `/api/job-seekers/profile` - Update job seeker profile (protected)
- GET `/api/job-seekers/profile/full` - Get profile with extensions (protected)

**Student Profile Extensions (Students Only):**
- POST `/api/job-seekers/student-profile` - Create student profile extension (protected)
- GET `/api/job-seekers/student-profile` - Get student profile extension (protected)
- PUT `/api/job-seekers/student-profile` - Update student profile extension (protected)
- DELETE `/api/job-seekers/student-profile` - Delete student profile extension (protected)

**Employer Profile Management:**
- POST `/api/employers/profile` - Create employer profile (protected)
- GET `/api/employers/profile` - Get employer profile (protected)
- PUT `/api/employers/profile` - Update employer profile (protected)
- GET `/api/employers/profile/full` - Get profile with extensions (protected)

**Firm Profile Extensions (Firms Only):**
- POST `/api/employers/firm-profile` - Create firm profile extension (protected)
- GET `/api/employers/firm-profile` - Get firm profile extension (protected)
- PUT `/api/employers/firm-profile` - Update firm profile extension (protected)
- DELETE `/api/employers/firm-profile` - Delete firm profile extension (protected)

**File Upload (Job Seekers Only):**
- POST `/api/upload/resume` - Upload resume (protected)
- POST `/api/upload/portfolio` - Upload portfolio (protected)

**Public Job Browsing:**
- GET `/api/jobs` - Browse all active jobs with filtering
- GET `/api/jobs/:id` - Get specific job details

**Employer Job Management:**
- POST `/api/employers/jobs` - Create new job (protected)
- GET `/api/employers/jobs` - Get employer's jobs with filtering (protected)
- GET `/api/employers/jobs/:id` - Get specific job (protected)
- PUT `/api/employers/jobs/:id` - Update job (protected)
- DELETE `/api/employers/jobs/:id` - Delete job (protected)
- PUT `/api/employers/jobs/:id/toggle` - Toggle job active/inactive (protected)

**Employer Dashboard:**
- GET `/api/employers/dashboard` - Get dashboard statistics (protected)

**Job Applications (Job Seekers):**
- POST `/api/applications` - Apply to a job (protected)
- GET `/api/applications` - Get my applications (protected)
- GET `/api/applications/stats` - Get application statistics (protected)
- DELETE `/api/applications/:id` - Withdraw application (protected)

**Application Management (Employers):**
- GET `/api/employers/jobs/:id/applications` - Get job applications (protected)
- GET `/api/employers/jobs/:id/applications/stats` - Get application stats for job (protected)
- PUT `/api/applications/:id/status` - Update application status (protected)

## Development Commands

```bash
# Start PostgreSQL container
docker compose up -d postgres

# Run database migrations (if needed)
go run cmd/migrate/main.go up
go run cmd/migrate/main.go status
go run cmd/migrate/main.go down

# Start development server
go run cmd/server/main.go

# Test health endpoint
curl http://localhost:8080/health

# Test API
curl http://localhost:8080/api/ping

# Run tests
go test ./test/services/ -v
```

# CLAUDE.md - Black Pages Development Instructions

## 🎯 Core Development Principles

### 🚨 CRITICAL: Frontend Code Responsibility

**CLAUDE IS FULLY RESPONSIBLE FOR FRONTEND CODE QUALITY, ARCHITECTURE, AND BEST PRACTICES**

- **User Focus**: User handles backend/business logic; Claude handles ALL frontend concerns
- **Code Quality**: Claude must ensure reusable, maintainable, scalable frontend code
- **Architecture**: Implement proper component structure, hooks, state management
- **Best Practices**: Follow React/Next.js conventions, TypeScript standards, accessibility
- **Performance**: Optimize for mobile-first, responsive design, loading states
- **Patterns**: Use established patterns (custom hooks, component composition, etc.)
- **NO Cutting Corners**: Even under time pressure, maintain code quality standards
- **Documentation**: Self-documenting code through clear naming and structure

### Code Style & Comments

- **NO method/function descriptions** - Go code should be self-documenting
- **Only add comments for complex business logic** - Keep them concise (1-line max)
- **Use descriptive variable/function names** instead of comments
- **Prefer explicit over clever** - Readability > brevity

### Speed & Pragmatism

- **MVP first** - Build working features, optimize later
- **Leverage frameworks** - Use GORM, Gin built-ins instead of custom solutions
- **No tests during initial development** - Focus on functionality first
- **Environment-driven config** - Use .env files, no hardcoding

## 🏗️ Design Patterns for Black Pages

### 1. Repository Pattern (Database Layer)

**Use for:** All database operations

```go
type UserRepository interface {
    Create(user *User) error
    GetByEmail(email string) (*User, error)
    Update(user *User) error
}

type userRepository struct {
    db *gorm.DB
}
```

**Why:** Clean separation of database logic, easy to mock/replace

### 2. Service Layer Pattern (Business Logic)

**Use for:** All business operations between handlers and repository

```go
type AuthService interface {
    Register(req RegisterRequest) (*User, error)
    Login(email, password string) (string, error)
}
```

**Why:** Keeps handlers thin, centralizes business rules

### 3. Factory Pattern (User Creation)

**Use for:** Creating different user types (JobSeeker, Employer, Student, Firm)

```go
type UserFactory interface {
    CreateUser(userType string, data interface{}) (interface{}, error)
}
```

**Why:** Handle different registration flows cleanly

### 4. Strategy Pattern (Profile Extensions)

**Use for:** Different profile completion flows

```go
type ProfileStrategy interface {
    CompleteProfile(userID uint, data interface{}) error
}

type StudentProfileStrategy struct{}
type FirmProfileStrategy struct{}
```

**Why:** Easy to add new user types without changing existing code

### 5. Builder Pattern (Query Filters)

**Use for:** Job search filters, application filters

```go
type JobQueryBuilder struct {
    query *gorm.DB
}

func (b *JobQueryBuilder) WithIndustry(industry string) *JobQueryBuilder
func (b *JobQueryBuilder) WithLocation(city string) *JobQueryBuilder
func (b *JobQueryBuilder) Build() *gorm.DB
```

**Why:** Complex filtering logic becomes readable and maintainable

### 6. Middleware Pattern (Cross-cutting Concerns)

**Use for:** Authentication, logging, validation, CORS

```go
func AuthMiddleware() gin.HandlerFunc
func LoggingMiddleware() gin.HandlerFunc
func ValidationMiddleware() gin.HandlerFunc
```

**Why:** Built into Gin, clean separation of concerns

## 📁 File Organization

```
internal/
├── handlers/          # HTTP handlers (thin)
│   ├── auth.go
│   ├── jobs.go
│   └── applications.go
├── services/          # Business logic
│   ├── auth_service.go
│   ├── job_service.go
│   └── user_service.go
├── repositories/      # Database operations
│   ├── user_repo.go
│   ├── job_repo.go
│   └── application_repo.go
├── models/           # Database models
│   ├── user.go
│   ├── job.go
│   └── application.go
├── middleware/       # Gin middleware
├── utils/           # Helper functions
└── config/          # Configuration
```

## 🚀 Quick Development Guidelines

### GORM Best Practices

```go
// Use auto-migration for speed
db.AutoMigrate(&User{}, &JobSeeker{}, &Employer{})

// Use preloading for relationships
db.Preload("JobSeeker").Find(&users)

// Use struct tags for validation
type User struct {
    Email    string `gorm:"uniqueIndex" json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}
```

### Error Handling Pattern

```go
type APIResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func SuccessResponse(data interface{}) APIResponse
func ErrorResponse(message string) APIResponse
```

### Environment Configuration

```go
type Config struct {
    DBHost     string `env:"DB_HOST" envDefault:"localhost"`
    DBPort     string `env:"DB_PORT" envDefault:"5432"`
    JWTSecret  string `env:"JWT_SECRET,required"`
    ServerPort string `env:"SERVER_PORT" envDefault:"8080"`
}
```

### Authentication Flow

```go
// JWT middleware for protected routes
func AuthRequired() gin.HandlerFunc

// Role-based access
func RequireRole(role string) gin.HandlerFunc

// User context in handlers
func GetCurrentUser(c *gin.Context) (*User, error)
```

## 📊 Database Conventions

### Model Structure

```go
type BaseModel struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
    BaseModel
    Email        string `gorm:"uniqueIndex" json:"email"`
    PasswordHash string `gorm:"column:password_hash" json:"-"`
}
```

### JSON Field Handling

```go
type JobSeeker struct {
    Skills datatypes.JSON `gorm:"type:json" json:"skills"`
}

// Helper methods for JSON fields
func (j *JobSeeker) GetSkills() []string
func (j *JobSeeker) SetSkills(skills []string)
```

## 🔧 API Development

### Handler Structure

```go
func (h *JobHandler) CreateJob(c *gin.Context) {
    var req CreateJobRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, ErrorResponse("Invalid request"))
        return
    }

    job, err := h.jobService.CreateJob(req)
    if err != nil {
        c.JSON(500, ErrorResponse(err.Error()))
        return
    }

    c.JSON(201, SuccessResponse(job))
}
```

### Validation Tags

```go
type CreateJobRequest struct {
    Title          string   `json:"title" binding:"required,max=100"`
    JobType        string   `json:"job_type" binding:"required,oneof=internship full_time contract"`
    RequiredSkills []string `json:"required_skills" binding:"required,min=1"`
    IsPaid         bool     `json:"is_paid" binding:"required"`
}
```

## 🔄 Development Workflow

### Phase Implementation Order

1. **Models first** - Define all structs and relationships
2. **Repository layer** - Basic CRUD operations
3. **Service layer** - Business logic implementation
4. **Handlers** - HTTP endpoints
5. **Middleware** - Auth, validation, logging
6. **Integration** - Connect all layers

### Quick Testing Commands

```bash
# Start with Docker
docker-compose up -d

# Run migrations
go run cmd/server/main.go -migrate

# Start development server
air # or go run cmd/server/main.go
```

### Database Seeding (Optional)

```go
func SeedDatabase(db *gorm.DB) {
    // Create sample users for testing
    // Only for development environment
}
```

## ⚡ Performance Considerations

### Database Optimization

- Use database indexes on frequently queried fields
- Implement pagination for list endpoints
- Use select specific fields when not needing full objects
- Cache frequently accessed data (Redis later)

### File Upload Strategy

```go
// Start with local storage, move to S3 later
func SaveFile(file *multipart.FileHeader, dest string) error

// Validate file types and sizes
func ValidateUpload(file *multipart.FileHeader) error
```

## 🚫 What NOT to Do

- **Don't create custom ORMs** - Use GORM
- **Don't write custom auth** - Use JWT + bcrypt
- **Don't over-engineer** - Keep it simple for MVP
- **Don't add caching initially** - Add when needed
- **Don't create custom validation** - Use Gin's binding
- **Don't write complex SQL** - Use GORM queries

## 📈 Future Expansion Points

When adding new features, consider:

- New user types: Add new strategy implementations
- New industries: Extend factory patterns
- New filters: Extend builder patterns
- New notification types: Add observer patterns

---

**Remember: Build fast, build working, optimize later. Focus on delivering the MVP as defined in the PRD first.**
