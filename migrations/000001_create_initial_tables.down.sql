-- Drop patients table and its indexes
DROP INDEX IF EXISTS idx_patient_contact_number;
DROP INDEX IF EXISTS idx_patient_email;
DROP INDEX IF EXISTS idx_patient_last_name;
DROP INDEX IF EXISTS idx_patient_first_name;
DROP TABLE IF EXISTS patients;

-- Drop users table
DROP TABLE IF EXISTS users; 