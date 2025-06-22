# Black Pages Backend

Backend API for the Black Pages job marketplace platform. Built with Go, Gin framework, and PostgreSQL.

## Features

- **Complete Authentication System**: JWT-based authentication with role-based access
- **Job Management**: Full CRUD operations for job postings with filtering
- **Application System**: Job application workflow with status tracking
- **Profile Management**: Comprehensive user profiles with extensions
- **File Upload**: Resume and portfolio upload functionality
- **Student/Firm Extensions**: Specialized profiles for architecture students and firms

## Technology Stack

- **Language**: Go 1.23
- **Framework**: Gin (HTTP web framework)
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT tokens
- **Containerization**: Docker with PostgreSQL container
- **Deployment**: Railway with Docker

## Architecture

- **Repository Pattern**: Clean database layer separation
- **Service Layer**: Business logic centralization
- **Middleware**: Authentication, CORS, and request logging
- **Migration System**: Version-controlled database schema

## Quick Start

### Prerequisites

- Go 1.23 or higher
- Docker and Docker Compose
- PostgreSQL (or use the provided Docker setup)

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/dekkaladiwakar/black-pages-backend.git
   cd black-pages-backend
   ```

2. **Set up environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start PostgreSQL with Docker**
   ```bash
   docker compose up -d postgres
   ```

4. **Run database migrations**
   ```bash
   go run cmd/migrate/main.go up
   ```

5. **Start the development server**
   ```bash
   go run cmd/server/main.go
   ```

6. **Test the API**
   ```bash
   curl http://localhost:8080/health
   ```

### Using Docker (Full Stack)

```bash
# Start both PostgreSQL and the Go application
docker compose up -d

# Check logs
docker compose logs -f app
```

## API Endpoints

### Authentication
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login  
- `GET /api/auth/me` - Get current user (protected)

### Job Management
- `GET /api/jobs` - Browse public jobs with filtering
- `GET /api/jobs/:id` - Get job details
- `POST /api/employers/jobs` - Create job (employers only)
- `GET /api/employers/jobs` - Get employer's jobs
- `PUT /api/employers/jobs/:id` - Update job
- `DELETE /api/employers/jobs/:id` - Delete job
- `PUT /api/employers/jobs/:id/toggle` - Toggle job status

### Profile Management
- `GET/POST/PUT /api/job-seekers/profile` - Job seeker profiles
- `GET/POST/PUT /api/employers/profile` - Employer profiles
- Profile extensions for students and firms

### Applications
- `POST /api/applications` - Apply to job
- `GET /api/applications` - Get user's applications
- `DELETE /api/applications/:id` - Withdraw application
- `PUT /api/applications/:id/status` - Update status (employers)

### File Upload
- `POST /api/upload/resume` - Upload resume
- `POST /api/upload/portfolio` - Upload portfolio

## Environment Variables

```bash
# Database
DATABASE_URL=postgres://postgres:postgres@localhost:5432/black_pages?sslmode=disable

# JWT
JWT_SECRET=your-secret-key-change-in-production

# Server
PORT=8080
```

## Database Migrations

```bash
# Run migrations
go run cmd/migrate/main.go up

# Check migration status
go run cmd/migrate/main.go status

# Rollback migrations
go run cmd/migrate/main.go down
```

## Testing

```bash
# Run tests
go test ./test/services/ -v

# Run with coverage
go test -cover ./test/services/
```

## Project Structure

```
├── cmd/
│   ├── migrate/          # Database migration tool
│   └── server/           # Main application server
├── internal/
│   ├── handlers/         # HTTP handlers
│   ├── middleware/       # Gin middleware
│   ├── models/          # Database models
│   ├── repositories/    # Database operations
│   ├── services/        # Business logic
│   └── utils/           # Helper functions
├── migrations/          # SQL migration files
├── test/               # Test files
├── docker-compose.yml  # Docker configuration
├── Dockerfile         # Container build
└── railway.json       # Railway deployment config
```

## Deployment

### Railway Deployment

1. **Connect your repository** to Railway
2. **Set environment variables** in Railway dashboard
3. **Deploy** - Railway will automatically build and deploy

### Environment Variables for Production

```bash
DATABASE_URL=postgresql://user:password@host:port/database
JWT_SECRET=your-secure-production-secret
PORT=8080
GIN_MODE=release
```

## Development Guidelines

- Follow the Repository + Service layer pattern
- Use GORM for database operations
- Implement proper error handling
- Add validation for all inputs
- Write tests for critical business logic
- Use environment variables for configuration

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is licensed under the MIT License.

## Related Repositories

- **Frontend**: [black-pages-frontend](https://github.com/dekkaladiwakar/black-pages-frontend)