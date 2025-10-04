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

func Start(conf config.Config, hospitalRepo repo.HospitalRepo) {
	manager := middleware.NewManager()
	manager.Use(
		middleware.Preflight,
		middleware.Cors,
		middleware.Logger,
	)

	mux := http.NewServeMux()

	wrappedMux := manager.WrapMux(mux)

	initRoutes(mux, manager, hospitalRepo)

	port := ":" + strconv.Itoa(conf.HttpPort)
	fmt.Println("Server running on port: ", port)
	err := http.ListenAndServe(port, wrappedMux)
	if err != nil {
		fmt.Println("Error starting the server:", err)
		os.Exit(1)
	}
}
