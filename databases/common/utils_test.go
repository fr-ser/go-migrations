package common

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestWithRunningDb(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("SELECT 1").WillReturnResult(sqlmock.NewResult(1, 1))

	if err = WaitForStart(db, 1000, 15); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestWithStartingDb(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("SELECT 1").WillReturnError(errors.New("sth"))
	mock.ExpectExec("SELECT 1").WillReturnError(errors.New("sth"))
	mock.ExpectExec("SELECT 1").WillReturnResult(sqlmock.NewResult(1, 1))

	if err = WaitForStart(db, 1, 3); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestWithBrokenDb(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("SELECT 1").WillReturnError(errors.New("some error"))

	if err = WaitForStart(db, 1, 3); err == nil {
		t.Error("Expected an error since the db is not up")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
