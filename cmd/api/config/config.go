// Package config provides a way to group all code related to the api
// critical behavior.
package config

import (
	"log"
	"sync"

	"github.com/bryryann/mantel/backend/cmd/api/helpers"
)

type JWT struct {
	Secret   string
	Issuer   string
	Audience string
}

// Configuration holds the values used to setup the application
type Configuration struct {
	Port int    // Port in which the API will be hosted.
	Env  string // Current application environment (DEVELOPMENT, PRODUCTION, etc).
	DSN  string
	JWT  JWT
}

var (
	instance *Configuration
	once     sync.Once
)

// Load returns a singleton Configuration instance, initializing it if necessary
// and defining it's default values.
func Load() *Configuration {
	once.Do(func() {
		// port
		port, err := helpers.GetEnvInt("PORT", 4000)
		if err != nil {
			log.Fatalf("Invalid PORT value: %v", err)
		}

		// db
		dsn := helpers.GetEnvString("MANTEL_DB_DSN", "")
		if dsn == "" {
			log.Fatal("Empty database dsn string\n")
		}

		// jwt
		secret := helpers.GetEnvString("JWT_SECRET", "")
		issuer := helpers.GetEnvString("JWT_ISSUER", "")
		audience := helpers.GetEnvString("JWT_AUDIENCE", "")
		if secret == "" || issuer == "" || audience == "" {
			log.Fatal("Missing jwt secret/issuer/audience\n")
		}

		instance = &Configuration{
			Port: port,
			Env:  helpers.GetEnvString("ENVIRONMENT", "DEVELOPMENT"),
			DSN:  dsn,
			JWT: JWT{
				Secret:   secret,
				Issuer:   issuer,
				Audience: audience,
			},
		}
	})

	return instance
}
