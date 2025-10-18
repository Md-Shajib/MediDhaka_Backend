package handlers

import (
	"encoding/json"
	"fmt"
	"medidhaka/repo"
	"medidhaka/util"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type HospitalDoctorHandler struct {
	repo repo.HospitalDoctorRepo
}

func NewHospitalDoctorHandler(r repo.HospitalDoctorRepo) *HospitalDoctorHandler {
	return &HospitalDoctorHandler{repo: r}
}

// Assign a doctor with a hospital
func (h *HospitalDoctorHandler) AssignDoctor(w http.ResponseWriter, r *http.Request) {
	var rel repo.HospitalDoctor
	if err := json.NewDecoder(r.Body).Decode(&rel); err != nil {
		util.SendData(w, map[string]string{"error": "Invalid input"}, http.StatusBadRequest)
		return
	}

	if err := h.repo.AssignDoctor(rel); err != nil {
		util.SendData(w, map[string]string{"error": "Failed to assign doctor"}, http.StatusInternalServerError)
		return
	}

	util.SendData(w, map[string]string{"message": "Doctor assigned successfully"}, http.StatusCreated)
}

// List doctors of a hospital
func (h *HospitalDoctorHandler) ListDoctorsByHospital(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idx, ok := vars["id"]
	if !ok {
		util.SendData(w, map[string]string{"error": "Missing hospital ID in URL"}, http.StatusBadRequest)
		return
	}
	id, errConv := strconv.Atoi(idx)
	if errConv != nil {
		util.SendData(w, map[string]string{"error": "Invalid hospital ID format"}, http.StatusBadRequest)
		return
	}
	fmt.Println("Hospital ID:", id)
	doctors, err := h.repo.ListDoctorsByHospital(id)
	if err != nil {
		util.SendData(w, map[string]string{"error": "Failed to fetch doctors"}, http.StatusInternalServerError)
		return
	}

	util.SendData(w, doctors, http.StatusOK)
}
