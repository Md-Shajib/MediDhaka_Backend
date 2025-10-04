package rest

import (
	"net/http"

	"medidhaka/repo"
	"medidhaka/rest/handlers"
	middleware "medidhaka/rest/middlewares"
)

func initRoutes(mux *http.ServeMux, manager *middleware.Manager, hospitalRepo repo.HospitalRepo) {

	hospitalHandler := handlers.NewHospitalHandler(hospitalRepo)

	const hospitalPath = "/v1/hospitals"
	const hospitalIDPath = hospitalPath + "/{id}"

	// Create Hospital
	mux.Handle("POST "+hospitalPath, manager.With(http.HandlerFunc(hospitalHandler.CreateHospital)))

	// All Hospitals
	mux.Handle("GET "+hospitalPath, manager.With(http.HandlerFunc(hospitalHandler.ListHospitals)))

	// Single Hospital
	mux.Handle("GET "+hospitalIDPath, manager.With(http.HandlerFunc(hospitalHandler.GetHospital)))

	// Update Hospital
	mux.Handle("PUT "+hospitalIDPath, manager.With(http.HandlerFunc(hospitalHandler.UpdateHospital)))

	// Delete Hospital
	mux.Handle("DELETE "+hospitalIDPath, manager.With(http.HandlerFunc(hospitalHandler.DeleteHospital)))
}
