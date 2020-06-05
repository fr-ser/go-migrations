package postgres

import (
	"database/sql"
	"fmt"
	"time"

	// import to register driver
	_ "github.com/jackc/pgx/stdlib"

	"go-migrations/database"
	"go-migrations/database/config"
)

// variables to allow mocking for tests
var (
	sqlOpen             = sql.Open
	commonWaitForStart  = database.WaitForStart
	commonBootstrap     = database.ApplyBootstrapMigration
	commonGetMigrations = database.GetMigrations
)

// Postgres is a model to apply migrations against a PostgreSQL database
type Postgres struct {
	config        config.Config
	connectionURL string
}

// WaitForStart tries to connect to the database within a timeout
func (pg *Postgres) WaitForStart(pollInterval time.Duration, retryCount int) error {
	db, err := sqlOpen("pgx", pg.connectionURL)
	if err != nil {
		return fmt.Errorf("Error opening database: %v", err)
	}
	defer db.Close()

	return commonWaitForStart(db, pollInterval, retryCount)
}

// Bootstrap applies the bootstrap migration
func (pg *Postgres) Bootstrap() error {
	db, err := sqlOpen("pgx", pg.connectionURL)
	if err != nil {
		return fmt.Errorf("Error opening database: %v", err)
	}
	defer db.Close()

	return commonBootstrap(db, pg.config.MigrationsPath)
}

// applyUpMigration applies the up migration in a transaction
// Depending on the config it also first runs the prepare script
// After the migration a verify script is executed and rolled back in a separate transaction.
// If the verify script fails the downmigration is executed (also in a transaction)
func (pg *Postgres) applyUpMigration(db *sql.DB, migration database.FileMigration) error {
	if err := database.ApplyUpMigration(db, migration); err != nil {
		return err
	}

	if err := database.InsertToChangelog(db, migration, "public.migrations_changelog"); err != nil {
		return err
	}

	if err := database.ApplyVerify(db, migration); err != nil {
		return err
	}

	return nil
}

// ApplyAllUpMigrations applies all up migrations
func (pg *Postgres) ApplyAllUpMigrations() (err error) {
	migrations, err := commonGetMigrations(pg.config.MigrationsPath, nil)
	if err != nil {
		return err
	}
	db, err := sqlOpen("pgx", pg.connectionURL)
	if err != nil {
		return fmt.Errorf("Error opening database: %v", err)
	}
	defer db.Close()

	for _, migration := range migrations {
		err = pg.applyUpMigration(db, migration)
		if err != nil {
			return err
		}
	}
	return nil
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