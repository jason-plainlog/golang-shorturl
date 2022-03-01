package config

import (
	"fmt"
	"os"
)

type Config struct {
	LISTEN_ADDR string // Address to listen for requests
	BASE_URL    string // Base URL of api

	MONGODB_URI string // MongoDB URI
}

var loadedFromEnv = false
var config = Config{
	LISTEN_ADDR: ":8000",
	BASE_URL:    "http://localhost",
}

// Get config loaded from environment variable or default value
func GetConfig() *Config {
	// Load configs from env only on the first call
	if !loadedFromEnv {
		if os.Getenv("LISTEN_ADDR") != "" {
			_, err := fmt.Sscanf(os.Getenv("LISTEN_ADDR"), "%d", &config.LISTEN_ADDR)
			if err != nil {
				panic(err)
			}
		}

		if os.Getenv("BASE_URL") != "" {
			config.BASE_URL = os.Getenv("BASE_URL")
		}

		if os.Getenv("MONGODB_URI") != "" {
			config.MONGODB_URI = os.Getenv("MONGODB_URI")
		}

		loadedFromEnv = true
	}

	return &config
}
