package handlers

import (
	"encoding/json"
	"errors"
	"medidhaka/repo"
	"medidhaka/util"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

	list, total, err := h.repo.List(search, offset, limit)
	if err != nil {
		util.SendData(w, map[string]string{"error": "Failed to fetch doctors"}, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":       list,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + limit - 1) / limit,
	}

	util.SendData(w, response, http.StatusOK)
}

func (h *DoctorHandler) GetDoctor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idx, ok := vars["id"]
	if !ok {
		util.SendData(w, map[string]string{"error": "Missing doctor ID in URL"}, http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idx)
	if err != nil {
		util.SendData(w, map[string]string{"error": "Invalid doctor ID format"}, http.StatusBadRequest)
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
	vars := mux.Vars(r)
	idx, ok := vars["id"]
	if !ok {
		util.SendData(w, map[string]string{"error": "Missing doctor ID in URL"}, http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idx)
	if err != nil {
		util.SendData(w, map[string]string{"error": "Invalid doctor ID format"}, http.StatusBadRequest)
		return
	}
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
	vars := mux.Vars(r)
	idx, ok := vars["id"]
	if !ok {
		util.SendData(w, map[string]string{"error": "Missing doctor ID in URL"}, http.StatusBadRequest)
		return
	}
	id, err_conversion := strconv.Atoi(idx)
	if err_conversion != nil {
		util.SendData(w, map[string]string{"error": "Invalid doctor ID format"}, http.StatusBadRequest)
		return
	}
	err := h.repo.Delete(id)
	if err != nil {
		util.SendData(w, map[string]string{"error": "Failed to delete"}, http.StatusInternalServerError)
		return
	}
	util.SendData(w, map[string]string{"message": "Doctor deleted successfully"}, http.StatusOK)
}
