package postgres

import (
	"database/sql"
	"fmt"
	"time"

	// import to register driver
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/lithammer/dedent"

	"go-migrations/database"
	"go-migrations/database/config"
	"go-migrations/internal/direction"
)

var (
	mockableSQLOpen                    = sql.Open
	mockableWaitForStart               = database.WaitForStart
	mockableBootstrap                  = database.ApplyBootstrapMigration
	mockableEnsureConsistentMigrations = database.EnsureConsistentMigrations
	mockableGetFileMigrations          = database.GetFileMigrations
	mockableGetAppliedMigrations       = database.GetAppliedMigrations
	mockableApplyMigration             = database.ApplyMigration
	mockableFilterMigrationsByText     = database.FilterMigrationsByText
	mockableFilterMigrationsByCount    = database.FilterMigrationsByCount
)

var changelogTable = "public.migrations_changelog"

var tracker progress.Tracker

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

// GetFileMigrations returns the available migrations found locally (sorted by ID)
func (pg *Postgres) GetFileMigrations() (migrations []database.FileMigration, err error) {
	if pg.fileMigrations != nil {
		return pg.fileMigrations, nil
	}

	pg.fileMigrations, err = mockableGetFileMigrations(pg.config.MigrationsPath)
	return pg.fileMigrations, err
}

// GetAppliedMigrations gets all applied migrations from the changelog (sorted by ID)
func (pg *Postgres) GetAppliedMigrations() (migrations []database.AppliedMigration, err error) {
	if pg.appliedMigrations != nil {
		return pg.appliedMigrations, nil
	}

	db, err := mockableSQLOpen("pgx", pg.connectionURL)
	if err != nil {
		return nil, fmt.Errorf("Error opening database: %v", err)
	}
	defer db.Close()

	pg.appliedMigrations, err = mockableGetAppliedMigrations(db, changelogTable)
	return pg.appliedMigrations, err
}

// ApplyAllUpMigrations applies all up migrations
func (pg *Postgres) ApplyAllUpMigrations(pw progress.Writer) (err error) {
	if pg.fileMigrations == nil {
		_, err = pg.GetFileMigrations()
		if err != nil {
			return err
		}

	}

	db, err := mockableSQLOpen("pgx", pg.connectionURL)
	if err != nil {
		return fmt.Errorf("Error opening database: %v", err)
	}
	defer db.Close()

	tracker = progress.Tracker{
		Message: "Applying migrations",
		Total:   int64(len(pg.fileMigrations)),
	}

	pw.AppendTracker(&tracker)

	for _, migration := range pg.fileMigrations {
		err = mockableApplyMigration(db, migration, changelogTable, direction.Up)
		if err != nil {
			return err
		}
		tracker.Increment(1)
	}
	tracker.MarkAsDone()

	return nil
}

// ApplySpecificMigration applies one migration by a filter
func (pg *Postgres) ApplySpecificMigration(
	filter string, direction direction.MigrateDirection,
) (err error) {
	if pg.fileMigrations == nil {
		_, err = pg.GetFileMigrations()
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
		_, err = pg.GetAppliedMigrations()
		if err != nil {
			return err
		}
	}

	migration, err := mockableFilterMigrationsByText(
		filter, direction, pg.fileMigrations, pg.appliedMigrations,
	)
	if err != nil {
		return err
	}

	err = mockableApplyMigration(db, migration, changelogTable, direction)
	if err != nil {
		return err
	}

	return nil
}

// ApplyMigrationsWithCount applies up migration by a count
func (pg *Postgres) ApplyMigrationsWithCount(
	count uint, all bool, dir direction.MigrateDirection,
) (err error) {
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
		_, err = pg.GetAppliedMigrations()
		if err != nil {
			return err
		}
	}

	migrations, err := mockableFilterMigrationsByCount(
		count, all, dir, pg.fileMigrations, pg.appliedMigrations,
	)
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		err = mockableApplyMigration(db, migration, changelogTable, dir)
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
		_, err = pg.GetFileMigrations()
		if err != nil {
			return err
		}

	}

	if pg.appliedMigrations == nil {
		_, err = pg.GetAppliedMigrations()
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
