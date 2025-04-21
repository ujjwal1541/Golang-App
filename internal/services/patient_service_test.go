package services

import (
	"errors"
	"testing"
	"time"

	"healthcare-app/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPatientRepository is a mock implementation of PatientRepository
type MockPatientRepository struct {
	mock.Mock
}

func (m *MockPatientRepository) Create(patient *models.Patient) error {
	args := m.Called(patient)
	return args.Error(0)
}

func (m *MockPatientRepository) FindByID(id uint) (*models.Patient, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Patient), args.Error(1)
}

func (m *MockPatientRepository) FindAll(limit, offset int) ([]models.Patient, int64, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]models.Patient), args.Get(1).(int64), args.Error(2)
}

func (m *MockPatientRepository) Update(patient *models.Patient) error {
	args := m.Called(patient)
	return args.Error(0)
}

func (m *MockPatientRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPatientRepository) SearchPatients(searchTerm string, limit, offset int) ([]models.Patient, int64, error) {
	args := m.Called(searchTerm, limit, offset)
	return args.Get(0).([]models.Patient), args.Get(1).(int64), args.Error(2)
}

func TestCreatePatient(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockPatientRepository)
	
	// Create test data
	dob, _ := time.Parse("2006-01-02", "1990-01-01")
	createReq := models.CreatePatientRequest{
		FirstName:     "John",
		LastName:      "Doe",
		DateOfBirth:   dob,
		Gender:        "male",
		ContactNumber: "1234567890",
		Email:         "john.doe@example.com",
		Address:       "123 Main St",
	}
	
	// Setup expectations
	mockRepo.On("Create", mock.AnythingOfType("*models.Patient")).Return(nil)
	
	// Create service with mock
	service := NewPatientService(mockRepo)
	
	// Test create patient
	patient, err := service.CreatePatient(createReq, 1)
	
	// Assert results
	assert.NoError(t, err)
	assert.NotNil(t, patient)
	assert.Equal(t, "John", patient.FirstName)
	assert.Equal(t, "Doe", patient.LastName)
	assert.Equal(t, dob, patient.DateOfBirth)
	assert.Equal(t, "male", patient.Gender)
	assert.Equal(t, "1234567890", patient.ContactNumber)
	assert.Equal(t, "john.doe@example.com", patient.Email)
	assert.Equal(t, "123 Main St", patient.Address)
	assert.Equal(t, uint(1), patient.RegisteredBy)
	
	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestGetPatient_Success(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockPatientRepository)
	
	// Create test data
	dob, _ := time.Parse("2006-01-02", "1990-01-01")
	patient := &models.Patient{
		ID:            1,
		FirstName:     "John",
		LastName:      "Doe",
		DateOfBirth:   dob,
		Gender:        "male",
		ContactNumber: "1234567890",
		Email:         "john.doe@example.com",
		Address:       "123 Main St",
		RegisteredBy:  1,
	}
	
	// Setup expectations
	mockRepo.On("FindByID", uint(1)).Return(patient, nil)
	
	// Create service with mock
	service := NewPatientService(mockRepo)
	
	// Test get patient
	result, err := service.GetPatient(1)
	
	// Assert results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, patient, result)
	
	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestGetPatient_NotFound(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockPatientRepository)
	
	// Setup expectations
	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("patient not found"))
	
	// Create service with mock
	service := NewPatientService(mockRepo)
	
	// Test get patient
	result, err := service.GetPatient(999)
	
	// Assert results
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrPatientNotFound, err)
	
	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestGetAllPatients(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockPatientRepository)
	
	// Create test data
	dob, _ := time.Parse("2006-01-02", "1990-01-01")
	patients := []models.Patient{
		{
			ID:            1,
			FirstName:     "John",
			LastName:      "Doe",
			DateOfBirth:   dob,
			Gender:        "male",
			ContactNumber: "1234567890",
			Email:         "john.doe@example.com",
			Address:       "123 Main St",
			RegisteredBy:  1,
		},
		{
			ID:            2,
			FirstName:     "Jane",
			LastName:      "Smith",
			DateOfBirth:   dob,
			Gender:        "female",
			ContactNumber: "0987654321",
			Email:         "jane.smith@example.com",
			Address:       "456 Oak St",
			RegisteredBy:  1,
		},
	}
	
	// Setup expectations
	mockRepo.On("FindAll", 10, 0).Return(patients, int64(2), nil)
	
	// Create service with mock
	service := NewPatientService(mockRepo)
	
	// Test get all patients
	result, err := service.GetAllPatients(1, 10)
	
	// Assert results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(2), result.TotalItems)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.PageSize)
	assert.Equal(t, 1, result.TotalPages)
	assert.Equal(t, patients, result.Items)
	
	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestUpdatePatient(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockPatientRepository)
	
	// Create test data
	dob, _ := time.Parse("2006-01-02", "1990-01-01")
	patient := &models.Patient{
		ID:            1,
		FirstName:     "John",
		LastName:      "Doe",
		DateOfBirth:   dob,
		Gender:        "male",
		ContactNumber: "1234567890",
		Email:         "john.doe@example.com",
		Address:       "123 Main St",
		RegisteredBy:  1,
	}
	
	updateReq := models.UpdatePatientRequest{
		FirstName: "Johnny",
		Email:     "johnny.doe@example.com",
	}
	
	// Setup expectations
	mockRepo.On("FindByID", uint(1)).Return(patient, nil)
	mockRepo.On("Update", patient).Return(nil)
	
	// Create service with mock
	service := NewPatientService(mockRepo)
	
	// Test update patient
	result, err := service.UpdatePatient(1, updateReq)
	
	// Assert results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Johnny", result.FirstName) // Updated field
	assert.Equal(t, "Doe", result.LastName)     // Unchanged field
	assert.Equal(t, "johnny.doe@example.com", result.Email) // Updated field
	
	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestDeletePatient(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockPatientRepository)
	
	// Setup expectations
	mockRepo.On("FindByID", uint(1)).Return(&models.Patient{ID: 1}, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)
	
	// Create service with mock
	service := NewPatientService(mockRepo)
	
	// Test delete patient
	err := service.DeletePatient(1)
	
	// Assert results
	assert.NoError(t, err)
	
	// Verify expectations
	mockRepo.AssertExpectations(t)
} 