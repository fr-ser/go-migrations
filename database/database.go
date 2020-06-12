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

	// ApplySpecificUpMigration applies one specific up migration based on a string search
	// of the filename
	ApplySpecificUpMigration(filter string) error
	// ApplyUpMigrationsWithCount applies a number of up migration starting from the last
	// by providing the "all" flag all remaining up migrations are applied
	ApplyUpMigrationsWithCount(count uint, all bool) error
	// ApplySpecificDownMigration applies one specific down migration based on a string search
	// of the filename
	ApplySpecificDownMigration(filter string) error
	// ApplyDownMigrationsWithCount applies a number of down migration starting from the last
	// applied
	// by providing the "all" flag all remaining down migrations, which have been rolled out,
	//  are applied
	ApplyDownMigrationsWithCount(count uint, all bool) error

	// EnsureMigrationsChangelog checks if a changelog table already exists and creates it if
	// necessary
	EnsureMigrationsChangelog() (created bool, err error)
	// EnsureConsistentMigrations checks if all applied migrations exist as local files
	// and if no local migration has been "skipped" (newer migrations applied)
	EnsureConsistentMigrations() error
	// Init initializes the database with the given configuration
	Init(config.Config) error
}
