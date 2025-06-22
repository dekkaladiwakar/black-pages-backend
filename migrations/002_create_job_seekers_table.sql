-- Create job_seekers table
CREATE TABLE job_seekers (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    full_name VARCHAR(255) NOT NULL,
    job_seeker_type VARCHAR(50) NOT NULL CHECK (job_seeker_type IN ('student', 'professional', 'freelancer')),
    current_city VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    desired_field VARCHAR(255) NOT NULL,
    resume_url VARCHAR(500) NOT NULL,
    portfolio_url VARCHAR(500),
    skills JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_job_seekers_user_id ON job_seekers(user_id);
CREATE INDEX idx_job_seekers_type ON job_seekers(job_seeker_type);
CREATE INDEX idx_job_seekers_city ON job_seekers(current_city);
CREATE INDEX idx_job_seekers_created_at ON job_seekers(created_at);