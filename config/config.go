package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Version     string
	ServiceName string
	HttpPort    int
}

var (
	config Config
	once   sync.Once
)

func loadConfig() {
	_ = godotenv.Load() // not fatal if missing in production

	version := os.Getenv("VERSION")
	serviceName := os.Getenv("SERVICE_NAME")
	httpPort := os.Getenv("HTTP_PORT")

	if version == "" || serviceName == "" || httpPort == "" {
		fmt.Println("Missing required environment variables")
		os.Exit(1)
	}

	port, err := strconv.Atoi(httpPort)
	if err != nil {
		fmt.Println("HTTP_PORT must be a number")
		os.Exit(1)
	}

	config = Config{
		Version:     version,
		ServiceName: serviceName,
		HttpPort:    port,
	}
}

func GetConfig() Config {
	once.Do(loadConfig)
	return config
}
