package cmd

import (
	"fmt"
	"medidhaka/config"
	"medidhaka/infra/db"
	"medidhaka/repo"
	"medidhaka/rest"
	"os"
)

func Serve() {
	conf := config.GetConfig()

	dbCon, err := db.NewConnection()
	if err != nil {
		fmt.Println("Database connection failed: ", err)
		os.Exit(1)
	}
	defer dbCon.Close() // connection closed

	hospitalRepo := repo.NewHospitalRepo(dbCon)

	rest.Start(conf, hospitalRepo)
}
