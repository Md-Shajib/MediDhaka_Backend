package rest

import (
	"net/http"

	"medidhaka/repo"
	"medidhaka/rest/handlers"
	middleware "medidhaka/rest/middlewares"
)

func initRoutes(mux *http.ServeMux, manager *middleware.Manager, hospitalRepo repo.HospitalRepo, doctorRepo repo.DoctorRepo, hospitalDoctorRepo repo.HospitalDoctorRepo) {
	// Initialize handlers
	hospitalHandler := handlers.NewHospitalHandler(hospitalRepo)
	doctorHandler := handlers.NewDoctorHandler(doctorRepo)
	hospitalDoctorHandler := handlers.NewHospitalDoctorHandler(hospitalDoctorRepo)

	// ---------- Hospital Routes ----------
	const hospitalPath = "/hospitals"
	const hospitalIDPath = hospitalPath + "/{id}"

	mux.Handle("POST "+hospitalPath, manager.With(http.HandlerFunc(hospitalHandler.CreateHospital)))
	mux.Handle("GET "+hospitalPath, manager.With(http.HandlerFunc(hospitalHandler.ListHospitals)))
	mux.Handle("GET "+hospitalIDPath, manager.With(http.HandlerFunc(hospitalHandler.GetHospital)))
	mux.Handle("PUT "+hospitalIDPath, manager.With(http.HandlerFunc(hospitalHandler.UpdateHospital)))
	mux.Handle("DELETE "+hospitalIDPath, manager.With(http.HandlerFunc(hospitalHandler.DeleteHospital)))

	// ---------- Doctor Routes ----------
	const doctorPath = "/doctors"
	const doctorIDPath = doctorPath + "/{id}"

	mux.Handle("POST "+doctorPath, manager.With(http.HandlerFunc(doctorHandler.CreateDoctor)))
	mux.Handle("GET "+doctorPath, manager.With(http.HandlerFunc(doctorHandler.ListDoctors)))
	mux.Handle("GET "+doctorIDPath, manager.With(http.HandlerFunc(doctorHandler.GetDoctor)))
	mux.Handle("PUT "+doctorIDPath, manager.With(http.HandlerFunc(doctorHandler.UpdateDoctor)))
	mux.Handle("DELETE "+doctorIDPath, manager.With(http.HandlerFunc(doctorHandler.DeleteDoctor)))

	// ---------- Hospital–Doctor Relation ----------
	const hospitalDoctorPath = "/hospital-doctor"

	mux.Handle("POST "+hospitalDoctorPath, manager.With(http.HandlerFunc(hospitalDoctorHandler.AssignDoctor)))
	mux.Handle("GET "+hospitalDoctorPath, manager.With(http.HandlerFunc(hospitalDoctorHandler.ListDoctorsByHospital)))
}
