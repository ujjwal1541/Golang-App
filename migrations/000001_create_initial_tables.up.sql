-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('receptionist', 'doctor')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create patients table
CREATE TABLE IF NOT EXISTS patients (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE NOT NULL,
    gender VARCHAR(20) NOT NULL CHECK (gender IN ('male', 'female', 'other')),
    contact_number VARCHAR(20) NOT NULL,
    email VARCHAR(255) UNIQUE,
    address TEXT NOT NULL,
    emergency_name VARCHAR(255),
    emergency_number VARCHAR(20),
    blood_group VARCHAR(10),
    allergies TEXT,
    medical_history TEXT,
    current_medication TEXT,
    notes TEXT,
    registered_by INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create index on patients for searching
CREATE INDEX idx_patient_first_name ON patients(first_name);
CREATE INDEX idx_patient_last_name ON patients(last_name);
CREATE INDEX idx_patient_email ON patients(email);
CREATE INDEX idx_patient_contact_number ON patients(contact_number);

-- Insert a default admin user with password 'admin123'
-- Password is hashed, this is just an example password for testing
INSERT INTO users (name, email, password, role) 
VALUES (
    'Admin Receptionist', 
    'receptionist@example.com', 
    '$2a$14$ajq8Q7fbtFRQvXpdCq7Jcuy.Rx1h/L4J60Otx.gyNLbAYctGMJ9tK', 
    'receptionist'
) ON CONFLICT DO NOTHING;

-- Insert a default doctor user with password 'doctor123'
INSERT INTO users (name, email, password, role) 
VALUES (
    'Admin Doctor', 
    'doctor@example.com', 
    '$2a$14$ajq8Q7fbtFRQvXpdCq7Jcuy.Rx1h/L4J60Otx.gyNLbAYctGMJ9tK', 
    'doctor'
) ON CONFLICT DO NOTHING; 