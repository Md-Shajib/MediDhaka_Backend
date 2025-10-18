package rest

import (
	"fmt"
	"medidhaka/config"
	"medidhaka/repo"
	middleware "medidhaka/rest/middlewares"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func Start(conf config.Config, hospitalRepo repo.HospitalRepo, doctorRepo repo.DoctorRepo, hospitalDoctorRepo repo.HospitalDoctorRepo) {
	manager := middleware.NewManager()
	manager.Use(
		middleware.Cors,
		middleware.Logger,
	)

	r := mux.NewRouter()

	initRoutes(r, manager, hospitalRepo, doctorRepo, hospitalDoctorRepo)

	handler := manager.WrapMux(r)

	port := ":" + strconv.Itoa(conf.HttpPort)

	fmt.Println("Server running on port:", port)

	if err := http.ListenAndServe(port, handler); err != nil {
		fmt.Println("Error starting the server:", err)
		os.Exit(1)
	}
}
