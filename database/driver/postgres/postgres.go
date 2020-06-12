package postgres

import (
	"database/sql"
	"fmt"
	"time"

	// import to register driver
	_ "github.com/jackc/pgx/stdlib"
	"github.com/lithammer/dedent"
	log "github.com/sirupsen/logrus"

	"go-migrations/database"
	"go-migrations/database/config"
)

var (
	mockableSQLOpen                     = sql.Open
	mockableWaitForStart                = database.WaitForStart
	mockableBootstrap                   = database.ApplyBootstrapMigration
	mockableEnsureConsistentMigrations  = database.EnsureConsistentMigrations
	mockableGetFileMigrations           = database.GetFileMigrations
	mockableGetAppliedMigrations        = database.GetAppliedMigrations
	mockableApplyUpMigration            = database.ApplyUpMigration
	mockableApplyDownMigration          = database.ApplyDownMigration
	mockableFilterUpMigrationsByText    = database.FilterUpMigrationsByText
	mockableFilterDownMigrationsByText  = database.FilterDownMigrationsByText
	mockableFilterUpMigrationsByCount   = database.FilterUpMigrationsByCount
	mockableFilterDownMigrationsByCount = database.FilterDownMigrationsByCount
)

var changelogTable = "public.migrations_changelog"

// Postgres is a model to apply migrations against a PostgreSQL database
type Postgres struct {
	config            config.Config
	connectionURL     string
	fileMigrations    []database.FileMigration
	appliedMigrations []database.AppliedMigration
}

// WaitForStart tries to connect to the database within a timeout
func (pg *Postgres) WaitForStart(pollInterval time.Duration, retryCount int) error {
	db, err := mockableSQLOpen("pgx", pg.connectionURL)
	if err != nil {
		return fmt.Errorf("Error opening database: %v", err)
	}
	defer db.Close()

	return mockableWaitForStart(db, pollInterval, retryCount)
}

// Bootstrap applies the bootstrap migration
func (pg *Postgres) Bootstrap() error {
	db, err := mockableSQLOpen("pgx", pg.connectionURL)
	if err != nil {
		return fmt.Errorf("Error opening database: %v", err)
	}
	defer db.Close()

	return mockableBootstrap(db, pg.config.MigrationsPath)
}

// ApplyAllUpMigrations applies all up migrations
func (pg *Postgres) ApplyAllUpMigrations() (err error) {
	if pg.fileMigrations == nil {
		pg.fileMigrations, err = mockableGetFileMigrations(pg.config.MigrationsPath)
		if err != nil {
			return err
		}

	}
	db, err := mockableSQLOpen("pgx", pg.connectionURL)
	if err != nil {
		return fmt.Errorf("Error opening database: %v", err)
	}
	defer db.Close()

	for _, migration := range pg.fileMigrations {
		err = mockableApplyUpMigration(db, migration, changelogTable)
		if err != nil {
			return err
		}
	}
	log.Infof("Applied %d migrations", len(pg.fileMigrations))
	return nil
}

// ApplySpecificUpMigration applies one up migration by a filter
func (pg *Postgres) ApplySpecificUpMigration(filter string) (err error) {
	if pg.fileMigrations == nil {
		pg.fileMigrations, err = mockableGetFileMigrations(pg.config.MigrationsPath)
		if err != nil {
			return err
		}

	}

	db, err := mockableSQLOpen("pgx", pg.connectionURL)
	if err != nil {
		return fmt.Errorf("Error opening database: %v", err)
	}
	defer db.Close()

	if pg.appliedMigrations == nil {
		pg.appliedMigrations, err = mockableGetAppliedMigrations(db, changelogTable)
		if err != nil {
			return err
		}
	}

	migration, err := mockableFilterUpMigrationsByText(filter, pg.fileMigrations, pg.appliedMigrations)
	if err != nil {
		return err
	}

	err = mockableApplyUpMigration(db, migration, changelogTable)
	if err != nil {
		return err
	}

	return nil
}

// ApplySpecificDownMigration applies one up migration by a filter
func (pg *Postgres) ApplySpecificDownMigration(filter string) (err error) {
	if pg.fileMigrations == nil {
		pg.fileMigrations, err = mockableGetFileMigrations(pg.config.MigrationsPath)
		if err != nil {
			return err
		}

	}

	db, err := mockableSQLOpen("pgx", pg.connectionURL)
	if err != nil {
		return fmt.Errorf("Error opening database: %v", err)
	}
	defer db.Close()

	if pg.appliedMigrations == nil {
		pg.appliedMigrations, err = mockableGetAppliedMigrations(db, changelogTable)
		if err != nil {
			return err
		}
	}

	migration, err := mockableFilterDownMigrationsByText(filter, pg.fileMigrations, pg.appliedMigrations)
	if err != nil {
		return err
	}

	err = mockableApplyDownMigration(db, migration, changelogTable)
	if err != nil {
		return err
	}

	return nil
}

// ApplyUpMigrationsWithCount applies up migration by a count
func (pg *Postgres) ApplyUpMigrationsWithCount(count uint, all bool) (err error) {
	if pg.fileMigrations == nil {
		pg.fileMigrations, err = mockableGetFileMigrations(pg.config.MigrationsPath)
		if err != nil {
			return err
		}

	}

	db, err := mockableSQLOpen("pgx", pg.connectionURL)
	if err != nil {
		return fmt.Errorf("Error opening database: %v", err)
	}
	defer db.Close()

	if pg.appliedMigrations == nil {
		pg.appliedMigrations, err = mockableGetAppliedMigrations(db, changelogTable)
		if err != nil {
			return err
		}
	}

	migrations, err := mockableFilterUpMigrationsByCount(
		count, all, pg.fileMigrations, pg.appliedMigrations,
	)
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		err = mockableApplyUpMigration(db, migration, changelogTable)
		if err != nil {
			return err
		}
	}
	return nil
}

// ApplyDownMigrationsWithCount applies down migration by a count
func (pg *Postgres) ApplyDownMigrationsWithCount(count uint, all bool) (err error) {
	if pg.fileMigrations == nil {
		pg.fileMigrations, err = mockableGetFileMigrations(pg.config.MigrationsPath)
		if err != nil {
			return err
		}

	}

	db, err := mockableSQLOpen("pgx", pg.connectionURL)
	if err != nil {
		return fmt.Errorf("Error opening database: %v", err)
	}
	defer db.Close()

	if pg.appliedMigrations == nil {
		pg.appliedMigrations, err = mockableGetAppliedMigrations(db, changelogTable)
		if err != nil {
			return err
		}
	}

	migrations, err := mockableFilterDownMigrationsByCount(
		count, all, pg.fileMigrations, pg.appliedMigrations,
	)
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		err = mockableApplyDownMigration(db, migration, changelogTable)
		if err != nil {
			return err
		}
	}
	return nil
}

// EnsureMigrationsChangelog creates a migrations changelog if necessary
func (pg *Postgres) EnsureMigrationsChangelog() (created bool, err error) {
	db, err := mockableSQLOpen("pgx", pg.connectionURL)
	if err != nil {
		return false, fmt.Errorf("Error opening database: %v", err)
	}
	defer db.Close()

	existRow := db.QueryRow(dedent.Dedent(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_schema = 'public'
				AND	table_name = 'migrations_changelog'
		) AS exists
	`))

	var exists bool
	err = existRow.Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("Error checking for migrations changelog existence: %v", err)
	}
	if exists {
		return false, nil
	}
	_, err = db.Exec(dedent.Dedent(`
		CREATE TABLE public.migrations_changelog (
			  id VARCHAR(14) NOT NULL PRIMARY KEY
			, name TEXT NOT NULL
			, applied_at timestamptz NOT NULL
		);
	`))
	if err != nil {
		return false, fmt.Errorf("Error creating migrations changelog: %v", err)
	}

	return true, nil
}

// EnsureConsistentMigrations checks for inconsistencies in the changelog
func (pg *Postgres) EnsureConsistentMigrations() (err error) {
	if pg.fileMigrations == nil {
		pg.fileMigrations, err = mockableGetFileMigrations(pg.config.MigrationsPath)
		if err != nil {
			return err
		}

	}

	if pg.appliedMigrations == nil {
		db, err := mockableSQLOpen("pgx", pg.connectionURL)
		if err != nil {
			return fmt.Errorf("Error opening database: %v", err)
		}
		defer db.Close()

		pg.appliedMigrations, err = mockableGetAppliedMigrations(db, changelogTable)
		if err != nil {
			return err
		}

	}

	return mockableEnsureConsistentMigrations(pg.fileMigrations, pg.appliedMigrations)
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
