package handlers

import (
	"encoding/json"
	"errors"
	"medidhaka/repo"
	"medidhaka/util"
	"net/http"
	"strconv"
)

type DoctorHandler struct {
	repo repo.DoctorRepo
}

func NewDoctorHandler(r repo.DoctorRepo) *DoctorHandler {
	return &DoctorHandler{repo: r}
}

func (h *DoctorHandler) CreateDoctor(w http.ResponseWriter, r *http.Request) {
	var doc repo.Doctor
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		util.SendData(w, map[string]string{"error": "Invalid input format"}, http.StatusBadRequest)
		return
	}
	created, err := h.repo.Create(doc)
	if err != nil {
		util.SendData(w, map[string]string{"error": "Failed to create doctor"}, http.StatusInternalServerError)
		return
	}
	util.SendData(w, created, http.StatusCreated)
}

func (h *DoctorHandler) ListDoctors(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.List()
	if err != nil {
		util.SendData(w, map[string]string{"error": "Failed to fetch doctors"}, http.StatusInternalServerError)
		return
	}
	util.SendData(w, list, http.StatusOK)
}

func (h *DoctorHandler) GetDoctor(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.SendData(w, map[string]string{"error": "Invalid ID"}, http.StatusBadRequest)
		return
	}
	doctor, err := h.repo.Get(id)
	if err != nil {
		if errors.Is(err, repo.ErrDoctorNotFound) {
			util.SendData(w, map[string]string{"error": "Doctor not found"}, http.StatusNotFound)
			return
		}
		util.SendData(w, map[string]string{"error": "Server error"}, http.StatusInternalServerError)
		return
	}
	util.SendData(w, doctor, http.StatusOK)
}

func (h *DoctorHandler) UpdateDoctor(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.Atoi(idStr)
	var doc repo.Doctor
	json.NewDecoder(r.Body).Decode(&doc)
	doc.DoctorID = id
	updated, err := h.repo.Update(doc)
	if err != nil {
		util.SendData(w, map[string]string{"error": "Failed to update"}, http.StatusInternalServerError)
		return
	}
	util.SendData(w, updated, http.StatusOK)
}

func (h *DoctorHandler) DeleteDoctor(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.Atoi(idStr)
	err := h.repo.Delete(id)
	if err != nil {
		util.SendData(w, map[string]string{"error": "Failed to delete"}, http.StatusInternalServerError)
		return
	}
	util.SendData(w, map[string]string{"message": "Doctor deleted successfully"}, http.StatusOK)
}
