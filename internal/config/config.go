package config

import (
	"fmt"
	"os"
)

type Config struct {
	LISTEN_ADDR string // Address to listen for requests
	BASE_URL    string // Base URL of api

	MAX_TOKEN int // Maximum tokens amount to generated offline

	MONGODB_URI string // MongoDB URI
}

var loadedFromEnv = false
var config = Config{
	LISTEN_ADDR: ":8000",
	BASE_URL:    "http://localhost",
	MAX_TOKEN:   1000,
	MONGODB_URI: "mongodb://localhost:27017",
}

// Get config loaded from environment variable or default value
func GetConfig() *Config {
	// Load configs from env only on the first call
	if !loadedFromEnv {
		if os.Getenv("LISTEN_ADDR") != "" {
			config.LISTEN_ADDR = os.Getenv("LISTEN_ADDR")
		}

		if os.Getenv("BASE_URL") != "" {
			config.BASE_URL = os.Getenv("BASE_URL")
		}

		if os.Getenv("MONGODB_URI") != "" {
			config.MONGODB_URI = os.Getenv("MONGODB_URI")
		}

		if os.Getenv("MAX_TOKEN") != "" {
			_, err := fmt.Sscanf(os.Getenv("MAX_TOKEN"), "%d", &config.MAX_TOKEN)
			if err != nil {
				panic(err)
			}
		}

		loadedFromEnv = true
	}

	return &config
}
