package postgres

import (
	"database/sql"
	"fmt"
	"go-migrations/database/common"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
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
	db, mock, err := sqlmock.New()
	sqlOpen = func(a, b string) (*sql.DB, error) { return db, err }

	commonGetMigrations = func(a string, b []string) (migrations []common.FileMigration, err error) {
		migrations = []common.FileMigration{
			{UpSQL: "SELECT 1", VerifySQL: "SELECT 12"},
			{UpSQL: "SELECT 2", VerifySQL: "SELECT 22"},
		}
		return migrations, err
	}

	mock.ExpectBegin()
	mock.ExpectExec("SELECT 1").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectBegin()
	mock.ExpectExec("SELECT 12").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectRollback()

	mock.ExpectBegin()
	mock.ExpectExec("SELECT 2").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
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
	db, mock, err := sqlmock.New()
	sqlOpen = func(a, b string) (*sql.DB, error) { return db, err }

	commonGetMigrations = func(a string, b []string) (migrations []common.FileMigration, err error) {
		migrations = []common.FileMigration{{UpSQL: "SELECT 1"}}
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
	db, mock, err := sqlmock.New()
	sqlOpen = func(a, b string) (*sql.DB, error) { return db, err }

	commonGetMigrations = func(a string, b []string) (migrations []common.FileMigration, err error) {
		migrations = []common.FileMigration{{UpSQL: "SELECT 1", VerifySQL: "SELECT 12"}}
		return migrations, err
	}

	mock.ExpectBegin()
	mock.ExpectExec("SELECT 1").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
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
