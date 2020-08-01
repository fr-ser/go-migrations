package database

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lithammer/dedent"
)

func TestWaitWithRunningDb(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectExec("SELECT 1").WillReturnResult(sqlmock.NewResult(1, 1))

	err := WaitForStart(db, 1000*time.Millisecond, 15)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestWaitWithStartingDb(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectExec("SELECT 1").WillReturnError(errors.New("sth"))
	mock.ExpectExec("SELECT 1").WillReturnError(errors.New("sth"))
	mock.ExpectExec("SELECT 1").WillReturnResult(sqlmock.NewResult(1, 1))

	err := WaitForStart(db, 1*time.Millisecond, 3)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestWaitWithBrokenDb(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectExec("SELECT 1").WillReturnError(errors.New("some error"))

	err := WaitForStart(db, 1*time.Millisecond, 3)
	if err == nil {
		t.Error("Expected an error since the db is not up")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestApplyBootstrap(t *testing.T) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)
	ioutil.WriteFile(filepath.Join(dir, "bootstrap.sql"), []byte("SELECT 1"), 0777)

	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectExec("SELECT 1").WillReturnResult(sqlmock.NewResult(1, 1))

	err := ApplyBootstrapMigration(db, dir)
	if err != nil {
		t.Fatalf("Received error during bootstrap: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}

func TestApplyBootstrapNoFile(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	err := ApplyBootstrapMigration(db, ".")
	if err != nil {
		t.Fatalf("Received error during bootstrap: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}

func TestApplyBootstrapFailure(t *testing.T) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)
	ioutil.WriteFile(filepath.Join(dir, "bootstrap.sql"), []byte("SELECT 1"), 0777)

	db, mock, _ := sqlmock.New()
	defer db.Close()

	expectedSQLErr := errors.New("my-err")
	mock.ExpectExec("SELECT 1").WillReturnError(expectedSQLErr)

	err := ApplyBootstrapMigration(db, dir)
	expectedError := fmt.Sprintf("Could not apply bootstrap.sql: %v", expectedSQLErr)
	if err.Error() != expectedError {
		t.Fatalf("Received different error during bootstrap: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}

func TestEnsureConsistentMigrationsSuccess(t *testing.T) {
	migrations := []struct {
		file    []FileMigration
		applied []AppliedMigration
	}{
		{file: []FileMigration{{ID: "a"}}, applied: []AppliedMigration{{ID: "a"}}},
		{file: []FileMigration{{ID: "a"}, {ID: "b"}}, applied: []AppliedMigration{{ID: "a"}}},
		{
			file:    []FileMigration{{ID: "a"}, {ID: "b"}, {ID: "c"}},
			applied: []AppliedMigration{{ID: "a"}, {ID: "b"}},
		},
	}
	for idx, migration := range migrations {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			err := EnsureConsistentMigrations(migration.file, migration.applied)
			if err != nil {
				t.Fatalf("Received error during EnsureConsistentMigrations: %v", err)
			}
		})
	}
}

func TestEnsureConsistentMigrationsError(t *testing.T) {
	migrations := []struct {
		file    []FileMigration
		applied []AppliedMigration
	}{
		{file: []FileMigration{}, applied: []AppliedMigration{{ID: "a"}}},
		{file: []FileMigration{{ID: "a"}}, applied: []AppliedMigration{{ID: "b"}}},
		{
			file:    []FileMigration{{ID: "a"}, {ID: "b"}, {ID: "c"}},
			applied: []AppliedMigration{{ID: "a"}, {ID: "c"}},
		},
	}
	for idx, migration := range migrations {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			err := EnsureConsistentMigrations(migration.file, migration.applied)
			if err == nil {
				t.Fatalf(
					dedent.Dedent(`
						Received no error during EnsureConsistentMigrations.
						FileMigrations: %v
						AppliedMigrations: %v
					`),
					migration.file, migration.applied,
				)
			}
		})
	}
}
