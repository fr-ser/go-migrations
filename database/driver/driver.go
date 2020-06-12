package driver

import (
	"fmt"

	"go-migrations/database"
	"go-migrations/database/config"
	"go-migrations/database/driver/postgres"
)

// variables to allow mocking for tests
var (
	loadConfig = config.LoadConfig
)

// LoadDB loads a configuration and initializes a database on top of it
func LoadDB(migrationsPath, environment string) (database.Database, error) {
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
