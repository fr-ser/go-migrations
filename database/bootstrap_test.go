package database

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestApplyBootstrap(t *testing.T) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)
	ioutil.WriteFile(filepath.Join(dir, "bootstrap.sql"), []byte("SELECT 1"), 0777)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("SELECT 1").WillReturnResult(sqlmock.NewResult(1, 1))

	err = ApplyBootstrapMigration(db, dir)
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

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	expectedError := errors.New("my-err")
	mock.ExpectExec("SELECT 1").WillReturnError(expectedError)

	err = ApplyBootstrapMigration(db, dir)
	if err != expectedError {
		t.Fatalf("Received different error during bootstrap: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}
