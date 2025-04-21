package repositories

import (
	"healthcare-app/internal/models"

	"gorm.io/gorm"
)

// PatientRepository handles patient data operations
type PatientRepository struct {
	db *gorm.DB
}

// NewPatientRepository creates a new PatientRepository
func NewPatientRepository(db *gorm.DB) *PatientRepository {
	return &PatientRepository{db: db}
}

// Create creates a new patient
func (r *PatientRepository) Create(patient *models.Patient) error {
	return r.db.Create(patient).Error
}

// FindByID finds a patient by ID
func (r *PatientRepository) FindByID(id uint) (*models.Patient, error) {
	var patient models.Patient
	err := r.db.Where("id = ?", id).First(&patient).Error
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

// FindAll finds all patients
func (r *PatientRepository) FindAll(limit, offset int) ([]models.Patient, int64, error) {
	var patients []models.Patient
	var count int64

	query := r.db.Model(&models.Patient{})
	
	// Get total count
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Get patients with pagination
	if err := query.Limit(limit).Offset(offset).Find(&patients).Error; err != nil {
		return nil, 0, err
	}

	return patients, count, nil
}

// Update updates a patient
func (r *PatientRepository) Update(patient *models.Patient) error {
	return r.db.Save(patient).Error
}

// Delete deletes a patient
func (r *PatientRepository) Delete(id uint) error {
	return r.db.Delete(&models.Patient{}, id).Error
}

// SearchPatients searches for patients
func (r *PatientRepository) SearchPatients(searchTerm string, limit, offset int) ([]models.Patient, int64, error) {
	var patients []models.Patient
	var count int64

	query := r.db.Model(&models.Patient{}).Where(
		"first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ? OR contact_number LIKE ?",
		"%"+searchTerm+"%", "%"+searchTerm+"%", "%"+searchTerm+"%", "%"+searchTerm+"%",
	)
	
	// Get total count
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Get patients with pagination
	if err := query.Limit(limit).Offset(offset).Find(&patients).Error; err != nil {
		return nil, 0, err
	}

	return patients, count, nil
} 