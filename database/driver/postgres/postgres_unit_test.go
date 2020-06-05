package postgres

import (
	"database/sql"
	"fmt"
	"go-migrations/database"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lithammer/dedent"
)

func TestWaitForStart(t *testing.T) {
	defer resetMockVariables()
	db, mock, err := sqlmock.New()
	mock.ExpectClose()
	sqlOpen = func(a, b string) (*sql.DB, error) { return db, err }

	fakeCalled := false
	commonWaitForStart = func(db *sql.DB, a time.Duration, b int) error {
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
	sqlOpen = func(a, b string) (*sql.DB, error) { return db, err }

	fakeCalled := false
	commonBootstrap = func(db *sql.DB, a string) error {
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

func TestApplyAllUpMigrations(t *testing.T) {
	defer resetMockVariables()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlOpen = func(a, b string) (*sql.DB, error) { return db, err }

	commonGetMigrations = func(a string) (migrations []database.FileMigration, err error) {
		migrations = []database.FileMigration{
			{UpSQL: "SELECT 1", VerifySQL: "SELECT 12", ID: "1", Description: "a"},
			{UpSQL: "SELECT 2", VerifySQL: "SELECT 22", ID: "2", Description: "b"},
		}
		return migrations, err
	}

	mock.ExpectBegin()
	mock.ExpectExec("SELECT 1").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectExec(
		"INSERT INTO public.migrations_changelog(id, name, applied_at) VALUES ('1', 'a', now())",
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectBegin()
	mock.ExpectExec("SELECT 12").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectRollback()

	mock.ExpectBegin()
	mock.ExpectExec("SELECT 2").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectExec(
		"INSERT INTO public.migrations_changelog(id, name, applied_at) VALUES ('2', 'b', now())",
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectBegin()
	mock.ExpectExec("SELECT 22").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectRollback()

	mock.ExpectClose()

	pg := Postgres{}
	err = pg.ApplyAllUpMigrations()
	if err != nil {
		t.Errorf("Expected no error for applying all up migrations, but got: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestApplyAllUpMigrationsUpMigrationError(t *testing.T) {
	defer resetMockVariables()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlOpen = func(a, b string) (*sql.DB, error) { return db, err }

	commonGetMigrations = func(a string) (migrations []database.FileMigration, err error) {
		migrations = []database.FileMigration{{UpSQL: "SELECT 1"}}
		return migrations, err
	}

	mock.ExpectBegin()
	mock.ExpectExec("SELECT 1").WillReturnError(fmt.Errorf("Some error"))
	mock.ExpectRollback()

	mock.ExpectClose()

	pg := Postgres{}
	err = pg.ApplyAllUpMigrations()
	if err == nil {
		t.Errorf("Expected error for applying all up migrations, but got nothing")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestApplyAllUpMigrationsVerifyError(t *testing.T) {
	defer resetMockVariables()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlOpen = func(a, b string) (*sql.DB, error) { return db, err }

	commonGetMigrations = func(a string) (migrations []database.FileMigration, err error) {
		migrations = []database.FileMigration{{UpSQL: "SELECT 1", VerifySQL: "SELECT 12", ID: "1", Description: "a"}}
		return migrations, err
	}

	mock.ExpectBegin()
	mock.ExpectExec("SELECT 1").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectExec(
		"INSERT INTO public.migrations_changelog(id, name, applied_at) VALUES ('1', 'a', now())",
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectBegin()
	mock.ExpectExec("SELECT 12").WillReturnError(fmt.Errorf("Verify error"))
	mock.ExpectRollback()

	mock.ExpectClose()

	pg := Postgres{}
	err = pg.ApplyAllUpMigrations()
	if err == nil {
		t.Errorf("Expected error for applying all up migrations, but got nothing")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestApplyAllUpMigrationsChangelogError(t *testing.T) {
	defer resetMockVariables()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlOpen = func(a, b string) (*sql.DB, error) { return db, err }

	commonGetMigrations = func(a string) (migrations []database.FileMigration, err error) {
		migrations = []database.FileMigration{{UpSQL: "SELECT 1", VerifySQL: "SELECT 12", ID: "1", Description: "a"}}
		return migrations, err
	}

	mock.ExpectBegin()
	mock.ExpectExec("SELECT 1").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectExec(
		"INSERT INTO public.migrations_changelog(id, name, applied_at) VALUES ('1', 'a', now())",
	).WillReturnError(fmt.Errorf("changelog error"))

	mock.ExpectClose()

	pg := Postgres{}
	err = pg.ApplyAllUpMigrations()
	if err == nil {
		t.Errorf("Expected error for applying all up migrations, but got nothing")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEnsureChangelogExists(t *testing.T) {
	defer resetMockVariables()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlOpen = func(a, b string) (*sql.DB, error) { return db, err }

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
	sqlOpen = func(a, b string) (*sql.DB, error) { return db, err }

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
