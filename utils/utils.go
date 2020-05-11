package utils

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// FileExists returns whether the given file or directory exists
func FileExists(path string) (exists bool, err error) {
	_, err = os.Stat(path)
	if err == nil {
		exists = true
	} else if os.IsNotExist(err) {
		exists = false
		err = nil
	}
	return
}

// GetEnvDefault returns the value from the environment for the specified key
// or returns the fallback
func GetEnvDefault(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

// GetEnvOrFail returns the value from the environment or panics with a log
func GetEnvOrFail(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("No environment variable found for: %s", key)
	}
	return value
}
