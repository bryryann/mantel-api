package helpers

import (
	"fmt"
	"os"
	"strconv"
)

func GetEnvString(key, defaultValue string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return val
}

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
