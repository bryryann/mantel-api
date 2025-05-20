package helpers

import (
	"fmt"
	"os"
	"strconv"
)

// GetEnvString returns the env variable with key. If
// variable doesn't exist, returns defaultValue.
func GetEnvString(key, defaultValue string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return val
}

// GetEnvInt returns an env variable with key, and parses it as an int.
// Returns defaultValue if variable is not found.
func GetEnvInt(key string, defaultValue int) (int, error) {
	val, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue, nil
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("environment variable %s is not a valid integer: %v", key, err)
	}

	return i, nil
}
