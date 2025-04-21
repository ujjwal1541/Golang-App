package models

import (
	"time"

	"gorm.io/gorm"
)

// Patient represents a patient in the system
type Patient struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	FirstName       string         `json:"first_name" gorm:"not null"`
	LastName        string         `json:"last_name" gorm:"not null"`
	DateOfBirth     time.Time      `json:"date_of_birth" gorm:"not null"`
	Gender          string         `json:"gender" gorm:"not null"`
	ContactNumber   string         `json:"contact_number" gorm:"not null"`
	Email           string         `json:"email" gorm:"uniqueIndex"`
	Address         string         `json:"address" gorm:"not null"`
	EmergencyName   string         `json:"emergency_name"`
	EmergencyNumber string         `json:"emergency_number"`
	BloodGroup      string         `json:"blood_group"`
	Allergies       string         `json:"allergies"`
	MedicalHistory  string         `json:"medical_history"`
	CurrentMedication string       `json:"current_medication"`
	Notes           string         `json:"notes"`
	RegisteredBy    uint           `json:"registered_by" gorm:"not null"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// CreatePatientRequest represents a request to create a patient
type CreatePatientRequest struct {
	FirstName       string    `json:"first_name" binding:"required"`
	LastName        string    `json:"last_name" binding:"required"`
	DateOfBirth     time.Time `json:"date_of_birth" binding:"required"`
	Gender          string    `json:"gender" binding:"required,oneof=male female other"`
	ContactNumber   string    `json:"contact_number" binding:"required"`
	Email           string    `json:"email" binding:"omitempty,email"`
	Address         string    `json:"address" binding:"required"`
	EmergencyName   string    `json:"emergency_name"`
	EmergencyNumber string    `json:"emergency_number"`
	BloodGroup      string    `json:"blood_group"`
	Allergies       string    `json:"allergies"`
	MedicalHistory  string    `json:"medical_history"`
	CurrentMedication string  `json:"current_medication"`
	Notes           string    `json:"notes"`
}

// UpdatePatientRequest represents a request to update a patient
type UpdatePatientRequest struct {
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	DateOfBirth     time.Time `json:"date_of_birth"`
	Gender          string    `json:"gender" binding:"omitempty,oneof=male female other"`
	ContactNumber   string    `json:"contact_number"`
	Email           string    `json:"email" binding:"omitempty,email"`
	Address         string    `json:"address"`
	EmergencyName   string    `json:"emergency_name"`
	EmergencyNumber string    `json:"emergency_number"`
	BloodGroup      string    `json:"blood_group"`
	Allergies       string    `json:"allergies"`
	MedicalHistory  string    `json:"medical_history"`
	CurrentMedication string  `json:"current_medication"`
	Notes           string    `json:"notes"`
}

// UpdatePatientMedicalRequest represents a request to update a patient's medical information by a doctor
type UpdatePatientMedicalRequest struct {
	BloodGroup        string `json:"blood_group"`
	Allergies         string `json:"allergies"`
	MedicalHistory    string `json:"medical_history"`
	CurrentMedication string `json:"current_medication"`
	Notes             string `json:"notes"`
}

// ApplyUpdates applies updates from an UpdatePatientRequest
func (p *Patient) ApplyUpdates(req UpdatePatientRequest) {
	if req.FirstName != "" {
		p.FirstName = req.FirstName
	}
	if req.LastName != "" {
		p.LastName = req.LastName
	}
	if !req.DateOfBirth.IsZero() {
		p.DateOfBirth = req.DateOfBirth
	}
	if req.Gender != "" {
		p.Gender = req.Gender
	}
	if req.ContactNumber != "" {
		p.ContactNumber = req.ContactNumber
	}
	if req.Email != "" {
		p.Email = req.Email
	}
	if req.Address != "" {
		p.Address = req.Address
	}
	p.EmergencyName = req.EmergencyName
	p.EmergencyNumber = req.EmergencyNumber
	p.BloodGroup = req.BloodGroup
	p.Allergies = req.Allergies
	p.MedicalHistory = req.MedicalHistory
	p.CurrentMedication = req.CurrentMedication
	p.Notes = req.Notes
}

// ApplyMedicalUpdates applies medical updates from an UpdatePatientMedicalRequest
func (p *Patient) ApplyMedicalUpdates(req UpdatePatientMedicalRequest) {
	p.BloodGroup = req.BloodGroup
	p.Allergies = req.Allergies
	p.MedicalHistory = req.MedicalHistory
	p.CurrentMedication = req.CurrentMedication
	p.Notes = req.Notes
} 