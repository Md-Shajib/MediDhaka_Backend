package handlers

import (
	"medidhaka/repo"
	"medidhaka/util"
	"net/http"
	"strings"
)

type SearchHandler struct {
	doctorRepo   repo.DoctorRepo
	hospitalRepo repo.HospitalRepo
}

func NewSearchHandler(dRepo repo.DoctorRepo, hRepo repo.HospitalRepo) *SearchHandler {
	return &SearchHandler{
		doctorRepo:   dRepo,
		hospitalRepo: hRepo,
	}
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		util.SendData(w, map[string]string{"error": "Search query is required"}, http.StatusBadRequest)
		return
	}

	// Fetch up to 3 doctors and hospitals
	doctors, _, err1 := h.doctorRepo.List(query, 0, 3)
	hospitals, _, err2 := h.hospitalRepo.List(query, 0, 3)

	if err1 != nil || err2 != nil {
		util.SendData(w, map[string]string{"error": "Failed to fetch search results"}, http.StatusInternalServerError)
		return
	}

	// Map results to simplified response
	type item struct {
		Name  string `json:"name"`
		Image string `json:"image"`
		ID    int    `json:"id"`
	}

	var doctorList []item
	for _, d := range doctors {
		doctorList = append(doctorList, item{
			Name:  d.Name,
			Image: d.ImageURL,
			ID:    d.DoctorID,
		})
	}

	var hospitalList []item
	for _, hsp := range hospitals {
		hospitalList = append(hospitalList, item{
			Name:  hsp.Name,
			Image: hsp.ImageURL,
			ID:    hsp.HospitalID,
		})
	}

	response := map[string]interface{}{
		"data": map[string]interface{}{
			"hospital": hospitalList,
			"doctor":   doctorList,
		},
	}

	util.SendData(w, response, http.StatusOK)
}
