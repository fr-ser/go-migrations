package database

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
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

	expectedError := errors.New("my-err")
	mock.ExpectExec("SELECT 1").WillReturnError(expectedError)

	err := ApplyBootstrapMigration(db, dir)
	if err != expectedError {
		t.Fatalf("Received different error during bootstrap: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}
