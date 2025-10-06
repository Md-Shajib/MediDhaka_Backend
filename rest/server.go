package rest

import (
	"fmt"
	"medidhaka/config"
	"medidhaka/repo"
	middleware "medidhaka/rest/middlewares"
	"net/http"
	"os"
	"strconv"
)

func Start(conf config.Config, hospitalRepo repo.HospitalRepo, doctorRepo repo.DoctorRepo, hospitalDoctorRepo repo.HospitalDoctorRepo) {
	manager := middleware.NewManager()
	manager.Use(
		middleware.Preflight,
		middleware.Cors,
		middleware.Logger,
	)

	mux := http.NewServeMux()

	wrappedMux := manager.WrapMux(mux)

	initRoutes(mux, manager, hospitalRepo, doctorRepo, hospitalDoctorRepo)

	port := ":" + strconv.Itoa(conf.HttpPort)
	fmt.Println("Server running on port:", port)

	if err := http.ListenAndServe(port, wrappedMux); err != nil {
		fmt.Println("Error starting the server:", err)
		os.Exit(1)
	}
}
