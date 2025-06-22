-- Create firm_profiles table (extension for employers with type 'firm')
CREATE TABLE firm_profiles (
    id SERIAL PRIMARY KEY,
    employer_id INTEGER UNIQUE NOT NULL REFERENCES employers(id) ON DELETE CASCADE,
    year_founded INTEGER,
    firm_size VARCHAR(50),
    legal_entity_type VARCHAR(100),
    primary_discipline VARCHAR(255) NOT NULL,
    secondary_disciplines JSON,
    instagram_url VARCHAR(500),
    linkedin_url VARCHAR(500),
    preferred_duration VARCHAR(100),
    stipend_range VARCHAR(100),
    project_images JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_firm_profiles_employer_id ON firm_profiles(employer_id);
CREATE INDEX idx_firm_profiles_primary_discipline ON firm_profiles(primary_discipline);
CREATE INDEX idx_firm_profiles_year_founded ON firm_profiles(year_founded);
CREATE INDEX idx_firm_profiles_created_at ON firm_profiles(created_at);