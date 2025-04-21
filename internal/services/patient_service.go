package services

import (
	"errors"

	"healthcare-app/internal/models"
	"healthcare-app/internal/repositories"
)

// Predefined errors
var (
	ErrPatientNotFound = errors.New("patient not found")
)

// PaginationResponse represents a paginated response
type PaginationResponse struct {
	TotalItems int64       `json:"totalItems"`
	Items      interface{} `json:"items"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	TotalPages int         `json:"totalPages"`
}

// PatientService handles patient business logic
type PatientService struct {
	patientRepo *repositories.PatientRepository
}

// NewPatientService creates a new PatientService
func NewPatientService(patientRepo *repositories.PatientRepository) *PatientService {
	return &PatientService{
		patientRepo: patientRepo,
	}
}

// CreatePatient creates a new patient
func (s *PatientService) CreatePatient(req models.CreatePatientRequest, registeredByID uint) (*models.Patient, error) {
	patient := &models.Patient{
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		DateOfBirth:       req.DateOfBirth,
		Gender:            req.Gender,
		ContactNumber:     req.ContactNumber,
		Email:             req.Email,
		Address:           req.Address,
		EmergencyName:     req.EmergencyName,
		EmergencyNumber:   req.EmergencyNumber,
		BloodGroup:        req.BloodGroup,
		Allergies:         req.Allergies,
		MedicalHistory:    req.MedicalHistory,
		CurrentMedication: req.CurrentMedication,
		Notes:             req.Notes,
		RegisteredBy:      registeredByID,
	}

	if err := s.patientRepo.Create(patient); err != nil {
		return nil, err
	}

	return patient, nil
}

// GetPatient gets a patient by ID
func (s *PatientService) GetPatient(id uint) (*models.Patient, error) {
	patient, err := s.patientRepo.FindByID(id)
	if err != nil {
		return nil, ErrPatientNotFound
	}

	return patient, nil
}

// GetAllPatients gets all patients with pagination
func (s *PatientService) GetAllPatients(page, pageSize int) (*PaginationResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	patients, totalItems, err := s.patientRepo.FindAll(pageSize, offset)
	if err != nil {
		return nil, err
	}

	totalPages := (int(totalItems) + pageSize - 1) / pageSize
	
	return &PaginationResponse{
		TotalItems: totalItems,
		Items:      patients,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// UpdatePatient updates a patient
func (s *PatientService) UpdatePatient(id uint, req models.UpdatePatientRequest) (*models.Patient, error) {
	patient, err := s.patientRepo.FindByID(id)
	if err != nil {
		return nil, ErrPatientNotFound
	}

	patient.ApplyUpdates(req)

	if err := s.patientRepo.Update(patient); err != nil {
		return nil, err
	}

	return patient, nil
}

// UpdatePatientMedicalInfo updates a patient's medical information
func (s *PatientService) UpdatePatientMedicalInfo(id uint, req models.UpdatePatientMedicalRequest) (*models.Patient, error) {
	patient, err := s.patientRepo.FindByID(id)
	if err != nil {
		return nil, ErrPatientNotFound
	}

	patient.ApplyMedicalUpdates(req)

	if err := s.patientRepo.Update(patient); err != nil {
		return nil, err
	}

	return patient, nil
}

// DeletePatient deletes a patient
func (s *PatientService) DeletePatient(id uint) error {
	_, err := s.patientRepo.FindByID(id)
	if err != nil {
		return ErrPatientNotFound
	}

	return s.patientRepo.Delete(id)
}

// SearchPatients searches for patients
func (s *PatientService) SearchPatients(searchTerm string, page, pageSize int) (*PaginationResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	patients, totalItems, err := s.patientRepo.SearchPatients(searchTerm, pageSize, offset)
	if err != nil {
		return nil, err
	}

	totalPages := (int(totalItems) + pageSize - 1) / pageSize
	
	return &PaginationResponse{
		TotalItems: totalItems,
		Items:      patients,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
} 