package postgres

import "go-migrations/databases/config"

// Postgres is a model to apply migrations against a PostgreSQL database
type Postgres struct {
	config config.Config
}

// WaitForStart tries to connect to the database within a timeout
func (db *Postgres) WaitForStart() error {
	return nil
}

// Bootstrap applies the bootstrap migration
func (db *Postgres) Bootstrap() error {
	return nil
}

// ApplyUpMigrations applies all up migrations
func (db *Postgres) ApplyUpMigrations() error {
	return nil
}

// Init initializes the database with the given configuration
func (db *Postgres) Init(config config.Config) error {
	db.config = config
	return nil
}
