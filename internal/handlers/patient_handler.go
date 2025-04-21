package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"healthcare-app/internal/models"
	"healthcare-app/internal/services"

	"github.com/gin-gonic/gin"
)

// PatientHandler handles patient requests
type PatientHandler struct {
	patientService *services.PatientService
}

// NewPatientHandler creates a new PatientHandler
func NewPatientHandler(patientService *services.PatientService) *PatientHandler {
	return &PatientHandler{
		patientService: patientService,
	}
}

// CreatePatient handles create patient requests
// @Summary Create patient
// @Description Create a new patient (Receptionist only)
// @Tags patients
// @Accept json
// @Produce json
// @Param request body models.CreatePatientRequest true "Create Patient Request"
// @Success 201 {object} models.Patient
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /patients [post]
func (h *PatientHandler) CreatePatient(c *gin.Context) {
	var req models.CreatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	userID := GetUserIDFromContext(c)
	patient, err := h.patientService.CreatePatient(req, userID)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, patient)
}

// GetPatient handles get patient requests
// @Summary Get patient
// @Description Get a patient by ID
// @Tags patients
// @Produce json
// @Param id path int true "Patient ID"
// @Success 200 {object} models.Patient
// @Failure 404 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /patients/{id} [get]
func (h *PatientHandler) GetPatient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid patient ID")
		return
	}

	patient, err := h.patientService.GetPatient(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, services.ErrPatientNotFound) {
			status = http.StatusNotFound
		}
		RespondWithError(c, status, err.Error())
		return
	}

	c.JSON(http.StatusOK, patient)
}

// GetAllPatients handles get all patients requests
// @Summary Get all patients
// @Description Get all patients with pagination
// @Tags patients
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} services.PaginationResponse
// @Failure 401 {object} ErrorResponse
// @Router /patients [get]
func (h *PatientHandler) GetAllPatients(c *gin.Context) {
	page, pageSize := GetPaginationParams(c)
	
	patients, err := h.patientService.GetAllPatients(page, pageSize)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, patients)
}

// UpdatePatient handles update patient requests
// @Summary Update patient
// @Description Update a patient (Receptionist only)
// @Tags patients
// @Accept json
// @Produce json
// @Param id path int true "Patient ID"
// @Param request body models.UpdatePatientRequest true "Update Patient Request"
// @Success 200 {object} models.Patient
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /patients/{id} [put]
func (h *PatientHandler) UpdatePatient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid patient ID")
		return
	}

	var req models.UpdatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	patient, err := h.patientService.UpdatePatient(uint(id), req)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, services.ErrPatientNotFound) {
			status = http.StatusNotFound
		}
		RespondWithError(c, status, err.Error())
		return
	}

	c.JSON(http.StatusOK, patient)
}

// UpdatePatientMedicalInfo handles update patient medical info requests
// @Summary Update patient medical info
// @Description Update a patient's medical information (Doctor only)
// @Tags doctor
// @Accept json
// @Produce json
// @Param id path int true "Patient ID"
// @Param request body models.UpdatePatientMedicalRequest true "Update Patient Medical Request"
// @Success 200 {object} models.Patient
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /doctor/patients/{id}/medical [put]
func (h *PatientHandler) UpdatePatientMedicalInfo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid patient ID")
		return
	}

	var req models.UpdatePatientMedicalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	patient, err := h.patientService.UpdatePatientMedicalInfo(uint(id), req)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, services.ErrPatientNotFound) {
			status = http.StatusNotFound
		}
		RespondWithError(c, status, err.Error())
		return
	}

	c.JSON(http.StatusOK, patient)
}

// DeletePatient handles delete patient requests
// @Summary Delete patient
// @Description Delete a patient (Receptionist only)
// @Tags patients
// @Param id path int true "Patient ID"
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /patients/{id} [delete]
func (h *PatientHandler) DeletePatient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid patient ID")
		return
	}

	err = h.patientService.DeletePatient(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, services.ErrPatientNotFound) {
			status = http.StatusNotFound
		}
		RespondWithError(c, status, err.Error())
		return
	}

	RespondWithSuccess(c, "Patient deleted successfully", nil)
}

// SearchPatients handles search patients requests
// @Summary Search patients
// @Description Search patients by search term
// @Tags patients
// @Produce json
// @Param q query string true "Search term"
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} services.PaginationResponse
// @Failure 401 {object} ErrorResponse
// @Router /patients/search [get]
func (h *PatientHandler) SearchPatients(c *gin.Context) {
	searchTerm := c.Query("q")
	page, pageSize := GetPaginationParams(c)
	
	patients, err := h.patientService.SearchPatients(searchTerm, page, pageSize)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, patients)
} 