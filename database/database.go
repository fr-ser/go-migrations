package database

import (
	"go-migrations/database/config"
	"time"
)

// Database is an abstraction over the underlying database and configuration models
type Database interface {
	// WaitForStart tries to connect to the database within a timeout
	WaitForStart(pollInterval time.Duration, retryCount int) error
	// Bootstrap applies the bootstrap migration
	Bootstrap() error
	// ApplyAllUpMigrations applies all up migrations
	ApplyAllUpMigrations() error

	// Init initializes the database with the given configuration
	Init(config.Config) error
}
