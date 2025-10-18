package rest

import (
	"net/http"

	"medidhaka/repo"
	"medidhaka/rest/handlers"
	middleware "medidhaka/rest/middlewares"

	"github.com/gorilla/mux"
)

func initRoutes(r *mux.Router, manager *middleware.Manager, hospitalRepo repo.HospitalRepo, doctorRepo repo.DoctorRepo, hospitalDoctorRepo repo.HospitalDoctorRepo) {
	// Initialize handlers
	hospitalHandler := handlers.NewHospitalHandler(hospitalRepo)
	doctorHandler := handlers.NewDoctorHandler(doctorRepo)
	hospitalDoctorHandler := handlers.NewHospitalDoctorHandler(hospitalDoctorRepo)

	// ---------- Hospital Routes ----------
	r.Handle("/hospitals", manager.With(http.HandlerFunc(hospitalHandler.CreateHospital))).Methods("POST", "OPTIONS")
	r.Handle("/hospitals", manager.With(http.HandlerFunc(hospitalHandler.ListHospitals))).Methods("GET", "OPTIONS")
	r.Handle("/hospitals/{id}", manager.With(http.HandlerFunc(hospitalHandler.GetHospital))).Methods("GET", "OPTIONS")
	r.Handle("/hospitals/{id}", manager.With(http.HandlerFunc(hospitalHandler.UpdateHospital))).Methods("PUT", "OPTIONS")
	r.Handle("/hospitals/{id}", manager.With(http.HandlerFunc(hospitalHandler.DeleteHospital))).Methods("DELETE", "OPTIONS")

	// ---------- Doctor Routes ----------
	r.Handle("/doctors", manager.With(http.HandlerFunc(doctorHandler.CreateDoctor))).Methods("POST", "OPTIONS")
	r.Handle("/doctors", manager.With(http.HandlerFunc(doctorHandler.ListDoctors))).Methods("GET", "OPTIONS")
	r.Handle("/doctors/{id}", manager.With(http.HandlerFunc(doctorHandler.GetDoctor))).Methods("GET", "OPTIONS")
	r.Handle("/doctors/{id}", manager.With(http.HandlerFunc(doctorHandler.UpdateDoctor))).Methods("PUT", "OPTIONS")
	r.Handle("/doctors/{id}", manager.With(http.HandlerFunc(doctorHandler.DeleteDoctor))).Methods("DELETE", "OPTIONS")

	// ---------- Hospital–Doctor Relation ----------
	r.Handle("/hospital-doctor", manager.With(http.HandlerFunc(hospitalDoctorHandler.AssignDoctor))).Methods("POST", "OPTIONS")
	r.Handle("/hospital-doctor/{id}", manager.With(http.HandlerFunc(hospitalDoctorHandler.ListDoctorsByHospital))).Methods("GET", "OPTIONS")
}
