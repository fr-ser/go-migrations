package postgres

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kylelemons/godebug/pretty"
	"github.com/lithammer/dedent"

	"go-migrations/database"
)

func TestWaitForStart(t *testing.T) {
	defer resetMockVariables()
	db, mock, err := sqlmock.New()
	mock.ExpectClose()
	mockableSQLOpen = func(a, b string) (*sql.DB, error) { return db, err }

	fakeCalled := false
	mockableWaitForStart = func(db *sql.DB, a time.Duration, b int) error {
		fakeCalled = true
		return nil
	}

	pg := Postgres{}
	pg.WaitForStart(time.Duration(1), 1)

	if !fakeCalled {
		t.Errorf("Expected WaitForStart to be called")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestBootstrap(t *testing.T) {
	defer resetMockVariables()
	db, mock, err := sqlmock.New()
	mock.ExpectClose()
	mockableSQLOpen = func(a, b string) (*sql.DB, error) { return db, err }

	fakeCalled := false
	mockableBootstrap = func(db *sql.DB, a string) error {
		fakeCalled = true
		return nil
	}

	pg := Postgres{}
	pg.Bootstrap()

	if !fakeCalled {
		t.Errorf("Expected Bootstrap to be called")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEnsureConsistentMigrations(t *testing.T) {
	defer resetMockVariables()
	db, mock, err := sqlmock.New()
	mock.ExpectClose()
	mockableSQLOpen = func(a, b string) (*sql.DB, error) { return db, err }

	var receivedFileMigrations []database.FileMigration
	var receivedAppliedMigrations []database.AppliedMigration
	mockableEnsureConsistentMigrations = func(
		a []database.FileMigration, b []database.AppliedMigration,
	) error {
		receivedFileMigrations = a
		receivedAppliedMigrations = b
		return nil
	}

	expectedFileMigrations := []database.FileMigration{
		{ID: "foo"}, {ID: "bar"},
	}
	mockableGetFileMigrations = func(a string) ([]database.FileMigration, error) {
		return expectedFileMigrations, nil
	}
	expectedAppliedMigrations := []database.AppliedMigration{
		{ID: "foo"}, {ID: "bar"},
	}
	mockableGetAppliedMigrations = func(a *sql.DB, b string) ([]database.AppliedMigration, error) {
		return expectedAppliedMigrations, nil
	}

	pg := Postgres{}
	pg.EnsureConsistentMigrations()

	if receivedAppliedMigrations == nil && receivedFileMigrations == nil {
		t.Errorf("Did not call EnsureConsistentMigrations")
	}
	if diff := pretty.Compare(expectedFileMigrations, receivedFileMigrations); diff != "" {
		t.Errorf("Did not pass right FileMigrations:\n%s", diff)
	}
	if diff := pretty.Compare(expectedAppliedMigrations, receivedAppliedMigrations); diff != "" {
		t.Errorf("Did not pass right AppliedMigrations:\n%s", diff)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEnsureChangelogExists(t *testing.T) {
	defer resetMockVariables()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	mockableSQLOpen = func(a, b string) (*sql.DB, error) { return db, err }

	mock.ExpectQuery(dedent.Dedent(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_schema = 'public'
				AND	table_name = 'migrations_changelog'
		) AS exists
	`)).WillReturnRows(
		sqlmock.NewRows([]string{"exists"}).AddRow(true),
	)
	mock.ExpectClose()

	pg := Postgres{}
	created, err := pg.EnsureMigrationsChangelog()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if created {
		t.Errorf("Expected the created flag to be false, but it was true")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEnsureChangelogNotExists(t *testing.T) {
	defer resetMockVariables()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	mockableSQLOpen = func(a, b string) (*sql.DB, error) { return db, err }

	mock.ExpectQuery(dedent.Dedent(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_schema = 'public'
				AND	table_name = 'migrations_changelog'
		) AS exists
	`)).WillReturnRows(
		sqlmock.NewRows([]string{"exists"}).AddRow(false),
	)
	mock.ExpectExec(dedent.Dedent(`
		CREATE TABLE public.migrations_changelog (
			  id VARCHAR(14) NOT NULL PRIMARY KEY
			, name TEXT NOT NULL
			, applied_at timestamptz NOT NULL
		);
	`)).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectClose()

	pg := Postgres{}
	created, err := pg.EnsureMigrationsChangelog()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if !created {
		t.Errorf("Expected the created flag to be true, but it was false")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPrintStatus(t *testing.T) {
	defer resetMockVariables()
	expectedRows := []database.MigrateStatusRow{{ID: "abc"}}
	expectedStatus := "All good"

	mockableGetMigrationStatus = func(
		fileMigrations []database.FileMigration, appliedMigrations []database.AppliedMigration,
	) (rows []database.MigrateStatusRow, statusNote string, err error) {
		return expectedRows, expectedStatus, nil
	}
	var gotRows []database.MigrateStatusRow
	var gotStatus string
	mockablePrintStatusTable = func(rows []database.MigrateStatusRow, statusNote string) {
		gotRows = rows
		gotStatus = statusNote
	}

	pg := Postgres{}
	pg.fileMigrations = []database.FileMigration{}
	pg.appliedMigrations = []database.AppliedMigration{}
	err := pg.PrintMigrationStatus()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if diff := pretty.Compare(expectedRows, gotRows); diff != "" {
		t.Errorf("Did not pass right rows for print:\n%s", diff)
	}
	if expectedStatus != gotStatus {
		t.Errorf("Expected status of '%s' but got '%s'", expectedStatus, gotStatus)
	}

}
