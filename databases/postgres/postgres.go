package postgres

import (
	"database/sql"
	"fmt"

	// import to register driver
	_ "github.com/jackc/pgx/stdlib"

	"go-migrations/databases/common"
	"go-migrations/databases/config"
)

// variables to allow mocking for tests
var (
	waitForDb = common.WaitForStart
)

// Postgres is a model to apply migrations against a PostgreSQL database
type Postgres struct {
	config        config.Config
	connectionURL string
}

// WaitForStart tries to connect to the database within a timeout
func (pg *Postgres) WaitForStart() error {
	db, err := sql.Open("pgx", pg.connectionURL)
	if err != nil {
		return fmt.Errorf("Error opening database: %v", err)
	}
	defer db.Close()

	return waitForDb(db, 1000, 10)
}

// Bootstrap applies the bootstrap migration
func (pg *Postgres) Bootstrap() error {
	db, err := sql.Open("pgx", pg.connectionURL)
	if err != nil {
		return fmt.Errorf("Error opening database: %v", err)
	}
	defer db.Close()

	return common.ApplyBootstrapMigration(db, pg.config.MigrationsPath)
}

// ApplyUpMigrations applies all up migrations
func (pg *Postgres) ApplyUpMigrations() error {
	return fmt.Errorf("Not implemented")
}

// Init initializes the database with the given configuration
func (pg *Postgres) Init(config config.Config) error {
	pg.config = config
	pg.connectionURL = fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		config.Db.User, config.Db.Password, config.Db.Host, config.Db.Port, config.Db.Name,
	)
	return nil
}
