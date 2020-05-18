package databases

// Database is an abstraction over the underlying database and configuration models
type Database interface {
	// try to connect to the database within a timeout
	WaitForStart() error
	// apply bootstrap migration
	Bootstrap() error
	// apply all up migrations
	ApplyUpMigrations() error
}

// LoadDb loads a configuration and initializes a database on top of it
func LoadDb(migrationsPath, environment string) (Database, error) {
	return nil, nil
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
