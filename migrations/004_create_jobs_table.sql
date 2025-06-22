-- Create jobs table
CREATE TABLE jobs (
    id SERIAL PRIMARY KEY,
    employer_id INTEGER NOT NULL REFERENCES employers(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    job_type VARCHAR(50) NOT NULL CHECK (job_type IN ('internship', 'full_time', 'contract')),
    industry VARCHAR(255) NOT NULL,
    target_audience VARCHAR(50) NOT NULL CHECK (target_audience IN ('students', 'professionals', 'any')),
    employment_mode VARCHAR(50) NOT NULL CHECK (employment_mode IN ('on_site', 'remote', 'hybrid')),
    start_month VARCHAR(50) NOT NULL,
    duration VARCHAR(100) NOT NULL,
    application_deadline TIMESTAMP NOT NULL,
    compensation_range VARCHAR(100),
    is_paid BOOLEAN NOT NULL,
    city VARCHAR(255) NOT NULL,
    state VARCHAR(255) NOT NULL,
    required_skills JSON NOT NULL,
    min_experience VARCHAR(100),
    portfolio_required BOOLEAN NOT NULL,
    resume_required BOOLEAN DEFAULT TRUE,
    description TEXT NOT NULL,
    about_team TEXT,
    contact_email VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_jobs_employer_id ON jobs(employer_id);
CREATE INDEX idx_jobs_job_type ON jobs(job_type);
CREATE INDEX idx_jobs_industry ON jobs(industry);
CREATE INDEX idx_jobs_target_audience ON jobs(target_audience);
CREATE INDEX idx_jobs_employment_mode ON jobs(employment_mode);
CREATE INDEX idx_jobs_city ON jobs(city);
CREATE INDEX idx_jobs_is_active ON jobs(is_active);
CREATE INDEX idx_jobs_application_deadline ON jobs(application_deadline);
CREATE INDEX idx_jobs_created_at ON jobs(created_at);