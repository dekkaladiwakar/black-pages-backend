-- Create employers table
CREATE TABLE employers (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company_name VARCHAR(255) NOT NULL,
    employer_type VARCHAR(50) NOT NULL CHECK (employer_type IN ('firm', 'corporation', 'startup')),
    industry VARCHAR(255) NOT NULL,
    primary_phone VARCHAR(20) NOT NULL,
    contact_person VARCHAR(255) NOT NULL,
    contact_person_desig VARCHAR(255) NOT NULL,
    city VARCHAR(255) NOT NULL,
    state VARCHAR(255) NOT NULL,
    pin_code VARCHAR(6) NOT NULL,
    website_url VARCHAR(500) NOT NULL,
    logo_url VARCHAR(500),
    is_hiring BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_employers_user_id ON employers(user_id);
CREATE INDEX idx_employers_type ON employers(employer_type);
CREATE INDEX idx_employers_industry ON employers(industry);
CREATE INDEX idx_employers_city ON employers(city);
CREATE INDEX idx_employers_is_hiring ON employers(is_hiring);
CREATE INDEX idx_employers_created_at ON employers(created_at);