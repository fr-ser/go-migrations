package database

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kylelemons/godebug/pretty"
)

func TestApplyUpMigration(t *testing.T) {
	expectedMigration := FileMigration{
		UpSQL: "SELECT 1", VerifySQL: "SELECT 12", ID: "1", Description: "a",
	}

	var migrateUpCall FileMigration
	mockableMigrateUp = func(a *sql.DB, b FileMigration) error {
		migrateUpCall = b
		return nil
	}

	var insertToChangelogCall FileMigration
	var calledChangelogTable string
	mockableInsertToChangelog = func(a *sql.DB, b FileMigration, c string) error {
		insertToChangelogCall = b
		calledChangelogTable = c
		return nil
	}

	var verifyCall FileMigration
	mockableApplyVerify = func(a *sql.DB, b FileMigration) error {
		verifyCall = b
		return nil
	}

	db, _, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	err := ApplyUpMigration(db, expectedMigration, "schema.sth")
	if err != nil {
		t.Errorf("Expected no error for applying up migrations, but got: %s", err)
	}

	if diff := pretty.Compare(expectedMigration, migrateUpCall); diff != "" {
		t.Errorf("Did not pass right FileMigrations to migrateUp:\n%s", diff)
	}

	if diff := pretty.Compare(expectedMigration, insertToChangelogCall); diff != "" {
		t.Errorf("Did not pass right FileMigrations to insertToChangelog:\n%s", diff)
	}
	if calledChangelogTable != "schema.sth" {
		t.Errorf(
			"Expected passed changelogTable to be %s, but got %s",
			"schema.sth", calledChangelogTable,
		)
	}
	if diff := pretty.Compare(expectedMigration, verifyCall); diff != "" {
		t.Errorf("Did not pass right FileMigrations to verify:\n%s", diff)
	}
}

func TestApplyUpMigrationUpMigrationError(t *testing.T) {
	var migrateUpCalled bool
	mockableMigrateUp = func(a *sql.DB, b FileMigration) error {
		migrateUpCalled = true
		return fmt.Errorf("test error")
	}

	var insertToChangelogCalled bool
	mockableInsertToChangelog = func(a *sql.DB, b FileMigration, c string) error {
		insertToChangelogCalled = true
		return nil
	}

	var verifyCalled bool
	mockableApplyVerify = func(a *sql.DB, b FileMigration) error {
		verifyCalled = true
		return nil
	}

	db, _, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	err := ApplyUpMigration(db, FileMigration{}, "sth")
	if err == nil {
		t.Errorf("Expected error for applying up migrations, but got nothing")
	}

	if !migrateUpCalled {
		t.Errorf("Expected migrateUp to be called once")
	}
	if insertToChangelogCalled {
		t.Errorf("Expected insertToChangelog not to be called")
	}
	if verifyCalled {
		t.Errorf("Expected verify not to be called")
	}
}

func TestApplyUpMigrationVerifyError(t *testing.T) {

	var migrateUpCalled bool
	mockableMigrateUp = func(a *sql.DB, b FileMigration) error {
		migrateUpCalled = true
		return nil
	}

	var insertToChangelogCalled bool
	mockableInsertToChangelog = func(a *sql.DB, b FileMigration, c string) error {
		insertToChangelogCalled = true
		return nil
	}

	var verifyCalled bool
	mockableApplyVerify = func(a *sql.DB, b FileMigration) error {
		verifyCalled = true
		return fmt.Errorf("test error")
	}

	db, _, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	err := ApplyUpMigration(db, FileMigration{}, "sth")
	if err == nil {
		t.Errorf("Expected error for applying up migrations, but got nothing")
	}

	if !migrateUpCalled {
		t.Errorf("Expected migrateUp to be called once")
	}
	if !insertToChangelogCalled {
		t.Errorf("Expected insertToChangelog to be called once")
	}
	if !verifyCalled {
		t.Errorf("Expected verify to be called once")
	}
}

func TestApplyUpMigrationChangelogError(t *testing.T) {
	var migrateUpCalled bool
	mockableMigrateUp = func(a *sql.DB, b FileMigration) error {
		migrateUpCalled = true
		return nil
	}

	var insertToChangelogCalled bool
	mockableInsertToChangelog = func(a *sql.DB, b FileMigration, c string) error {
		insertToChangelogCalled = true
		return fmt.Errorf("test error")
	}

	var verifyCalled bool
	mockableApplyVerify = func(a *sql.DB, b FileMigration) error {
		verifyCalled = true
		return nil
	}

	db, _, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	err := ApplyUpMigration(db, FileMigration{}, "sth")
	if err == nil {
		t.Errorf("Expected error for applying up migrations, but got nothing")
	}

	if !migrateUpCalled {
		t.Errorf("Expected migrateUp to be called once")
	}
	if !insertToChangelogCalled {
		t.Errorf("Expected insertToChangelog  to be called once")
	}
	if verifyCalled {
		t.Errorf("Expected verify not to be called")
	}
}

func TestApplyUpSQL(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	migration := FileMigration{
		UpSQL: "SELECT 1", VerifySQL: "SELECT 12", ID: "1", Description: "a",
	}

	mock.ExpectBegin()
	mock.ExpectExec("SELECT 1").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = ApplyUpSQL(db, migration)
	if err != nil {
		t.Errorf("Expected no error, but got: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestApplyUpSQLMigrationError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	migration := FileMigration{UpSQL: "SELECT 1"}

	mock.ExpectBegin()
	mock.ExpectExec("SELECT 1").WillReturnError(fmt.Errorf("Some error"))
	mock.ExpectRollback()

	err = ApplyUpSQL(db, migration)
	if err == nil {
		t.Errorf("Expected error, but got nothing")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsertToChangelog(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	migration := FileMigration{ID: "1", Description: "a"}

	mock.ExpectExec(
		`INSERT INTO sth (id, name, applied_at) VALUES ('1', 'a', now())`,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	err = InsertToChangelog(db, migration, "sth")
	if err != nil {
		t.Errorf("Expected no error, but got: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestApplyVerify(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	migration := FileMigration{
		UpSQL: "SELECT 1", VerifySQL: "SELECT 12", ID: "1", Description: "a",
	}

	mock.ExpectBegin()
	mock.ExpectExec("SELECT 12").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectRollback()

	err = ApplyVerify(db, migration)
	if err != nil {
		t.Errorf("Expected no error, but got: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestApplyVerifyError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	migration := FileMigration{UpSQL: "SELECT 1", VerifySQL: "SELECT 12", ID: "1", Description: "a"}

	mock.ExpectBegin()
	mock.ExpectExec("SELECT 12").WillReturnError(fmt.Errorf("Verify error"))
	mock.ExpectRollback()

	err = ApplyVerify(db, migration)
	if err == nil {
		t.Errorf("Expected error, but got nothing")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFilterUpMigrationsByText(t *testing.T) {
	testCases := []struct {
		filter            string
		expectedFilename  string
		fileMigrations    []FileMigration
		appliedMigrations []AppliedMigration
	}{
		{
			filter: "2017", expectedFilename: "20171101000003_buz.sql",
			fileMigrations: []FileMigration{
				{ID: "20161101000001", Filename: "20161101000001_foo.sql"},
				{ID: "20171101000003", Filename: "20171101000003_buz.sql"},
			},
			appliedMigrations: []AppliedMigration{},
		},
		{
			filter: "01_start", expectedFilename: "20161101000001_start_a.sql",
			fileMigrations: []FileMigration{
				{ID: "20161101000001", Filename: "20161101000001_start_a.sql"},
				{ID: "20171101000003", Filename: "20171101000003_start_b.sql"},
			},
			appliedMigrations: []AppliedMigration{},
		},
		{
			filter: "2016", expectedFilename: "20161101000001_a.sql",
			fileMigrations: []FileMigration{
				{ID: "20161101000001", Filename: "20161101000001_a.sql"},
				{ID: "20161101000002", Filename: "20161101000002_b.sql"},
			},
			appliedMigrations: []AppliedMigration{{ID: "20161101000002"}},
		},
	}

	for _, testCase := range testCases {
		migration, err := FilterUpMigrationsByText(
			testCase.filter, testCase.fileMigrations, testCase.appliedMigrations,
		)
		if err != nil {
			t.Errorf("Expected no error, but got: %s", err)
		}

		if migration.Filename != testCase.expectedFilename {
			t.Errorf(
				"Wrong migration: Expected file '%s', but got '%s'",
				testCase.expectedFilename, migration.Filename,
			)
		}
	}
}

func TestFilterUpMigrationsByTextError(t *testing.T) {
	testCases := []struct {
		filter            string
		fileMigrations    []FileMigration
		appliedMigrations []AppliedMigration
	}{

		{
			filter: "2017",
			fileMigrations: []FileMigration{
				{ID: "20161101000001", Filename: "20161101000001_foo.sql"},
			},
			appliedMigrations: []AppliedMigration{},
		},
		{
			filter: "2017",
			fileMigrations: []FileMigration{
				{ID: "20171101000001", Filename: "20171101000001_foo.sql"},
				{ID: "20171101000003", Filename: "20171101000003_buz.sql"},
			},
			appliedMigrations: []AppliedMigration{},
		},

		{
			filter: "2017",
			fileMigrations: []FileMigration{
				{ID: "20171101000001", Filename: "20171101000001_foo.sql"},
			},
			appliedMigrations: []AppliedMigration{{ID: "20171101000001"}},
		},
	}

	for _, testCase := range testCases {
		_, err := FilterUpMigrationsByText(
			testCase.filter, testCase.fileMigrations, testCase.appliedMigrations,
		)
		if err == nil {
			t.Errorf("Expected error, but got none")
		}
	}
}

func TestFilterUpMigrationsByCount(t *testing.T) {
	testCases := []struct {
		count       uint
		all         bool
		fileIDs     []string
		appliedIDs  []string
		expectedIDs []string
	}{
		{
			count: 1, all: false, fileIDs: []string{"1"}, appliedIDs: []string{},
			expectedIDs: []string{"1"},
		},
		{
			count: 1, all: false, fileIDs: []string{"1", "2", "3"}, appliedIDs: []string{"1"},
			expectedIDs: []string{"2"},
		},
		{
			count: 2, all: false, fileIDs: []string{"1", "2", "3"}, appliedIDs: []string{"1"},
			expectedIDs: []string{"2", "3"},
		},
		{
			count: 99, all: false, fileIDs: []string{"1", "2", "3"}, appliedIDs: []string{"1"},
			expectedIDs: []string{"2", "3"},
		},
		{
			count: 0, all: true, fileIDs: []string{"1", "2"}, appliedIDs: []string{},
			expectedIDs: []string{"1", "2"},
		},
		{
			count: 0, all: true, fileIDs: []string{"1", "2", "3"}, appliedIDs: []string{"1"},
			expectedIDs: []string{"2", "3"},
		},
	}

	for _, testCase := range testCases {
		fileMigrations := []FileMigration{}
		for _, id := range testCase.fileIDs {
			fileMigrations = append(fileMigrations, FileMigration{ID: id})
		}
		appliedMigrations := []AppliedMigration{}
		for _, id := range testCase.appliedIDs {
			appliedMigrations = append(appliedMigrations, AppliedMigration{ID: id})
		}
		expectedMigrations := []FileMigration{}
		for _, id := range testCase.expectedIDs {
			expectedMigrations = append(expectedMigrations, FileMigration{ID: id})
		}

		migrations, err := FilterUpMigrationsByCount(
			testCase.count, testCase.all, fileMigrations, appliedMigrations,
		)
		if err != nil {
			t.Errorf("Expected no error, but got: %s", err)
		}

		if diff := pretty.Compare(expectedMigrations, migrations); diff != "" {
			t.Errorf("Did not pass right FileMigrations to migrateUp:\n%s", diff)
		}
	}
}

func TestFilterUpMigrationsByCountNothingLeft(t *testing.T) {
	_, err := FilterUpMigrationsByCount(
		3, false, []FileMigration{{ID: "1"}}, []AppliedMigration{{ID: "1"}},
	)
	if err == nil {
		t.Errorf("Expected error, but got  nothing")
	}
}
