package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"medidhaka/repo"
	"medidhaka/util"
	"net/http"
	"strconv"
)

// HospitalHandler holds the dependency on the HospitalRepo interface.
type HospitalHandler struct {
	repo repo.HospitalRepo
}

// NewHospitalHandler creates and returns a new HospitalHandler instance.
func NewHospitalHandler(r repo.HospitalRepo) *HospitalHandler {
	return &HospitalHandler{repo: r}
}

// ====================================================================
// CRUD IMPLEMENTATIONS
// ====================================================================

// CreateHospital handles POST requests to create a new Hospital record. (C)
func (h *HospitalHandler) CreateHospital(w http.ResponseWriter, r *http.Request) {
	var hospital repo.Hospital
	if err := json.NewDecoder(r.Body).Decode(&hospital); err != nil {
		util.SendData(w, map[string]string{"error": "Invalid input format"}, http.StatusBadRequest)
		return
	}

	// Basic input validation
	if hospital.Name == "" {
		util.SendData(w, map[string]string{"error": "Hospital name is required"}, http.StatusBadRequest)
		return
	}

	createdHospital, err := h.repo.Create(hospital)
	if err != nil {
		log.Printf("Failed to create hospital: %v", err)
		util.SendData(w, map[string]string{"error": "Failed to create hospital record"}, http.StatusInternalServerError)
		return
	}

	util.SendData(w, createdHospital, http.StatusCreated)
	log.Printf("✅ Hospital created: %s (ID: %d)", createdHospital.Name, createdHospital.HospitalID)
}

// ListHospitals handles GET requests to retrieve a list of all Hospital records. (R - All)
func (h *HospitalHandler) ListHospitals(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	search := query.Get("search")
	page := 1
	limit := 10

	if p := query.Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if l := query.Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}

	offset := (page - 1) * limit
	hospitals, total, err := h.repo.List(search, offset, limit)

	if err != nil {
		log.Printf("Failed to list hospitals: %v", err)
		util.SendData(w, map[string]string{"error": "Internal server error listing hospitals"}, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":       hospitals,
		"limit":      limit,
		"page":       page,
		"total":      total,
		"totalPages": (total + limit - 1) / limit,
	}

	// Returns an empty JSON array if no records are found, which is standard.
	util.SendData(w, response, http.StatusOK)
}

// GetHospital handles GET requests to retrieve a single Hospital by ID. (R - Single)
func (h *HospitalHandler) GetHospital(w http.ResponseWriter, r *http.Request) {
	// 1. Extract ID from URL path using r.PathValue (standard for Go 1.22 mux)
	idStr := r.PathValue("id")
	if idStr == "" {
		util.SendData(w, map[string]string{"error": "Missing hospital ID in URL"}, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.SendData(w, map[string]string{"error": "Invalid hospital ID format"}, http.StatusBadRequest)
		return
	}

	hospital, err := h.repo.Get(id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			util.SendData(w, map[string]string{"error": fmt.Sprintf("Hospital with ID %d not found", id)}, http.StatusNotFound)
			return
		}
		log.Printf("Failed to get hospital ID %d: %v", id, err)
		util.SendData(w, map[string]string{"error": "Internal server error fetching hospital"}, http.StatusInternalServerError)
		return
	}

	util.SendData(w, hospital, http.StatusOK)
}

// UpdateHospital handles PUT requests to modify an existing Hospital record. (U)
func (h *HospitalHandler) UpdateHospital(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		util.SendData(w, map[string]string{"error": "Missing hospital ID in URL"}, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.SendData(w, map[string]string{"error": "Invalid hospital ID format"}, http.StatusBadRequest)
		return
	}

	var hospital repo.Hospital
	if err := json.NewDecoder(r.Body).Decode(&hospital); err != nil {
		util.SendData(w, map[string]string{"error": "Invalid input format"}, http.StatusBadRequest)
		return
	}

	// Ensure the ID from the URL is used for the update operation
	hospital.HospitalID = id

	updatedHospital, err := h.repo.Update(hospital)
	if err != nil {
		if errors.Is(err, repo.ErrFailedUpdate) {
			util.SendData(w, map[string]string{"error": fmt.Sprintf("Hospital with ID %d not found for update", id)}, http.StatusNotFound)
			return
		}
		log.Printf("Failed to update hospital ID %d: %v", id, err)
		util.SendData(w, map[string]string{"error": "Internal server error updating hospital"}, http.StatusInternalServerError)
		return
	}

	util.SendData(w, updatedHospital, http.StatusOK)
	log.Printf("🔄 Hospital updated: %s (ID: %d)", updatedHospital.Name, updatedHospital.HospitalID)
}

// DeleteHospital handles DELETE requests to remove a Hospital record by ID. (D)
func (h *HospitalHandler) DeleteHospital(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		util.SendData(w, map[string]string{"error": "Missing hospital ID in URL"}, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.SendData(w, map[string]string{"error": "Invalid hospital ID format"}, http.StatusBadRequest)
		return
	}

	err = h.repo.Delete(id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			util.SendData(w, map[string]string{"error": fmt.Sprintf("Hospital with ID %d not found for deletion", id)}, http.StatusNotFound)
			return
		}
		log.Printf("Failed to delete hospital ID %d: %v", id, err)
		util.SendData(w, map[string]string{"error": "Internal server error deleting hospital"}, http.StatusInternalServerError)
		return
	}

	util.SendData(w, map[string]string{"message": fmt.Sprintf("Hospital ID %d deleted successfully", id)}, http.StatusOK)
	log.Printf("🗑️ Hospital deleted: ID %d", id)
}
