package database

import (
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/progress"

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
	ApplyAllUpMigrations(pw progress.Writer) error

	// GenerateSeedSQL writes all migration into a single file as an SQL seed
	GenerateSeedSQL(f *os.File) error

	// GetFileMigrations returns the available migrations found locally (sorted by ID)
	GetFileMigrations() ([]FileMigration, error)
	// GetAppliedMigrations gets all applied migrations from the changelog (sorted by ID)
	GetAppliedMigrations() ([]AppliedMigration, error)

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
