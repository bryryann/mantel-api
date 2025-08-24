// Package responses provides functionality for managing responses with logging capabilities.
package responses

import (
	"log/slog"
	"os"
	"sync"
)

// Responses struct encapsulates a logger for structured logging.
type Responses struct {
	logger *slog.Logger
}

var (
	// Singleton instance of Responses.
	instance *Responses
	// Ensures the singleton instance is initialized only once.
	once sync.Once
)

// Get initializes and returns the singleton instance of Responses.
// It uses sync.Once to ensure thread-safe lazy initialization.
func Get() *Responses {
	once.Do(func() {
		instance = &Responses{
			// Creating a new logger with a text handler that writes to standard output.
			logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
		}
	})

	return instance
}
