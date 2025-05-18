// Package config provides a way to group all code related to the api
// critical behavior.
package config

import (
	"log"
	"os"
	"strconv"
	"sync"
)

// Configuration holds the values used to setup the application
type Configuration struct {
	Port int    // Port in which the API will be hosted.
	Env  string // Current application environment (DEVELOPMENT, PRODUCTION, etc).
}

var (
	instance *Configuration
	once     sync.Once
)

// Load returns a singleton Configuration instance, initializing it if necessary
// and defining it's default values.
func Load() *Configuration {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Invalid PORT value: %v", err)
	}

	once.Do(func() {
		instance = &Configuration{
			Port: port,
			Env:  os.Getenv("ENVIRONMENT"),
		}
	})

	return instance
}
