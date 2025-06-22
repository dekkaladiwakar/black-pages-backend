-- Create student_profiles table (extension for job_seekers with type 'student')
CREATE TABLE student_profiles (
    id SERIAL PRIMARY KEY,
    job_seeker_id INTEGER UNIQUE NOT NULL REFERENCES job_seekers(id) ON DELETE CASCADE,
    college_name VARCHAR(255) NOT NULL,
    degree VARCHAR(255) NOT NULL,
    year_semester VARCHAR(100) NOT NULL,
    software_proficiency JSON,
    previous_internships JSON,
    freelance_projects JSON,
    preferred_start_month VARCHAR(50),
    preferred_duration VARCHAR(100),
    willing_to_relocate BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_student_profiles_job_seeker_id ON student_profiles(job_seeker_id);
CREATE INDEX idx_student_profiles_college ON student_profiles(college_name);
CREATE INDEX idx_student_profiles_degree ON student_profiles(degree);
CREATE INDEX idx_student_profiles_created_at ON student_profiles(created_at);