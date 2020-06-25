package database

import (
	"time"

	"go-migrations/database/config"
	"go-migrations/internal/direction"
)

// Database is an abstraction over the underlying database and configuration models
type Database interface {
	// WaitForStart tries to connect to the database within a timeout
	WaitForStart(pollInterval time.Duration, retryCount int) error
	// Bootstrap applies the bootstrap migration
	Bootstrap() error
	// ApplyAllUpMigrations applies all up migrations
	ApplyAllUpMigrations() error

	// PrintMigrationStatus prints a human readable table about applied and unapplied migrations
	PrintMigrationStatus() error

	// ApplySpecificMigration applies one migration based on a string search of the filename
	ApplySpecificMigration(filter string, direction direction.MigrateDirection) error
	// ApplyUpMigrationsWithCount applies a number of up migration starting from the last
	// by providing the "all" flag all remaining up migrations are applied
	ApplyMigrationsWithCount(count uint, all bool, direction direction.MigrateDirection) error

	// EnsureMigrationsChangelog checks if a changelog table already exists and creates it if
	// necessary
	EnsureMigrationsChangelog() (created bool, err error)
	// EnsureConsistentMigrations checks if all applied migrations exist as local files
	// and if no local migration has been "skipped" (newer migrations applied)
	EnsureConsistentMigrations() error
	// Init initializes the database with the given configuration
	Init(config.Config) error
}
