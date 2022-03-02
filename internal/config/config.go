package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	LISTEN_ADDR string // Address to listen for requests
	BASE_URL    string // Base URL of api

	MAX_TOKEN          int           // Maximum tokens amount to generated offline
	MAX_URL_LEN        int           // Maximum length of url to shrink
	MAX_ALIVE_DURATION time.Duration // Maximum record alive time from creation time

	MONGODB_URI       string // MongoDB URI
	DATABASE          string // Database to be used for shorturl service
	RECORD_COLLECTION string // Collection to store records

	MEMCACHED_ADDRS []string // Memcached addresses
}

var loadedFromEnv = false
var config = Config{
	LISTEN_ADDR: ":8000",
	BASE_URL:    "http://localhost:8000",

	MAX_TOKEN:          5000,
	MAX_URL_LEN:        1024,
	MAX_ALIVE_DURATION: time.Hour * 24 * 365, // an year

	MONGODB_URI:       "mongodb://localhost:27017",
	DATABASE:          "shorturl",
	RECORD_COLLECTION: "records",

	MEMCACHED_ADDRS: []string{"localhost:11211"},
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

		if os.Getenv("MAX_ALIVE_TIME") != "" {
			duration, err := time.ParseDuration(os.Getenv("MAX_ALIVE_TIME"))
			if err != nil {
				panic(err)
			}

			config.MAX_ALIVE_DURATION = duration
		}

		if os.Getenv("DATABASE") != "" {
			config.DATABASE = os.Getenv("DATABASE")
		}

		if os.Getenv("RECORD_COLLECTION") != "" {
			config.RECORD_COLLECTION = os.Getenv("RECORD_COLLECTION")
		}

		if os.Getenv("MEMCACHED_ADDRS") != "" {
			config.MEMCACHED_ADDRS = strings.Split(os.Getenv("MEMCACHED_ADDRS"), ",")
		}

		loadedFromEnv = true
	}

	return &config
}
