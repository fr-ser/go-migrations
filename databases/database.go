package databases

import (
	"fmt"
	"time"

	"go-migrations/databases/config"
	"go-migrations/databases/postgres"
)

// variables to allow mocking for tests
var (
	loadConfig = config.LoadConfig
)

// Database is an abstraction over the underlying database and configuration models
type Database interface {
	// WaitForStart tries to connect to the database within a timeout
	WaitForStart(pollInterval time.Duration, retryCount int) error
	// Bootstrap applies the bootstrap migration
	Bootstrap() error
	// ApplyUpMigrations applies all up migrations
	ApplyUpMigrations() error

	// Init initializes the database with the given configuration
	Init(config.Config) error
}

// LoadDb loads a configuration and initializes a database on top of it
func LoadDb(migrationsPath, environment string) (Database, error) {
	configPath := fmt.Sprintf("%s/_environments/%s.yaml", migrationsPath, environment)
	config, err := loadConfig(configPath, migrationsPath, environment)
	if err != nil {
		return nil, err
	}

	db := &postgres.Postgres{}
	if db.Init(config); err != nil {
		return nil, err
	}

	return db, err
}

// Pseudo Code: migrate_up
//
// var db Database
//
// db.load_config(environment)
// db.wait_for_db_to_start()
//
// var appFilter = []string{'sth', 'sth_else'}
// db.get_file_migrations(appFilter)
//
// db.filter_up_migrations(all=false, only="", count=2)
//
// db.apply_up_migrations()
